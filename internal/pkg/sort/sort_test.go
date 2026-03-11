package sortutil

import (
	"testing"

	"luke-chu-site-api/internal/constant"
)

func TestNormalizeSortField(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "合法字段", in: "view_count", want: "view_count"},
		{name: "合法字段有空格", in: "  like_count  ", want: "like_count"},
		{name: "非法字段回落默认", in: "unknown", want: constant.DefaultSort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSortField(tt.in)
			if got != tt.want {
				t.Fatalf("NormalizeSortField() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNormalizeSortOrder(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "合法asc", in: "asc", want: "asc"},
		{name: "合法desc", in: "DESC", want: "desc"},
		{name: "非法值回落默认", in: "invalid", want: constant.DefaultSortOrder},
		{name: "空值回落默认", in: "", want: constant.DefaultSortOrder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSortOrder(tt.in)
			if got != tt.want {
				t.Fatalf("NormalizeSortOrder() = %s, want %s", got, tt.want)
			}
		})
	}
}
