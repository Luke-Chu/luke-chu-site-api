package request

import (
	"strings"

	"luke-chu-site-api/internal/pkg/pager"
	"luke-chu-site-api/internal/pkg/search"
	sortutil "luke-chu-site-api/internal/pkg/sort"
)

type PhotoListRequest struct {
	Q           string `form:"q" validate:"omitempty,max=120"`
	Page        int    `form:"page" validate:"omitempty,min=1"`
	PageSize    int    `form:"pageSize" validate:"omitempty,min=1,max=100"`
	Sort        string `form:"sort" validate:"omitempty,max=30"`
	Order       string `form:"order" validate:"omitempty,oneof=asc desc ASC DESC"`
	Tags        string `form:"tags" validate:"omitempty,max=200"`
	Orientation string `form:"orientation" validate:"omitempty,oneof=landscape portrait square"`
	Year        int    `form:"year" validate:"omitempty,min=1900,max=2100"`
	Month       int    `form:"month" validate:"omitempty,min=1,max=12"`
}

func (r *PhotoListRequest) Normalize() {
	r.Page, r.PageSize = pager.Normalize(r.Page, r.PageSize)
	r.Sort, r.Order = sortutil.Normalize(r.Sort, r.Order)
}

func (r PhotoListRequest) Keywords() []string {
	return search.ParseKeywords(r.Q)
}

func (r PhotoListRequest) TagList() []string {
	if strings.TrimSpace(r.Tags) == "" {
		return nil
	}
	return search.ParseKeywords(r.Tags)
}
