package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type StatsRepository interface {
	Ping(ctx context.Context) error
}

type SQLXStatsRepository struct {
	db *sqlx.DB
}

func NewStatsRepository(db *sqlx.DB) StatsRepository {
	return &SQLXStatsRepository{db: db}
}

func (r *SQLXStatsRepository) Ping(ctx context.Context) error {
	if r.db == nil {
		return ErrRepositoryNotReady
	}
	if err := r.db.PingContext(ctx); err != nil {
		return fmt.Errorf("stats repository ping failed: %w", err)
	}
	return nil
}
