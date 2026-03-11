package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/model"
	"luke-chu-site-api/internal/pkg/pager"
	"luke-chu-site-api/internal/repository"
)

var ErrPhotoNotFound = errors.New("photo not found")

type PhotoService interface {
	ListPhotos(ctx context.Context, req request.PhotoListRequest) (*response.PhotoListData, error)
	GetPhotoDetail(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error)
}

type photoService struct {
	photoRepo repository.PhotoRepository
}

func NewPhotoService(photoRepo repository.PhotoRepository) PhotoService {
	return &photoService{photoRepo: photoRepo}
}

func (s *photoService) ListPhotos(ctx context.Context, req request.PhotoListRequest) (*response.PhotoListData, error) {
	req.Normalize()

	// TODO: future query SQL will use parsed keywords.
	_ = req.Keywords()
	_ = req.TagList()

	photos, listErr := s.photoRepo.ListPhotos(ctx, req)
	if listErr != nil &&
		!errors.Is(listErr, repository.ErrNotImplemented) &&
		!errors.Is(listErr, repository.ErrRepositoryNotReady) {
		return nil, fmt.Errorf("list photos failed: %w", listErr)
	}

	total, countErr := s.photoRepo.CountPhotos(ctx, req)
	if countErr != nil &&
		!errors.Is(countErr, repository.ErrNotImplemented) &&
		!errors.Is(countErr, repository.ErrRepositoryNotReady) {
		return nil, fmt.Errorf("count photos failed: %w", countErr)
	}

	if errors.Is(listErr, repository.ErrNotImplemented) || errors.Is(listErr, repository.ErrRepositoryNotReady) {
		photos = []model.Photo{}
	}
	if errors.Is(countErr, repository.ErrNotImplemented) || errors.Is(countErr, repository.ErrRepositoryNotReady) {
		total = 0
	}

	items := make([]response.PhotoListItem, 0, len(photos))
	for _, photo := range photos {
		items = append(items, response.PhotoListItem{
			UUID:          photo.UUID.String(),
			TitleCN:       ptrString(photo.TitleCN),
			TitleEN:       ptrString(photo.TitleEN),
			Orientation:   photo.Orientation,
			ThumbURL:      ptrString(photo.ThumbURL),
			DisplayURL:    ptrString(photo.DisplayURL),
			LikeCount:     photo.LikeCount,
			ViewCount:     photo.ViewCount,
			DownloadCount: photo.DownloadCount,
			ShotTime:      formatTime(photo.ShotTime),
		})
	}

	return &response.PhotoListData{
		Items: items,
		Pagination: response.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: pager.TotalPages(total, req.PageSize),
		},
	}, nil
}

func (s *photoService) GetPhotoDetail(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error) {
	if _, err := uuid.Parse(photoUUID); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	photo, err := s.photoRepo.GetPhotoByUUID(ctx, photoUUID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrPhotoNotFound
		case errors.Is(err, repository.ErrNotImplemented):
			return mockPhotoDetail(photoUUID), nil
		default:
			return nil, fmt.Errorf("get photo detail failed: %w", err)
		}
	}

	return &response.PhotoDetailData{
		UUID:            photo.UUID.String(),
		Filename:        photo.Filename,
		TitleCN:         ptrString(photo.TitleCN),
		TitleEN:         ptrString(photo.TitleEN),
		Description:     ptrString(photo.Description),
		Category:        ptrString(photo.Category),
		ShotTime:        formatTime(photo.ShotTime),
		Width:           photo.Width,
		Height:          photo.Height,
		Orientation:     photo.Orientation,
		Resolution:      ptrString(photo.Resolution),
		CameraModel:     ptrString(photo.CameraModel),
		LensModel:       ptrString(photo.LensModel),
		Aperture:        ptrString(photo.Aperture),
		ShutterSpeed:    ptrString(photo.ShutterSpeed),
		ISO:             ptrInt(photo.ISO),
		FocalLength:     ptrString(photo.FocalLength),
		FocalLength35mm: ptrString(photo.FocalLength35mm),
		MeteringMode:    ptrString(photo.MeteringMode),
		ExposureProgram: ptrString(photo.ExposureProgram),
		WhiteBalance:    ptrString(photo.WhiteBalance),
		Flash:           ptrString(photo.Flash),
		ThumbURL:        ptrString(photo.ThumbURL),
		DisplayURL:      ptrString(photo.DisplayURL),
		OriginalURL:     ptrString(photo.OriginalURL),
		LikeCount:       photo.LikeCount,
		DownloadCount:   photo.DownloadCount,
		ViewCount:       photo.ViewCount,
	}, nil
}

func mockPhotoDetail(photoUUID string) *response.PhotoDetailData {
	return &response.PhotoDetailData{
		UUID:          photoUUID,
		Filename:      "mock.jpg",
		TitleCN:       "示例照片",
		TitleEN:       "Sample Photo",
		Orientation:   "landscape",
		Width:         1920,
		Height:        1280,
		LikeCount:     0,
		ViewCount:     0,
		DownloadCount: 0,
	}
}

func ptrString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func ptrInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
