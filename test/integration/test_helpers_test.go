package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"luke-chu-site-api/internal/app"
	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/handler"
	"luke-chu-site-api/internal/service"
)

type stubPhotoService struct {
	listFn   func(ctx context.Context, req *request.PhotoListRequest) (*response.PhotoListData, error)
	detailFn func(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error)
}

func (s *stubPhotoService) ListPhotos(ctx context.Context, req *request.PhotoListRequest) (*response.PhotoListData, error) {
	if s.listFn != nil {
		return s.listFn(ctx, req)
	}
	return &response.PhotoListData{}, nil
}

func (s *stubPhotoService) GetPhotoDetail(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error) {
	if s.detailFn != nil {
		return s.detailFn(ctx, photoUUID)
	}
	return &response.PhotoDetailData{}, nil
}

type stubBehaviorService struct {
	viewFn     func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoViewData, error)
	likeFn     func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error)
	unlikeFn   func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoUnlikeData, error)
	downloadFn func(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoDownloadData, error)
}

func (s *stubBehaviorService) ViewPhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoViewData, error) {
	if s.viewFn != nil {
		return s.viewFn(ctx, photoUUID, visitorHash)
	}
	return &response.PhotoViewData{UUID: photoUUID, ViewCount: 1, Counted: true}, nil
}

func (s *stubBehaviorService) LikePhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error) {
	if s.likeFn != nil {
		return s.likeFn(ctx, photoUUID, visitorHash)
	}
	return &response.PhotoLikeData{UUID: photoUUID, Liked: true, LikeCount: 1}, nil
}

func (s *stubBehaviorService) UnlikePhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoUnlikeData, error) {
	if s.unlikeFn != nil {
		return s.unlikeFn(ctx, photoUUID, visitorHash)
	}
	return &response.PhotoUnlikeData{UUID: photoUUID, Unliked: true, LikeCount: 0}, nil
}

func (s *stubBehaviorService) DownloadPhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoDownloadData, error) {
	if s.downloadFn != nil {
		return s.downloadFn(ctx, photoUUID, visitorHash)
	}
	return &response.PhotoDownloadData{UUID: photoUUID, DownloadCount: 1, DownloadURL: "https://example.com/a.jpg", Counted: true}, nil
}

type stubTagService struct {
	listFn func(ctx context.Context) (*response.TagListData, error)
}

func (s *stubTagService) ListTags(ctx context.Context) (*response.TagListData, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}
	return &response.TagListData{}, nil
}

type stubFilterService struct {
	getFn func(ctx context.Context) (*response.FilterData, error)
}

func (s *stubFilterService) GetFilters(ctx context.Context) (*response.FilterData, error) {
	if s.getFn != nil {
		return s.getFn(ctx)
	}
	return &response.FilterData{}, nil
}

type testServices struct {
	photo    service.PhotoService
	behavior service.BehaviorService
	tag      service.TagService
	filter   service.FilterService
}

func newTestRouter(s testServices) http.Handler {
	gin.SetMode(gin.TestMode)

	if s.photo == nil {
		s.photo = &stubPhotoService{}
	}
	if s.behavior == nil {
		s.behavior = &stubBehaviorService{}
	}
	if s.tag == nil {
		s.tag = &stubTagService{}
	}
	if s.filter == nil {
		s.filter = &stubFilterService{}
	}

	healthHandler := handler.NewHealthHandler("luke-chu-site-api")
	photoHandler := handler.NewPhotoHandler(s.photo, s.behavior, validator.New())
	tagHandler := handler.NewTagHandler(s.tag)
	filterHandler := handler.NewFilterHandler(s.filter)

	return app.NewRouter(zap.NewNop(), healthHandler, photoHandler, tagHandler, filterHandler)
}

func newRequestRecorder(router http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("User-Agent", "integration-test-ua")
	req.Header.Set("Accept-Language", "zh-CN")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func decodeBody[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, rr.Body.String())
	}
	return out
}
