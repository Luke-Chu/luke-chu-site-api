package sortutil

import (
	"strings"

	"luke-chu-site-api/internal/constant"
)

var allowedFields = map[string]struct{}{
	constant.SortShotTime:  {},
	constant.SortLikeCount: {},
	constant.SortViewCount: {},
	constant.SortDownload:  {},
	constant.SortCreatedAt: {},
}

func Normalize(sortField, order string) (string, string) {
	return NormalizeSortField(sortField), NormalizeSortOrder(order)
}

func NormalizeSortField(sortField string) string {
	sortField = strings.ToLower(strings.TrimSpace(sortField))
	if _, ok := allowedFields[sortField]; !ok {
		sortField = constant.DefaultSort
	}
	return sortField
}

func NormalizeSortOrder(order string) string {
	order = strings.ToLower(strings.TrimSpace(order))
	if order != "asc" && order != "desc" {
		order = constant.DefaultSortOrder
	}
	return order
}

func IsAllowedField(sortField string) bool {
	_, ok := allowedFields[strings.ToLower(strings.TrimSpace(sortField))]
	return ok
}
