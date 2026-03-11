package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"luke-chu-site-api/internal/constant"
	"luke-chu-site-api/internal/dto/request"
	"luke-chu-site-api/internal/dto/response"
	"luke-chu-site-api/internal/model"
	"luke-chu-site-api/internal/pkg/pager"
	sortutil "luke-chu-site-api/internal/pkg/sort"
)

var (
	ErrRepositoryNotReady = errors.New("repository not ready")
	ErrNotImplemented     = errors.New("repository sql not implemented for current schema")
)

var photoSortColumns = map[string]string{
	constant.SortShotTime:  "p.shot_time",
	constant.SortLikeCount: "p.like_count",
	constant.SortViewCount: "p.view_count",
	constant.SortDownload:  "p.download_count",
	constant.SortCreatedAt: "p.created_at",
}

type PhotoRepository interface {
	ListPhotos(ctx context.Context, req *request.PhotoListRequest) ([]*model.Photo, error)
	CountPhotos(ctx context.Context, req *request.PhotoListRequest) (int64, error)
	ListPhotoTagsByPhotoIDs(ctx context.Context, photoIDs []int64) (map[int64][]response.PhotoTagItem, error)
	GetPhotoDetailByUUID(ctx context.Context, uuid string) (*model.Photo, error)
	GetPublishedPhotoBaseByUUID(ctx context.Context, uuid string) (*model.Photo, error)
	GetPhotoTagsByPhotoID(ctx context.Context, photoID int64) ([]response.TagItem, error)
	GetPhotoByUUID(ctx context.Context, uuid string) (*model.Photo, error)
	IncrementViewCount(ctx context.Context, uuid, visitorHash string) (int64, bool, error)
	IncrementDownloadCount(ctx context.Context, uuid, visitorHash string) (int64, string, bool, error)
	AddLike(ctx context.Context, uuid, visitorHash string) (bool, int64, error)
	RemoveLike(ctx context.Context, uuid, visitorHash string) (bool, int64, error)
}

type SQLXPhotoRepository struct {
	db *sqlx.DB
}

func NewPhotoRepository(db *sqlx.DB) PhotoRepository {
	return &SQLXPhotoRepository{db: db}
}

func (r *SQLXPhotoRepository) ListPhotos(ctx context.Context, req *request.PhotoListRequest) ([]*model.Photo, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}
	if req == nil {
		req = &request.PhotoListRequest{}
	}
	req.Normalize()

	whereSQL, args, next := r.buildPhotoListWhere(req, 1)
	sortColumn := normalizePhotoSortColumn(req.Sort)
	sortOrder := strings.ToUpper(sortutil.NormalizeSortOrder(req.Order))

	query := fmt.Sprintf(`
SELECT
	p.id, p.uuid, p.filename, p.title_cn, p.title_en, p.thumb_url, p.display_url,
	p.width, p.height, p.orientation, p.shot_time, p.aperture, p.shutter_speed,
	p.iso, p.like_count, p.view_count, p.download_count
FROM photos p
%s
ORDER BY %s %s, p.id DESC
LIMIT $%d OFFSET $%d
`, whereSQL, sortColumn, sortOrder, next, next+1)

	args = append(args, req.PageSize, pager.Offset(req.Page, req.PageSize))

	rows := make([]model.Photo, 0)
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("list photos failed: %w", err)
	}

	result := make([]*model.Photo, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

func (r *SQLXPhotoRepository) CountPhotos(ctx context.Context, req *request.PhotoListRequest) (int64, error) {
	if r.db == nil {
		return 0, ErrRepositoryNotReady
	}
	if req == nil {
		req = &request.PhotoListRequest{}
	}
	req.Normalize()

	whereSQL, args, _ := r.buildPhotoListWhere(req, 1)
	query := fmt.Sprintf(`SELECT COUNT(1) FROM photos p %s`, whereSQL)

	var total int64
	if err := r.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, fmt.Errorf("count photos failed: %w", err)
	}

	return total, nil
}

