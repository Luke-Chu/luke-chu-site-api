package service

import (
	"context"
	"errors"
	"fmt"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/repository"
)

type TagService interface {
	ListTags(ctx context.Context) (*response.TagListData, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) ListTags(ctx context.Context) (*response.TagListData, error) {
	tags, err := s.tagRepo.ListTags(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotImplemented) || errors.Is(err, repository.ErrRepositoryNotReady) {
			return &response.TagListData{
				Items: []response.TagItem{
					{Name: "subject", TagType: constant.TagTypeSubject},
					{Name: "element", TagType: constant.TagTypeElement},
					{Name: "mood", TagType: constant.TagTypeMood},
				},
			}, nil
		}
		return nil, fmt.Errorf("list tags failed: %w", err)
	}

	items := make([]response.TagItem, 0, len(tags))
	for _, tag := range tags {
		items = append(items, response.TagItem{
			Name:    tag.Name,
			TagType: tag.TagType,
		})
	}

	return &response.TagListData{Items: items}, nil
}
