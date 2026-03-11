package ossutil

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	aliyunoss "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"luke-chu-site-api/internal/config"
)

const defaultPresignExpireSeconds = 300

// DownloadURLSigner signs public object url to a temporary downloadable url.
type DownloadURLSigner interface {
	SignDownloadURL(ctx context.Context, sourceURL string) (string, error)
}

type PresignDownloadURLSigner struct {
	client         *aliyunoss.Client
	bucketName     string
	publicBaseHost string
	publicBasePath string
	expires        time.Duration
}

func NewPresignDownloadURLSigner(cfg config.OSSConfig) (*PresignDownloadURLSigner, error) {
	bucketName := strings.TrimSpace(cfg.BucketName)
	if bucketName == "" {
		return nil, fmt.Errorf("oss bucket_name is required")
	}

	endpoint, err := normalizeEndpoint(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = inferRegionFromEndpoint(endpoint)
	}
	if region == "" {
		return nil, fmt.Errorf("oss region is required")
	}

	expires := time.Duration(cfg.PresignExpireSecond) * time.Second
	if expires <= 0 {
		expires = defaultPresignExpireSeconds * time.Second
	}

	publicBaseHost, publicBasePath, err := parsePublicBaseURL(cfg.PublicBaseURL)
	if err != nil {
		return nil, err
	}

	credProvider := credentials.NewEnvironmentVariableCredentialsProvider()
	if _, err := credProvider.GetCredentials(context.Background()); err != nil {
		return nil, fmt.Errorf("load oss credentials from env failed: %w", err)
	}

	ossCfg := aliyunoss.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithDisableSSL(false)

	return &PresignDownloadURLSigner{
		client:         aliyunoss.NewClient(ossCfg),
		bucketName:     bucketName,
		publicBaseHost: publicBaseHost,
		publicBasePath: publicBasePath,
		expires:        expires,
	}, nil
}

func (s *PresignDownloadURLSigner) SignDownloadURL(ctx context.Context, sourceURL string) (string, error) {
	objectKey, err := s.objectKeyFromSourceURL(sourceURL)
	if err != nil {
		return "", err
	}

	result, err := s.client.Presign(ctx, &aliyunoss.GetObjectRequest{
		Bucket: aliyunoss.Ptr(s.bucketName),
		Key:    aliyunoss.Ptr(objectKey),
	}, aliyunoss.PresignExpires(s.expires))
	if err != nil {
		return "", fmt.Errorf("presign oss get object failed: %w", err)
	}

	return result.URL, nil
}

func (s *PresignDownloadURLSigner) objectKeyFromSourceURL(sourceURL string) (string, error) {
	raw := strings.TrimSpace(sourceURL)
	if raw == "" {
		return "", fmt.Errorf("source url is empty")
	}

	if !strings.Contains(raw, "://") {
		objectKey := strings.Trim(strings.TrimPrefix(raw, "/"), " ")
		if objectKey == "" {
			return "", fmt.Errorf("invalid object key")
		}
		return objectKey, nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse source url failed: %w", err)
	}

	objectKey := strings.TrimPrefix(u.EscapedPath(), "/")
	if objectKey == "" {
		objectKey = strings.TrimPrefix(u.Path, "/")
	}
	if objectKey == "" {
		return "", fmt.Errorf("missing object key in source url")
	}

	if s.publicBaseHost != "" && strings.EqualFold(u.Hostname(), s.publicBaseHost) && s.publicBasePath != "" {
		prefix := s.publicBasePath + "/"
		if strings.HasPrefix(objectKey, prefix) {
			objectKey = strings.TrimPrefix(objectKey, prefix)
		} else if objectKey == s.publicBasePath {
			objectKey = ""
		}
	}

	objectKey = strings.TrimPrefix(objectKey, "/")
	if objectKey == "" {
		return "", fmt.Errorf("invalid object key in source url")
	}

	return objectKey, nil
}

func normalizeEndpoint(endpoint string) (string, error) {
	value := strings.TrimSpace(endpoint)
	if value == "" {
		return "", fmt.Errorf("oss endpoint is required")
	}

	if strings.Contains(value, "://") {
		u, err := url.Parse(value)
		if err != nil {
			return "", fmt.Errorf("parse oss endpoint failed: %w", err)
		}
		value = strings.TrimSpace(u.Host)
	}

	value = strings.TrimSuffix(value, "/")
	if value == "" {
		return "", fmt.Errorf("oss endpoint is empty")
	}
	return value, nil
}

func inferRegionFromEndpoint(endpoint string) string {
	host := strings.TrimSpace(strings.ToLower(endpoint))
	if host == "" {
		return ""
	}

	parts := strings.Split(host, ".")
	if len(parts) == 0 {
		return ""
	}
	if strings.HasPrefix(parts[0], "oss-") {
		return strings.TrimPrefix(parts[0], "oss-")
	}

	return ""
}

func parsePublicBaseURL(raw string) (string, string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", "", nil
	}

	u, err := url.Parse(value)
	if err != nil {
		return "", "", fmt.Errorf("parse oss public_base_url failed: %w", err)
	}
	if u.Hostname() == "" {
		return "", "", fmt.Errorf("invalid oss public_base_url: host is empty")
	}

	return strings.ToLower(u.Hostname()), strings.Trim(strings.TrimSpace(u.Path), "/"), nil
}
