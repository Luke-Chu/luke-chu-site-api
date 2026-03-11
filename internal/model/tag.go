package model

import "time"

type Tag struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	TagType   string    `db:"tag_type"`
	CreatedAt time.Time `db:"created_at"`
}
