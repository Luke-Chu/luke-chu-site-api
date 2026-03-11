package response

type PhotoDetailData struct {
	ID              int64     `json:"id"`
	UUID            string    `json:"uuid"`
	Filename        string    `json:"filename,omitempty"`
	TitleCN         string    `json:"titleCn,omitempty"`
	TitleEN         string    `json:"titleEn,omitempty"`
	Description     string    `json:"description,omitempty"`
	Category        string    `json:"category,omitempty"`
	ShotTime        string    `json:"shotTime,omitempty"`
	Width           int       `json:"width"`
	Height          int       `json:"height"`
	Orientation     string    `json:"orientation,omitempty"`
	Resolution      string    `json:"resolution,omitempty"`
	CameraModel     string    `json:"cameraModel,omitempty"`
	LensModel       string    `json:"lensModel,omitempty"`
	Aperture        string    `json:"aperture,omitempty"`
	ShutterSpeed    string    `json:"shutterSpeed,omitempty"`
	ISO             int       `json:"iso,omitempty"`
	FocalLength     float64   `json:"focalLength,omitempty"`
	FocalLength35mm float64   `json:"focalLength35mm,omitempty"`
	MeteringMode    string    `json:"meteringMode,omitempty"`
	ExposureComp    string    `json:"exposureCompensation,omitempty"`
	ExposureProgram string    `json:"exposureProgram,omitempty"`
	WhiteBalance    string    `json:"whiteBalance,omitempty"`
	Flash           string    `json:"flash,omitempty"`
	ThumbURL        string    `json:"thumbUrl,omitempty"`
	DisplayURL      string    `json:"displayUrl,omitempty"`
	OriginalURL     string    `json:"originalUrl,omitempty"`
	LikeCount       int64     `json:"likeCount"`
	DownloadCount   int64     `json:"downloadCount"`
	ViewCount       int64     `json:"viewCount"`
	CreatedAt       string    `json:"createdAt,omitempty"`
	UpdatedAt       string    `json:"updatedAt,omitempty"`
	Tags            []TagItem `json:"tags"`
}
