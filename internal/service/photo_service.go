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
	"luke-chu-site-api/internal/pkg/pager"
	"luke-chu-site-api/internal/repository"
)

var ErrPhotoNotFound = errors.New("photo not found")

type PhotoService interface {
	ListPhotos(ctx context.Context, req *request.PhotoListRequest) (*response.PhotoListData, error)
	GetPhotoDetail(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error)
}

type photoService struct {
	photoRepo repository.PhotoRepository
}

func NewPhotoService(photoRepo repository.PhotoRepository) PhotoService {
	return &photoService{photoRepo: photoRepo}
}

func (s *photoService) ListPhotos(ctx context.Context, req *request.PhotoListRequest) (*response.PhotoListData, error) {
	if req == nil {
		req = &request.PhotoListRequest{}
	}
	req.Normalize()

	photos, err := s.photoRepo.ListPhotos(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list photos failed: %w", err)
	}

	total, err := s.photoRepo.CountPhotos(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("count photos failed: %w", err)
	}

	photoIDs := make([]int64, 0, len(photos))
	for _, photo := range photos {
		photoIDs = append(photoIDs, photo.ID)
	}

	tagMap, err := s.photoRepo.ListPhotoTagsByPhotoIDs(ctx, photoIDs)
	if err != nil {
		return nil, fmt.Errorf("list photo tags failed: %w", err)
	}

	list := make([]response.PhotoListItem, 0, len(photos))
	for _, photo := range photos {
		list = append(list, response.PhotoListItem{
			ID:            photo.ID,
			UUID:          photo.UUID.String(),
			Filename:      photo.Filename,
			TitleCN:       ptrString(photo.TitleCN),
			TitleEN:       ptrString(photo.TitleEN),
			ThumbURL:      ptrString(photo.ThumbURL),
			DisplayURL:    ptrString(photo.DisplayURL),
			Width:         photo.Width,
			Height:        photo.Height,
			Orientation:   photo.Orientation,
			ShotTime:      formatTime(photo.ShotTime),
			Aperture:      ptrString(photo.Aperture),
			ShutterSpeed:  ptrString(photo.ShutterSpeed),
			ISO:           ptrInt(photo.ISO),
			LikeCount:     photo.LikeCount,
			ViewCount:     photo.ViewCount,
			DownloadCount: photo.DownloadCount,
			Tags:          tagMap[photo.ID],
		})
	}

	return &response.PhotoListData{
		List: list,
		Pagination: response.Pagination{
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: pager.TotalPages(total, req.PageSize),
		},
		Query: response.PhotoListQuery{
			Q:           req.Q,
			Keywords:    req.KeywordList(),
			Sort:        req.Sort,
			Order:       req.Order,
			Tags:        req.TagList(),
			TagMode:     req.TagMode,
			Orientation: req.Orientation,
			Year:        req.Year,
			Month:       req.Month,
			Category:    req.Category,
		},
	}, nil
}

func (s *photoService) GetPhotoDetail(ctx context.Context, photoUUID string) (*response.PhotoDetailData, error) {
	if _, err := uuid.Parse(photoUUID); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	photo, err := s.photoRepo.GetPhotoDetailByUUID(ctx, photoUUID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrPhotoNotFound
		default:
			return nil, fmt.Errorf("get photo detail failed: %w", err)
		}
	}

	tags, err := s.photoRepo.GetPhotoTagsByPhotoID(ctx, photo.ID)
	if err != nil {
		return nil, fmt.Errorf("get photo detail tags failed: %w", err)
	}

	return &response.PhotoDetailData{
		ID:              photo.ID,
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
		FocalLength:     ptrFloat(photo.FocalLength),
		FocalLength35mm: ptrFloat(photo.FocalLength35mm),
		MeteringMode:    ptrString(photo.MeteringMode),
		ExposureComp:    ptrString(photo.ExposureComp),
		ExposureProgram: ptrString(photo.ExposureProgram),
		WhiteBalance:    ptrString(photo.WhiteBalance),
		Flash:           ptrString(photo.Flash),
		ThumbURL:        ptrString(photo.ThumbURL),
		DisplayURL:      ptrString(photo.DisplayURL),
		OriginalURL:     ptrString(photo.OriginalURL),
		LikeCount:       photo.LikeCount,
		DownloadCount:   photo.DownloadCount,
		ViewCount:       photo.ViewCount,
		CreatedAt:       photo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       photo.UpdatedAt.Format(time.RFC3339),
		Tags:            tags,
	}, nil
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

func ptrFloat(v *float64) float64 {
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
