package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/response"
)

type FilterRepository interface {
	ListAvailableYears(ctx context.Context) ([]int, error)
	ListAvailableCategories(ctx context.Context) ([]string, error)
	ListOrientationCounts(ctx context.Context) ([]response.OrientationOption, error)
	ListAllTagsGrouped(ctx context.Context) (map[string][]response.TagItem, error)
}

type SQLXFilterRepository struct {
	db *sqlx.DB
}

func NewFilterRepository(db *sqlx.DB) FilterRepository {
	return &SQLXFilterRepository{db: db}
}

func (r *SQLXFilterRepository) ListAvailableYears(ctx context.Context) ([]int, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}
	years := make([]int, 0)
	err := r.db.SelectContext(ctx, &years, `
SELECT DISTINCT p.year
FROM photos p
WHERE p.is_published = TRUE
  AND p.year IS NOT NULL
  AND p.year > 0
ORDER BY p.year DESC
`)
	if err != nil {
		return nil, fmt.Errorf("list years failed: %w", err)
	}
	return years, nil
}

func (r *SQLXFilterRepository) ListAvailableCategories(ctx context.Context) ([]string, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}
	categories := make([]string, 0)
	err := r.db.SelectContext(ctx, &categories, `
SELECT DISTINCT p.category
FROM photos p
WHERE p.is_published = TRUE
  AND COALESCE(TRIM(p.category), '') <> ''
ORDER BY p.category ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list categories failed: %w", err)
	}
	return categories, nil
}

func (r *SQLXFilterRepository) ListOrientationCounts(ctx context.Context) ([]response.OrientationOption, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}
	type row struct {
		Name  string `db:"name"`
		Count int64  `db:"count"`
	}
	rows := make([]row, 0)
	err := r.db.SelectContext(ctx, &rows, `
SELECT p.orientation AS name, COUNT(1) AS count
FROM photos p
WHERE p.is_published = TRUE
  AND p.orientation IN ('landscape', 'portrait', 'square')
GROUP BY p.orientation
`)
	if err != nil {
		return nil, fmt.Errorf("list orientation counts failed: %w", err)
	}

	countMap := map[string]int64{
		constant.OrientationLandscape: 0,
		constant.OrientationPortrait:  0,
		constant.OrientationSquare:    0,
	}
	for _, item := range rows {
		name := strings.ToLower(strings.TrimSpace(item.Name))
		if _, ok := countMap[name]; ok {
			countMap[name] = item.Count
		}
	}

	return []response.OrientationOption{
		{Name: constant.OrientationLandscape, Count: countMap[constant.OrientationLandscape]},
		{Name: constant.OrientationPortrait, Count: countMap[constant.OrientationPortrait]},
		{Name: constant.OrientationSquare, Count: countMap[constant.OrientationSquare]},
	}, nil
}

func (r *SQLXFilterRepository) ListAllTagsGrouped(ctx context.Context) (map[string][]response.TagItem, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}
	type row struct {
		ID      int64  `db:"id"`
		Name    string `db:"name"`
		TagType string `db:"tag_type"`
	}
	rows := make([]row, 0)
	err := r.db.SelectContext(ctx, &rows, `
SELECT DISTINCT
	t.id,
	t.name,
	t.tag_type
FROM tags t
JOIN photo_tags pt ON pt.tag_id = t.id
JOIN photos p ON p.id = pt.photo_id
WHERE p.is_published = TRUE
ORDER BY t.tag_type ASC, t.name ASC, t.id ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list grouped tags failed: %w", err)
	}

	grouped := map[string][]response.TagItem{
		constant.TagTypeSubject: {},
		constant.TagTypeElement: {},
		constant.TagTypeMood:    {},
	}

	for _, item := range rows {
		tagType := strings.ToLower(strings.TrimSpace(item.TagType))
		if _, ok := grouped[tagType]; !ok {
			grouped[tagType] = []response.TagItem{}
		}
		grouped[tagType] = append(grouped[tagType], response.TagItem{
			ID:      item.ID,
			Name:    item.Name,
			TagType: tagType,
		})
	}

	return grouped, nil
}
