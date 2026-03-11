package model

type PhotoTag struct {
	PhotoID int64 `db:"photo_id"`
	TagID   int64 `db:"tag_id"`
}
