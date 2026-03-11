package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"luke-chu-site-api/internal/model"
)

type TagRepository interface {
	ListTags(ctx context.Context) ([]model.Tag, error)
}

type SQLXTagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) TagRepository {
	return &SQLXTagRepository{db: db}
}

func (r *SQLXTagRepository) ListTags(ctx context.Context) ([]model.Tag, error) {
	if r.db == nil {
		return nil, ErrRepositoryNotReady
	}

	ok, err := r.tableExists(ctx, "tags")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotImplemented
	}

	tags := make([]model.Tag, 0)
	err = r.db.SelectContext(ctx, &tags, `
SELECT id, name, tag_type, created_at
FROM tags
ORDER BY name ASC
`)
	if err != nil {
		return nil, fmt.Errorf("list tags failed: %w", err)
	}

	return tags, nil
}

func (r *SQLXTagRepository) tableExists(ctx context.Context, table string) (bool, error) {
	var tableName sql.NullString
	if err := r.db.GetContext(ctx, &tableName, `SELECT to_regclass($1)`, fmt.Sprintf("public.%s", table)); err != nil {
		return false, fmt.Errorf("check table %s failed: %w", table, err)
	}
	return tableName.Valid, nil
}
