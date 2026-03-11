package middleware

import (
	"github.com/gin-gonic/gin"

	visitorkey "luke-chu-site-api/internal/pkg/visitor"
)

const VisitorHashKey = "visitor_hash"

func Visitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := visitorkey.Hash(
			c.ClientIP(),
			c.GetHeader("User-Agent"),
			c.GetHeader("Accept-Language"),
		)
		c.Set(VisitorHashKey, hash)
		c.Next()
	}
}
