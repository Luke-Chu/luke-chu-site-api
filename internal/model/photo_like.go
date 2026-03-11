package model

import "time"

type PhotoLike struct {
	ID          int64     `db:"id"`
	PhotoID     int64     `db:"photo_id"`
	VisitorHash string    `db:"visitor_hash"`
	CreatedAt   time.Time `db:"created_at"`
}
