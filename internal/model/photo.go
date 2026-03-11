package model

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID              int64      `db:"id"`
	UUID            uuid.UUID  `db:"uuid"`
	Filename        string     `db:"filename"`
	TitleCN         *string    `db:"title_cn"`
	TitleEN         *string    `db:"title_en"`
	Description     *string    `db:"description"`
	Category        *string    `db:"category"`
	ShotTime        *time.Time `db:"shot_time"`
	Width           int        `db:"width"`
	Height          int        `db:"height"`
	Orientation     string     `db:"orientation"`
	Resolution      *string    `db:"resolution"`
	CameraModel     *string    `db:"camera_model"`
	LensModel       *string    `db:"lens_model"`
	Aperture        *string    `db:"aperture"`
	ShutterSpeed    *string    `db:"shutter_speed"`
	ISO             *int       `db:"iso"`
	FocalLength     *string    `db:"focal_length"`
	FocalLength35mm *string    `db:"focal_length_35mm"`
	MeteringMode    *string    `db:"metering_mode"`
	ExposureProgram *string    `db:"exposure_program"`
	WhiteBalance    *string    `db:"white_balance"`
	Flash           *string    `db:"flash"`
	ThumbURL        *string    `db:"thumb_url"`
	DisplayURL      *string    `db:"display_url"`
	OriginalURL     *string    `db:"original_url"`
	LikeCount       int64      `db:"like_count"`
	DownloadCount   int64      `db:"download_count"`
	ViewCount       int64      `db:"view_count"`
	IsPublished     bool       `db:"is_published"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}
