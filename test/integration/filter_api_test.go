package integration_test

import (
	"context"
	"net/http"
	"testing"

	"luke-chu-site-api/internal/dto/response"
)

func TestFiltersAPI(t *testing.T) {
	filterSvc := &stubFilterService{
		getFn: func(ctx context.Context) (*response.FilterData, error) {
			return &response.FilterData{
				Years:      []int{2026, 2025},
				Categories: []string{"旅行", "城市"},
				Orientations: []response.OrientationOption{
					{Name: "landscape", Count: 10},
				},
				TagTypes: []string{"subject", "element", "mood"},
				Tags: map[string][]response.TagItem{
					"subject": {
						{ID: 1, Name: "风光", TagType: "subject"},
					},
				},
			}, nil
		},
	}
	router := newTestRouter(testServices{filter: filterSvc})
	rr := newRequestRecorder(router, http.MethodGet, "/api/v1/filters")

	if rr.Code != http.StatusOK {
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}

	resp := decodeBody[commonEnvelope](t, rr)
	if resp.Code != 0 {
		t.Fatalf("code = %d, want 0", resp.Code)
	}

	data := resp.Data
	for _, key := range []string{"years", "categories", "orientations", "tagTypes", "tags"} {
		if _, ok := data[key]; !ok {
			t.Fatalf("响应缺少字段: %s", key)
		}
	}
}
