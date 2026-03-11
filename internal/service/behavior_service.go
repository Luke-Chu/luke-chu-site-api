package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"luke-chu-site-api/internal/dto/response"
	ossutil "luke-chu-site-api/internal/pkg/oss"
	"luke-chu-site-api/internal/repository"
)

type BehaviorService interface {
	ViewPhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoViewData, error)
	LikePhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error)
	UnlikePhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoUnlikeData, error)
	DownloadPhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoDownloadData, error)
}

type behaviorService struct {
	photoRepo repository.PhotoRepository
	signer    ossutil.DownloadURLSigner
}

func NewBehaviorService(photoRepo repository.PhotoRepository, signer ossutil.DownloadURLSigner) BehaviorService {
	return &behaviorService{
		photoRepo: photoRepo,
		signer:    signer,
	}
}

func (s *behaviorService) ViewPhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoViewData, error) {
	if visitorHash == "" {
		return nil, fmt.Errorf("visitor hash is required")
	}

	count, counted, err := s.photoRepo.IncrementViewCount(ctx, photoUUID, visitorHash)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		return &response.PhotoViewData{UUID: photoUUID, ViewCount: count, Counted: counted}, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPhotoNotFound
	}
	return nil, fmt.Errorf("view photo failed: %w", err)
}

func (s *behaviorService) LikePhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoLikeData, error) {
	if visitorHash == "" {
		return nil, fmt.Errorf("visitor hash is required")
	}

	liked, count, err := s.photoRepo.AddLike(ctx, photoUUID, visitorHash)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		return &response.PhotoLikeData{
			UUID:      photoUUID,
			Liked:     liked,
			LikeCount: count,
		}, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPhotoNotFound
	}
	return nil, fmt.Errorf("like photo failed: %w", err)
}

func (s *behaviorService) UnlikePhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoUnlikeData, error) {
	if visitorHash == "" {
		return nil, fmt.Errorf("visitor hash is required")
	}

	unliked, count, err := s.photoRepo.RemoveLike(ctx, photoUUID, visitorHash)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		return &response.PhotoUnlikeData{
			UUID:      photoUUID,
			Unliked:   unliked,
			LikeCount: count,
		}, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPhotoNotFound
	}
	return nil, fmt.Errorf("unlike photo failed: %w", err)
}

func (s *behaviorService) DownloadPhoto(ctx context.Context, photoUUID, visitorHash string) (*response.PhotoDownloadData, error) {
	if visitorHash == "" {
		return nil, fmt.Errorf("visitor hash is required")
	}
	if s.signer == nil {
		return nil, fmt.Errorf("download signer is not configured")
	}

	count, url, counted, err := s.photoRepo.IncrementDownloadCount(ctx, photoUUID, visitorHash)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		signedURL, signErr := s.signer.SignDownloadURL(ctx, url)
		if signErr != nil {
			return nil, fmt.Errorf("sign download url failed: %w", signErr)
		}
		return &response.PhotoDownloadData{
			UUID:          photoUUID,
			DownloadCount: count,
			DownloadURL:   signedURL,
			Counted:       counted,
		}, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPhotoNotFound
	}
	return nil, fmt.Errorf("download photo failed: %w", err)
}
