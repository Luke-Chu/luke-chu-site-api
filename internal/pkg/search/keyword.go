package search

import (
	"regexp"
	"strings"
)

var keywordSplitter = regexp.MustCompile(`[,\s，、]+`)

func ParseKeywords(raw string) []string {
	parts := keywordSplitter.Split(strings.TrimSpace(raw), -1)
	if len(parts) == 0 {
		return nil
	}

	const maxKeywords = 5
	seen := make(map[string]struct{}, maxKeywords)
	result := make([]string, 0, maxKeywords)

	for _, part := range parts {
		word := strings.TrimSpace(part)
		if word == "" {
			continue
		}
		if _, ok := seen[word]; ok {
			continue
		}
		seen[word] = struct{}{}
		result = append(result, word)
		if len(result) >= maxKeywords {
			break
		}
	}

	return result
}
