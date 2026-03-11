package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"luke-chu-site-api/internal/app/middleware"
	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/service"
)

func TestPhotosListAPI(t *testing.T) {
	var captured request.PhotoListRequest
	photoSvc := &stubPhotoService{
		listFn: func(ctx context.Context, req *request.PhotoListRequest) (*response.PhotoListData, error) {
			captured = *req
			return &response.PhotoListData{
				List: []response.PhotoListItem{
					{ID: 1, UUID: "550e8400-e29b-41d4-a716-446655440000", Filename: "a.jpg"},
				},
				Pagination: response.Pagination{Page: 1, PageSize: 30, Total: 1, TotalPages: 1},
			}, nil
		},
	}

	router := newTestRouter(testServices{photo: photoSvc})
	rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos?q=天空,风筝&page=2&pageSize=20&sort=view_count&order=asc")

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	resp := decodeBody[commonEnvelope](t, rr)
	if resp.Code != 0 {
		t.Fatalf("code = %d, want 0", resp.Code)
	}
	if captured.Q != "天空,风筝" || captured.Page != 2 || captured.PageSize != 20 || captured.Sort != "view_count" || captured.Order != "asc" {
		t.Fatalf("unexpected bound query: %+v", captured)
	}
}

func TestPhotoDetailAPI(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	photoSvc := &stubPhotoService{
		detailFn: func(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error) {
			if photoUUID == validUUID {
				return &response.PhotoDetailData{
					ID:           10,
					UUID:         validUUID,
					Filename:     "detail.jpg",
					ExposureComp: "0EV",
				}, nil
			}
			return nil, service.ErrPhotoNotFound
		},
	}
	router := newTestRouter(testServices{photo: photoSvc})

	t.Run("valid uuid returns 200", func(t *testing.T) {
		rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos/"+validUUID)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("invalid uuid returns 400", func(t *testing.T) {
		rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos/not-a-uuid")
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("nonexistent returns 404", func(t *testing.T) {
		rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos/2d4f2e5b-4aa1-4ed1-a32a-111111111111")
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestPhotoViewAPI(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	behaviorSvc := &stubBehaviorService{
		viewFn: func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoViewData, error) {
			if photoUUID != validUUID {
				return nil, service.ErrPhotoNotFound
			}
			if visitorHash == "" {
				t.Fatal("visitorHash should not be empty")
			}
			return &response.PhotoViewData{UUID: photoUUID, ViewCount: 11, Counted: true}, nil
		},
	}
	router := newTestRouter(testServices{behavior: behaviorSvc})

	rr := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/"+validUUID+"/view")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	bad := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/bad/view")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", bad.Code, http.StatusBadRequest)
	}
}

func TestPhotoLikeAPI(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	var capturedHash string
	behaviorSvc := &stubBehaviorService{
		likeFn: func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error) {
			capturedHash = visitorHash
			return &response.PhotoLikeData{UUID: photoUUID, Liked: true, LikeCount: 3}, nil
		},
	}
	router := newTestRouter(testServices{behavior: behaviorSvc})

	rr := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/"+validUUID+"/like")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if len(capturedHash) == 0 {
		t.Fatal("visitor hash should be passed to behavior service")
	}

	bad := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/bad/like")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", bad.Code, http.StatusBadRequest)
	}
}

func TestPhotoDownloadAPI(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	behaviorSvc := &stubBehaviorService{
		downloadFn: func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoDownloadData, error) {
			if visitorHash == "" {
				t.Fatal("download visitorHash should not be empty")
			}
			return &response.PhotoDownloadData{
				UUID:          photoUUID,
				DownloadCount: 20,
				DownloadURL:   "https://example.com/original.jpg",
				Counted:       true,
			}, nil
		},
	}
	router := newTestRouter(testServices{behavior: behaviorSvc})

	rr := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/"+validUUID+"/download")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	bad := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/not-uuid/download")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", bad.Code, http.StatusBadRequest)
	}
}

func TestLikeMiddlewareHashLooksValid(t *testing.T) {
	behaviorSvc := &stubBehaviorService{
		likeFn: func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error) {
			if len(visitorHash) != 64 {
				t.Fatalf("visitor hash length = %d, want 64", len(visitorHash))
			}
			if strings.Trim(visitorHash, "0123456789abcdef") != "" {
				t.Fatalf("visitor hash should be hex string: %s", visitorHash)
			}
			return &response.PhotoLikeData{UUID: photoUUID, Liked: true, LikeCount: 1}, nil
		},
	}
	router := newTestRouter(testServices{behavior: behaviorSvc})
	rr := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/like")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestBehaviorGuardIPRateLimit(t *testing.T) {
	router := newTestRouter(testServices{
		behavior: &stubBehaviorService{},
		behaviorGuard: middleware.BehaviorGuardConfig{
			Enabled:                    true,
			WindowSeconds:              60,
			IPLimitPerWindow:           1,
			SuspiciousIPLimitPerWindow: 5,
		},
	})

	first := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/view")
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	second := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/view")
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}

	body := decodeBody[commonEnvelope](t, second)
	if body.Code != constant.CodeTooManyBehaviorRequests {
		t.Fatalf("code = %d, want %d", body.Code, constant.CodeTooManyBehaviorRequests)
	}
}

func TestBehaviorGuardSuspiciousUserAgentLimit(t *testing.T) {
	router := newTestRouter(testServices{
		behavior: &stubBehaviorService{},
		behaviorGuard: middleware.BehaviorGuardConfig{
			Enabled:                    true,
			WindowSeconds:              60,
			IPLimitPerWindow:           10,
			SuspiciousIPLimitPerWindow: 1,
		},
	})

	path := "/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/like"
	first := newRequestRecorderWithUA(router, http.MethodPost, path, "python-requests/2.31")
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	second := newRequestRecorderWithUA(router, http.MethodPost, path, "python-requests/2.31")
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}
	body := decodeBody[commonEnvelope](t, second)
	if body.Code != constant.CodeSuspiciousBehavior {
		t.Fatalf("code = %d, want %d", body.Code, constant.CodeSuspiciousBehavior)
	}
}

func newRequestRecorderWithUA(router http.Handler, method, path, userAgent string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Language", "zh-CN")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
