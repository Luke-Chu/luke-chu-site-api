package integration_test

import (
	"net/http"
	"testing"
)

type commonEnvelope struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func TestHealthAPI(t *testing.T) {
	router := newTestRouter(testServices{})
	rr := newRequestRecorder(router, http.MethodGet, "/api/v1/health")

	if rr.Code != http.StatusOK {
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}

	resp := decodeBody[commonEnvelope](t, rr)
	if resp.Code != 0 {
		t.Fatalf("code = %d, want 0", resp.Code)
	}
	if resp.Data["status"] != "ok" {
		t.Fatalf("status = %v, want ok", resp.Data["status"])
	}
}
