package integration_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

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
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}

	resp := decodeBody[commonEnvelope](t, rr)
	if resp.Code != 0 {
		t.Fatalf("code = %d, want 0", resp.Code)
	}
	if captured.Q != "天空,风筝" || captured.Page != 2 || captured.PageSize != 20 || captured.Sort != "view_count" || captured.Order != "asc" {
		t.Fatalf("query 参数绑定不正确: %+v", captured)
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

	t.Run("合法uuid返回200", func(t *testing.T) {
		rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos/"+validUUID)
		if rr.Code != http.StatusOK {
			t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("非法uuid返回400", func(t *testing.T) {
		rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos/not-a-uuid")
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("不存在返回404", func(t *testing.T) {
		rr := newRequestRecorder(router, http.MethodGet, "/api/v1/photos/2d4f2e5b-4aa1-4ed1-a32a-111111111111")
		if rr.Code != http.StatusNotFound {
			t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusNotFound)
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
				t.Fatal("visitorHash 不应为空")
			}
			return &response.PhotoViewData{UUID: photoUUID, ViewCount: 11, Counted: true}, nil
		},
	}
	router := newTestRouter(testServices{behavior: behaviorSvc})

	rr := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/"+validUUID+"/view")
	if rr.Code != http.StatusOK {
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}

	bad := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/bad/view")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("状态码 = %d, want %d", bad.Code, http.StatusBadRequest)
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
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}
	if len(capturedHash) == 0 {
		t.Fatal("visitor hash 未传入 behavior service")
	}

	bad := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/bad/like")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("状态码 = %d, want %d", bad.Code, http.StatusBadRequest)
	}
}

func TestPhotoDownloadAPI(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	behaviorSvc := &stubBehaviorService{
		downloadFn: func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoDownloadData, error) {
			if visitorHash == "" {
				t.Fatal("download 接口 visitorHash 不应为空")
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
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}

	bad := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/not-uuid/download")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("状态码 = %d, want %d", bad.Code, http.StatusBadRequest)
	}
}

func TestLikeMiddlewareHashLooksValid(t *testing.T) {
	behaviorSvc := &stubBehaviorService{
		likeFn: func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error) {
			if len(visitorHash) != 64 {
				t.Fatalf("visitor hash 长度 = %d, want 64", len(visitorHash))
			}
			if strings.Trim(visitorHash, "0123456789abcdef") != "" {
				t.Fatalf("visitor hash 不是十六进制: %s", visitorHash)
			}
			return &response.PhotoLikeData{UUID: photoUUID, Liked: true, LikeCount: 1}, nil
		},
	}
	router := newTestRouter(testServices{behavior: behaviorSvc})
	rr := newRequestRecorder(router, http.MethodPost, "/api/v1/photos/550e8400-e29b-41d4-a716-446655440000/like")
	if rr.Code != http.StatusOK {
		t.Fatalf("状态码 = %d, want %d", rr.Code, http.StatusOK)
	}
}