func (r *SQLXPhotoRepository) ListPhotoTagsByPhotoIDs(ctx context.Context, photoIDs []int64) (map[int64][]response.PhotoTagItem, error) {
	result := make(map[int64][]response.PhotoTagItem, len(photoIDs))
	for _, photoID := range photoIDs {
		result[photoID] = []response.PhotoTagItem{}
	}
	if len(photoIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, 0, len(photoIDs))
	args := make([]any, 0, len(photoIDs))
	for i, photoID := range photoIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, photoID)
	}

	query := fmt.Sprintf(`
SELECT
	pt.photo_id,
	t.id,
	t.name,
	t.tag_type
FROM photo_tags pt
JOIN tags t ON t.id = pt.tag_id
WHERE pt.photo_id IN (%s)
ORDER BY pt.photo_id ASC, t.id ASC
`, strings.Join(placeholders, ","))

	type tagRow struct {
		PhotoID int64  `db:"photo_id"`
		TagID   int64  `db:"id"`
		Name    string `db:"name"`
		TagType string `db:"tag_type"`
	}

	rows := make([]tagRow, 0)
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("list photo tags failed: %w", err)
	}

	for _, row := range rows {
		result[row.PhotoID] = append(result[row.PhotoID], response.PhotoTagItem{
			ID:      row.TagID,
			Name:    row.Name,
			TagType: row.TagType,
		})
	}

	return result, nil
}

func (r *SQLXPhotoRepository) GetPhotoByUUID(ctx context.Context, uuid string) (*model.Photo, error) {
	return r.GetPhotoDetailByUUID(ctx, uuid)
}

func (r *SQLXPhotoRepository) GetPhotoDetailByUUID(ctx context.Context, uuid string) (*model.Photo, error) {
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
	id, uuid, filename, title_cn, title_en, description, category, shot_time, width, height, resolution, orientation,
	camera_model, lens_model, focal_length, focal_length_35mm, aperture, shutter_speed, iso, metering_mode,
	exposure_compensation, exposure_program, white_balance, flash, thumb_url, display_url, original_url, like_count, download_count,
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

func (r *SQLXPhotoRepository) GetPhotoTagsByPhotoID(ctx context.Context, photoID int64) ([]response.TagItem, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}

	type tagRow struct {
		ID      int64  `db:"id"`
		Name    string `db:"name"`
		TagType string `db:"tag_type"`
	}

	rows := make([]tagRow, 0)
	err := r.db.SelectContext(ctx, &rows, `
SELECT t.id, t.name, t.tag_type
FROM photo_tags pt
JOIN tags t ON t.id = pt.tag_id
WHERE pt.photo_id = $1
ORDER BY t.tag_type ASC, t.id ASC
`, photoID)
	if err != nil {
		return nil, fmt.Errorf("get photo tags failed: %w", err)
	}

	result := make([]response.TagItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, response.TagItem{
			ID:      row.ID,
			Name:    row.Name,
			TagType: row.TagType,
		})
	}
	return result, nil
}

