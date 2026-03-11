package request

import (
	"regexp"
	"strings"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/pkg/search"
	sortutil "luke-chu-site-api/internal/pkg/sort"
)

var tagSplitter = regexp.MustCompile(`[,\s，、]+`)

type PhotoListRequest struct {
	Q           string `form:"q" validate:"omitempty,max=120"`
	Page        int    `form:"page" validate:"omitempty,min=1"`
	PageSize    int    `form:"pageSize" validate:"omitempty,min=1,max=60"`
	Sort        string `form:"sort" validate:"omitempty,max=30"`
	Order       string `form:"order" validate:"omitempty,max=10"`
	Tags        string `form:"tags" validate:"omitempty,max=200"`
	Orientation string `form:"orientation" validate:"omitempty,max=20"`
	Year        int    `form:"year" validate:"omitempty,min=1900,max=2100"`
	Month       int    `form:"month" validate:"omitempty,min=1,max=12"`
	Category    string `form:"category" validate:"omitempty,max=100"`
	TagMode     string `form:"tagMode" validate:"omitempty,max=10"`
}

func (r *PhotoListRequest) Normalize() {
	r.Q = strings.TrimSpace(r.Q)
	r.Tags = strings.TrimSpace(r.Tags)
	r.Orientation = strings.ToLower(strings.TrimSpace(r.Orientation))
	r.Category = strings.TrimSpace(r.Category)
	r.TagMode = strings.ToLower(strings.TrimSpace(r.TagMode))

	if r.Page < 1 {
		r.Page = 1
	}
	switch {
	case r.PageSize <= 0:
		r.PageSize = 30
	case r.PageSize > 60:
		r.PageSize = 60
	}

	r.Sort = sortutil.NormalizeSortField(r.Sort)
	r.Order = sortutil.NormalizeSortOrder(r.Order)

	if r.TagMode != "any" && r.TagMode != "all" {
		r.TagMode = "any"
	}

	switch r.Orientation {
	case constant.OrientationLandscape, constant.OrientationPortrait, constant.OrientationSquare:
	default:
		r.Orientation = ""
	}

	if r.Month < 1 || r.Month > 12 {
		r.Month = 0
	}
	if r.Year < 1900 || r.Year > 2100 {
		r.Year = 0
	}
}

func (r PhotoListRequest) KeywordList() []string {
	return search.ParseKeywords(r.Q)
}

func (r PhotoListRequest) TagList() []string {
	if strings.TrimSpace(r.Tags) == "" {
		return nil
	}
	parts := tagSplitter.Split(r.Tags, -1)
	seen := make(map[string]struct{}, len(parts))
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		lower := strings.ToLower(name)
		if _, ok := seen[lower]; ok {
			continue
		}
		seen[lower] = struct{}{}
		result = append(result, name)
	}
	return result
}
