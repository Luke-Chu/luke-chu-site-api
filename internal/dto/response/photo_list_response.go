package response

type PhotoTagItem struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	TagType string `json:"tagType"`
}

type PhotoListItem struct {
	ID            int64          `json:"id"`
	UUID          string         `json:"uuid"`
	Filename      string         `json:"filename"`
	TitleCN       string         `json:"titleCn,omitempty"`
	TitleEN       string         `json:"titleEn,omitempty"`
	ThumbURL      string         `json:"thumbUrl,omitempty"`
	DisplayURL    string         `json:"displayUrl,omitempty"`
	Width         int            `json:"width"`
	Height        int            `json:"height"`
	Orientation   string         `json:"orientation,omitempty"`
	ShotTime      string         `json:"shotTime,omitempty"`
	Aperture      string         `json:"aperture,omitempty"`
	ShutterSpeed  string         `json:"shutterSpeed,omitempty"`
	ISO           int            `json:"iso,omitempty"`
	LikeCount     int64          `json:"likeCount"`
	ViewCount     int64          `json:"viewCount"`
	DownloadCount int64          `json:"downloadCount"`
	Tags          []PhotoTagItem `json:"tags"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type PhotoListQuery struct {
	Q           string   `json:"q"`
	Keywords    []string `json:"keywords"`
	Sort        string   `json:"sort"`
	Order       string   `json:"order"`
	Tags        []string `json:"tags"`
	TagMode     string   `json:"tagMode"`
	Orientation string   `json:"orientation,omitempty"`
	Year        int      `json:"year,omitempty"`
	Month       int      `json:"month,omitempty"`
	Category    string   `json:"category,omitempty"`
}

type PhotoListData struct {
	List       []PhotoListItem `json:"list"`
	Pagination Pagination      `json:"pagination"`
	Query      PhotoListQuery  `json:"query"`
}
