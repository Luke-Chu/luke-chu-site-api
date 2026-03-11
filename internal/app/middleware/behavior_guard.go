package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/response"
)

type BehaviorGuardConfig struct {
	Enabled                    bool
	WindowSeconds              int64
	IPLimitPerWindow           int
	SuspiciousIPLimitPerWindow int
}

func DefaultBehaviorGuardConfig() BehaviorGuardConfig {
	return BehaviorGuardConfig{
		Enabled:                    true,
		WindowSeconds:              60,
		IPLimitPerWindow:           120,
		SuspiciousIPLimitPerWindow: 20,
	}
}

type ipWindowCounter struct {
	windowSeconds int64
	mu            sync.Mutex
	counters      map[string]counterValue
}

type counterValue struct {
	Bucket int64
	Count  int
}

func newIPWindowCounter(windowSeconds int64) *ipWindowCounter {
	if windowSeconds <= 0 {
		windowSeconds = 60
	}
	return &ipWindowCounter{
		windowSeconds: windowSeconds,
		counters:      make(map[string]counterValue, 1024),
	}
}

func (c *ipWindowCounter) increment(key string, now time.Time) int {
	bucket := now.Unix() / c.windowSeconds

	c.mu.Lock()
	defer c.mu.Unlock()

	current := c.counters[key]
	if current.Bucket != bucket {
		current.Bucket = bucket
		current.Count = 0
	}
	current.Count++
	c.counters[key] = current

	// Best-effort cleanup to control memory growth.
	if len(c.counters) > 10000 {
		expireBefore := bucket - 2
		for k, v := range c.counters {
			if v.Bucket < expireBefore {
				delete(c.counters, k)
			}
		}
	}

	return current.Count
}

func BehaviorGuard(logger *zap.Logger, cfg BehaviorGuardConfig) gin.HandlerFunc {
	if logger == nil {
		logger = zap.NewNop()
	}
	if cfg.WindowSeconds <= 0 {
		cfg.WindowSeconds = 60
	}
	if cfg.IPLimitPerWindow <= 0 {
		cfg.IPLimitPerWindow = 120
	}
	if cfg.SuspiciousIPLimitPerWindow <= 0 {
		cfg.SuspiciousIPLimitPerWindow = 20
	}

	ipLimiter := newIPWindowCounter(cfg.WindowSeconds)
	suspiciousLimiter := newIPWindowCounter(cfg.WindowSeconds)

	return func(c *gin.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}

		ip := c.ClientIP()
		ua := strings.TrimSpace(c.GetHeader("User-Agent"))
		visitorHash, _ := c.Get(VisitorHashKey)
		hash, _ := visitorHash.(string)

		currentCount := ipLimiter.increment(ip, time.Now())
		if currentCount > cfg.IPLimitPerWindow {
			logger.Warn("behavior guard blocked by ip rate limit",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("visitor_hash", hash),
				zap.Int("count_in_window", currentCount),
				zap.Int("limit", cfg.IPLimitPerWindow),
				zap.Int64("window_seconds", cfg.WindowSeconds),
			)
			response.Error(c, 429, constant.CodeTooManyBehaviorRequests, constant.MsgTooManyBehaviorRequests)
			c.Abort()
			return
		}

		suspicious, reason := isSuspiciousUserAgent(ua)
		if suspicious {
			suspiciousCount := suspiciousLimiter.increment(ip, time.Now())
			logger.Warn("behavior guard suspicious user-agent detected",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("visitor_hash", hash),
				zap.String("reason", reason),
				zap.String("user_agent", ua),
				zap.Int("count_in_window", suspiciousCount),
				zap.Int("limit", cfg.SuspiciousIPLimitPerWindow),
				zap.Int64("window_seconds", cfg.WindowSeconds),
			)

			if suspiciousCount > cfg.SuspiciousIPLimitPerWindow {
				response.Error(c, 429, constant.CodeSuspiciousBehavior, constant.MsgSuspiciousBehavior)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func isSuspiciousUserAgent(userAgent string) (bool, string) {
	ua := strings.TrimSpace(strings.ToLower(userAgent))
	if ua == "" {
		return true, "empty user-agent"
	}

	patterns := []string{
		"python-requests",
		"go-http-client",
		"scrapy",
		"sqlmap",
		"nmap",
		"masscan",
		"nikto",
		"zgrab",
		"crawler",
		"spider",
		"bot/",
	}
	for _, pattern := range patterns {
		if strings.Contains(ua, pattern) {
			return true, "matched pattern: " + pattern
		}
	}

	return false, ""
}
