package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/model"
	"luke-chu-site-api/internal/pkg/pager"
	sortutil "luke-chu-site-api/internal/pkg/sort"
)

var (
	ErrRepositoryNotReady = errors.New("repository not ready")
	ErrNotImplemented     = errors.New("repository sql not implemented for current schema")
)

type PhotoRepository interface {
	ListPhotos(ctx context.Context, req request.PhotoListRequest) ([]model.Photo, error)
	CountPhotos(ctx context.Context, req request.PhotoListRequest) (int64, error)
	GetPhotoByUUID(ctx context.Context, uuid string) (*model.Photo, error)
	IncrementViewCount(ctx context.Context, uuid string) error
	IncrementDownloadCount(ctx context.Context, uuid string) error
	AddLike(ctx context.Context, uuid, visitorHash string) error
}

type SQLXPhotoRepository struct {
	db *sqlx.DB
}

func NewPhotoRepository(db *sqlx.DB) PhotoRepository {
	return &SQLXPhotoRepository{db: db}
}

func (r *SQLXPhotoRepository) ListPhotos(ctx context.Context, req request.PhotoListRequest) ([]model.Photo, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}

	ok, err := r.tableExists(ctx, "photos")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotImplemented
	}

	sortField, order := sortutil.Normalize(req.Sort, req.Order)
	query := fmt.Sprintf(`
SELECT
	id, uuid, filename, title_cn, title_en, description, category, shot_time, width, height, orientation, resolution,
	camera_model, lens_model, aperture, shutter_speed, iso, focal_length, focal_length_35mm, metering_mode,
	exposure_program, white_balance, flash, thumb_url, display_url, original_url, like_count, download_count,
	view_count, is_published, created_at, updated_at
FROM photos
WHERE is_published = TRUE
ORDER BY %s %s
LIMIT $1 OFFSET $2
`, sortField, order)

	items := make([]model.Photo, 0)
	err = r.db.SelectContext(ctx, &items, query, req.PageSize, pager.Offset(req.Page, req.PageSize))
	if err != nil {
		return nil, fmt.Errorf("list photos failed: %w", err)
	}

	return items, nil
}

func (r *SQLXPhotoRepository) CountPhotos(ctx context.Context, _ request.PhotoListRequest) (int64, error) {
	if r.db == nil {
		return 0, ErrRepositoryNotReady
	}

	ok, err := r.tableExists(ctx, "photos")
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, ErrNotImplemented
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(1) FROM photos WHERE is_published = TRUE`); err != nil {
		return 0, fmt.Errorf("count photos failed: %w", err)
	}
	return total, nil
}

func (r *SQLXPhotoRepository) GetPhotoByUUID(ctx context.Context, uuid string) (*model.Photo, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}

	ok, err := r.tableExists(ctx, "photos")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotImplemented
	}

	var photo model.Photo
	err = r.db.GetContext(ctx, &photo, `
SELECT
	id, uuid, filename, title_cn, title_en, description, category, shot_time, width, height, orientation, resolution,
	camera_model, lens_model, aperture, shutter_speed, iso, focal_length, focal_length_35mm, metering_mode,
	exposure_program, white_balance, flash, thumb_url, display_url, original_url, like_count, download_count,
	view_count, is_published, created_at, updated_at
FROM photos
WHERE uuid = $1 AND is_published = TRUE
LIMIT 1
`, uuid)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}

func (r *SQLXPhotoRepository) IncrementViewCount(ctx context.Context, uuid string) error {
	return r.incrementCounter(ctx, uuid, "view_count")
}

func (r *SQLXPhotoRepository) IncrementDownloadCount(ctx context.Context, uuid string) error {
	return r.incrementCounter(ctx, uuid, "download_count")
}

func (r *SQLXPhotoRepository) AddLike(ctx context.Context, uuid, visitorHash string) error {
	if r.db == nil {
		return ErrRepositoryNotReady
	}

	ok, err := r.tableExists(ctx, "photos")
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotImplemented
	}

	ok, err = r.tableExists(ctx, "photo_like")
	if err != nil {
		return err
	}
	if !ok {
		// TODO: schema will be finalized in docs/sql/init.sql and migrations.
		return ErrNotImplemented
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var photoID int64
	if err := tx.GetContext(ctx, &photoID, `SELECT id FROM photos WHERE uuid = $1 LIMIT 1`, uuid); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO photo_like (photo_id, visitor_hash, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT DO NOTHING
`, photoID, visitorHash)
	if err != nil {
		return fmt.Errorf("insert photo_like failed: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
UPDATE photos
SET like_count = (SELECT COUNT(1) FROM photo_like WHERE photo_id = $1),
	updated_at = NOW()
WHERE id = $1
`, photoID)
	if err != nil {
		return fmt.Errorf("update like_count failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx failed: %w", err)
	}

	return nil
}

func (r *SQLXPhotoRepository) incrementCounter(ctx context.Context, uuid, field string) error {
	if r.db == nil {
		return ErrRepositoryNotReady
	}

	ok, err := r.tableExists(ctx, "photos")
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotImplemented
	}

	query := fmt.Sprintf(`UPDATE photos SET %s = %s + 1, updated_at = NOW() WHERE uuid = $1`, field, field)
	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return fmt.Errorf("increment %s failed: %w", field, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read rows affected failed: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLXPhotoRepository) tableExists(ctx context.Context, table string) (bool, error) {
	var tableName sql.NullString
	if err := r.db.GetContext(ctx, &tableName, `SELECT to_regclass($1)`, fmt.Sprintf("public.%s", table)); err != nil {
		return false, fmt.Errorf("check table %s failed: %w", table, err)
	}
	return tableName.Valid, nil
}
