package response

type PhotoListItem struct {
	UUID          string `json:"uuid"`
	TitleCN       string `json:"titleCn,omitempty"`
	TitleEN       string `json:"titleEn,omitempty"`
	Orientation   string `json:"orientation,omitempty"`
	ThumbURL      string `json:"thumbUrl,omitempty"`
	DisplayURL    string `json:"displayUrl,omitempty"`
	LikeCount     int64  `json:"likeCount"`
	ViewCount     int64  `json:"viewCount"`
	DownloadCount int64  `json:"downloadCount"`
	ShotTime      string `json:"shotTime,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type PhotoListData struct {
	Items      []PhotoListItem `json:"items"`
	Pagination Pagination      `json:"pagination"`
}
