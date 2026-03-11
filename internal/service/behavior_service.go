package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"luke-chu-site-api/internal/repository"
)

type BehaviorService interface {
	ViewPhoto(ctx context.Context, photoUUID string) error
	LikePhoto(ctx context.Context, photoUUID, visitorHash string) error
	DownloadPhoto(ctx context.Context, photoUUID string) error
}

type behaviorService struct {
	photoRepo repository.PhotoRepository
}

func NewBehaviorService(photoRepo repository.PhotoRepository) BehaviorService {
	return &behaviorService{photoRepo: photoRepo}
}

func (s *behaviorService) ViewPhoto(ctx context.Context, photoUUID string) error {
	err := s.photoRepo.IncrementViewCount(ctx, photoUUID)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPhotoNotFound
	}
	return fmt.Errorf("view photo failed: %w", err)
}

func (s *behaviorService) LikePhoto(ctx context.Context, photoUUID, visitorHash string) error {
	err := s.photoRepo.AddLike(ctx, photoUUID, visitorHash)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPhotoNotFound
	}
	return fmt.Errorf("like photo failed: %w", err)
}

func (s *behaviorService) DownloadPhoto(ctx context.Context, photoUUID string) error {
	err := s.photoRepo.IncrementDownloadCount(ctx, photoUUID)
	if err == nil || errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPhotoNotFound
	}
	return fmt.Errorf("download photo failed: %w", err)
}
