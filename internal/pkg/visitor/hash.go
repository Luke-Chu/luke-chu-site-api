package visitor

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Hash(ip, userAgent, acceptLanguage string) string {
	raw := strings.Join([]string{
		strings.TrimSpace(ip),
		strings.TrimSpace(userAgent),
		strings.TrimSpace(acceptLanguage),
	}, "|")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
