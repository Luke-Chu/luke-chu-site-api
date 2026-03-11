package response

type PhotoViewData struct {
	UUID      string `json:"uuid"`
	ViewCount int64  `json:"viewCount"`
	Counted   bool   `json:"counted"`
}

type PhotoLikeData struct {
	UUID      string `json:"uuid"`
	Liked     bool   `json:"liked"`
	LikeCount int64  `json:"likeCount"`
}

type PhotoUnlikeData struct {
	UUID      string `json:"uuid"`
	Unliked   bool   `json:"unliked"`
	LikeCount int64  `json:"likeCount"`
}

type PhotoDownloadData struct {
	UUID          string `json:"uuid"`
	DownloadCount int64  `json:"downloadCount"`
	DownloadURL   string `json:"downloadUrl"`
}
