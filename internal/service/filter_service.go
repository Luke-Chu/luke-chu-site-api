package service

import (
	"context"
	"fmt"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/repository"
)

type FilterService interface {
	GetFilters(ctx context.Context) (*response.FilterData, error)
}

type filterService struct {
	filterRepo repository.FilterRepository
}

func NewFilterService(filterRepo repository.FilterRepository) FilterService {
	return &filterService{filterRepo: filterRepo}
}

func (s *filterService) GetFilters(ctx context.Context) (*response.FilterData, error) {
	years, err := s.filterRepo.ListAvailableYears(ctx)
	if err != nil {
		return nil, fmt.Errorf("list years failed: %w", err)
	}

	categories, err := s.filterRepo.ListAvailableCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories failed: %w", err)
	}

	orientations, err := s.filterRepo.ListOrientationCounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("list orientations failed: %w", err)
	}

	tags, err := s.filterRepo.ListAllTagsGrouped(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tags failed: %w", err)
	}

	return &response.FilterData{
		Years:        years,
		Categories:   categories,
		Orientations: orientations,
		TagTypes: []string{
			constant.TagTypeSubject,
			constant.TagTypeElement,
			constant.TagTypeMood,
		},
		Tags: tags,
	}, nil
}