func (r *SQLXPhotoRepository) GetPublishedPhotoBaseByUUID(ctx context.Context, uuid string) (*model.Photo, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}

	var photo model.Photo
	err := r.db.GetContext(ctx, &photo, `
SELECT id, uuid, like_count
FROM photos
WHERE uuid = $1
  AND is_published = TRUE
LIMIT 1
`, uuid)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *SQLXPhotoRepository) IncrementViewCount(ctx context.Context, uuid, visitorHash string) (int64, bool, error) {
	if r.db == nil {
		return 0, false, ErrRepositoryNotReady
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, false, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var photoID int64
	if err := tx.GetContext(ctx, &photoID, `
SELECT id
FROM photos
WHERE uuid = $1
  AND is_published = TRUE
LIMIT 1
`, uuid); err != nil {
		return 0, false, err
	}

	var inWindow bool
	if err := tx.GetContext(ctx, &inWindow, `
SELECT EXISTS (
	SELECT 1
	FROM photo_views
	WHERE photo_id = $1
	  AND visitor_hash = $2
	  AND viewed_at >= NOW() - INTERVAL '10 minutes'
)
`, photoID, visitorHash); err != nil {
		return 0, false, fmt.Errorf("check photo_views window failed: %w", err)
	}

	var count int64
	if inWindow {
		if err := tx.GetContext(ctx, &count, `SELECT view_count FROM photos WHERE id = $1`, photoID); err != nil {
			return 0, false, fmt.Errorf("query view_count failed: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return 0, false, fmt.Errorf("commit tx failed: %w", err)
		}
		return count, false, nil
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO photo_views (photo_id, visitor_hash, viewed_at)
VALUES ($1, $2, NOW())
`, photoID, visitorHash); err != nil {
		return 0, false, fmt.Errorf("insert photo_views failed: %w", err)
	}

	if err := tx.GetContext(ctx, &count, `
UPDATE photos
SET view_count = view_count + 1,
	updated_at = NOW()
WHERE id = $1
RETURNING view_count
`, photoID); err != nil {
		return 0, false, err
	}

	if err := tx.Commit(); err != nil {
		return 0, false, fmt.Errorf("commit tx failed: %w", err)
	}

	return count, true, nil
}

func (r *SQLXPhotoRepository) IncrementDownloadCount(ctx context.Context, uuid, visitorHash string) (int64, string, bool, error) {
	if r.db == nil {
		return 0, "", false, ErrRepositoryNotReady
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, "", false, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	type baseRow struct {
		PhotoID int64          `db:"id"`
		Count   int64          `db:"download_count"`
		URL     sql.NullString `db:"original_url"`
	}
	var base baseRow
	err = tx.GetContext(ctx, &base, `
SELECT id, download_count, original_url
FROM photos
WHERE uuid = $1
  AND is_published = TRUE
LIMIT 1
`, uuid)
	if err != nil {
		return 0, "", false, err
	}

	var inWindow bool
	if err := tx.GetContext(ctx, &inWindow, `
SELECT EXISTS (
	SELECT 1
	FROM photo_downloads
	WHERE photo_id = $1
	  AND visitor_hash = $2
	  AND downloaded_at >= NOW() - INTERVAL '30 minutes'
)
`, base.PhotoID, visitorHash); err != nil {
		return 0, "", false, fmt.Errorf("check photo_downloads window failed: %w", err)
	}

	if inWindow {
		if err := tx.Commit(); err != nil {
			return 0, "", false, fmt.Errorf("commit tx failed: %w", err)
		}
		return base.Count, base.URL.String, false, nil
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO photo_downloads (photo_id, visitor_hash, downloaded_at)
VALUES ($1, $2, NOW())
`, base.PhotoID, visitorHash); err != nil {
		return 0, "", false, fmt.Errorf("insert photo_downloads failed: %w", err)
	}

	var count int64
	if err := tx.GetContext(ctx, &count, `
UPDATE photos
SET download_count = download_count + 1,
	updated_at = NOW()
WHERE id = $1
RETURNING download_count
`, base.PhotoID); err != nil {
		return 0, "", false, err
	}

	if err := tx.Commit(); err != nil {
		return 0, "", false, fmt.Errorf("commit tx failed: %w", err)
	}

	return count, base.URL.String, true, nil
}

func (r *SQLXPhotoRepository) AddLike(ctx context.Context, uuid, visitorHash string) (bool, int64, error) {
	if r.db == nil {
		return false, 0, ErrRepositoryNotReady
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return false, 0, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var photoID int64
	if err := tx.GetContext(ctx, &photoID, `
SELECT id
FROM photos
WHERE uuid = $1
  AND is_published = TRUE
LIMIT 1
`, uuid); err != nil {
		return false, 0, err
	}

	result, err := tx.ExecContext(ctx, `
INSERT INTO photo_likes (photo_id, visitor_hash, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT DO NOTHING
`, photoID, visitorHash)
	if err != nil {
		return false, 0, fmt.Errorf("insert photo_likes failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, 0, fmt.Errorf("check like insert rows failed: %w", err)
	}

	liked := rowsAffected > 0
	var likeCount int64
	if liked {
		if err := tx.GetContext(ctx, &likeCount, `
UPDATE photos
SET like_count = like_count + 1,
	updated_at = NOW()
WHERE id = $1
RETURNING like_count
`, photoID); err != nil {
			return false, 0, fmt.Errorf("increment like_count failed: %w", err)
		}
	} else {
		if err := tx.GetContext(ctx, &likeCount, `
SELECT like_count
FROM photos
WHERE id = $1
`, photoID); err != nil {
			return false, 0, fmt.Errorf("query like_count failed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, 0, fmt.Errorf("commit tx failed: %w", err)
	}

	return liked, likeCount, nil
}

func (r *SQLXPhotoRepository) RemoveLike(ctx context.Context, uuid, visitorHash string) (bool, int64, error) {
	if r.db == nil {
		return false, 0, ErrRepositoryNotReady
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return false, 0, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var photoID int64
	if err := tx.GetContext(ctx, &photoID, `
SELECT id
FROM photos
WHERE uuid = $1
  AND is_published = TRUE
LIMIT 1
`, uuid); err != nil {
		return false, 0, err
	}

	result, err := tx.ExecContext(ctx, `
DELETE FROM photo_likes
WHERE photo_id = $1
  AND visitor_hash = $2
`, photoID, visitorHash)
	if err != nil {
		return false, 0, fmt.Errorf("delete photo_likes failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, 0, fmt.Errorf("check unlike rows failed: %w", err)
	}

	unliked := rowsAffected > 0
	var likeCount int64
	if unliked {
		if err := tx.GetContext(ctx, &likeCount, `
UPDATE photos
SET like_count = GREATEST(like_count - 1, 0),
	updated_at = NOW()
WHERE id = $1
RETURNING like_count
`, photoID); err != nil {
			return false, 0, fmt.Errorf("decrement like_count failed: %w", err)
		}
	} else {
		if err := tx.GetContext(ctx, &likeCount, `
SELECT like_count
FROM photos
WHERE id = $1
`, photoID); err != nil {
			return false, 0, fmt.Errorf("query like_count failed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return false, 0, fmt.Errorf("commit tx failed: %w", err)
	}

	return unliked, likeCount, nil
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

func (r *SQLXPhotoRepository) buildPhotoListWhere(req *request.PhotoListRequest, startIndex int) (string, []any, int) {
	clauses := []string{"p.is_published = TRUE"}
	args := make([]any, 0, 16)
	idx := startIndex

	for _, keyword := range req.KeywordList() {
		keywordLike := "%" + keyword + "%"
		clauses = append(clauses, fmt.Sprintf(`(
COALESCE(p.title_cn, '') ILIKE $%d
OR COALESCE(p.title_en, '') ILIKE $%d
OR p.filename ILIKE $%d
OR EXISTS (
	SELECT 1
	FROM photo_tags pt
	JOIN tags t ON t.id = pt.tag_id
	WHERE pt.photo_id = p.id
	  AND t.name ILIKE $%d
)
)`, idx, idx, idx, idx))
		args = append(args, keywordLike)
		idx++
	}

	tags := req.TagList()
	if len(tags) > 0 {
		tagNames := make([]string, 0, len(tags))
		seen := make(map[string]struct{}, len(tags))
		for _, tag := range tags {
			lower := strings.ToLower(strings.TrimSpace(tag))
			if lower == "" {
				continue
			}
			if _, ok := seen[lower]; ok {
				continue
			}
			seen[lower] = struct{}{}
			tagNames = append(tagNames, lower)
		}

		if len(tagNames) > 0 {
			placeholders := make([]string, 0, len(tagNames))
			for _, tagName := range tagNames {
				placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
				args = append(args, tagName)
				idx++
			}

			switch req.TagMode {
			case "all":
				clauses = append(clauses, fmt.Sprintf(`
p.id IN (
	SELECT pt.photo_id
	FROM photo_tags pt
	JOIN tags t ON t.id = pt.tag_id
	WHERE LOWER(t.name) IN (%s)
	GROUP BY pt.photo_id
	HAVING COUNT(DISTINCT LOWER(t.name)) = $%d
)`, strings.Join(placeholders, ","), idx))
				args = append(args, len(tagNames))
				idx++
			default:
				clauses = append(clauses, fmt.Sprintf(`
EXISTS (
	SELECT 1
	FROM photo_tags pt
	JOIN tags t ON t.id = pt.tag_id
	WHERE pt.photo_id = p.id
	  AND LOWER(t.name) IN (%s)
)`, strings.Join(placeholders, ",")))
			}
		}
	}

	if req.Orientation != "" {
		clauses = append(clauses, fmt.Sprintf("p.orientation = $%d", idx))
		args = append(args, req.Orientation)
		idx++
	}
	if req.Year > 0 {
		clauses = append(clauses, fmt.Sprintf("p.year = $%d", idx))
		args = append(args, req.Year)
		idx++
	}
	if req.Month > 0 {
		clauses = append(clauses, fmt.Sprintf("p.month = $%d", idx))
		args = append(args, req.Month)
		idx++
	}
	if req.Category != "" {
		clauses = append(clauses, fmt.Sprintf("p.category = $%d", idx))
		args = append(args, req.Category)
		idx++
	}

	return "WHERE " + strings.Join(clauses, " AND "), args, idx
}

func normalizePhotoSortColumn(sortField string) string {
	field := sortutil.NormalizeSortField(sortField)
	col, ok := photoSortColumns[field]
	if !ok {
		return photoSortColumns[constant.DefaultSort]
	}
	return col
}

func (r *SQLXPhotoRepository) tableExists(ctx context.Context, table string) (bool, error) {
	var tableName sql.NullString
	if err := r.db.GetContext(ctx, &tableName, `SELECT to_regclass($1)`, fmt.Sprintf("public.%s", table)); err != nil {
		return false, fmt.Errorf("check table %s failed: %w", table, err)
	}
	return tableName.Valid, nil
}
