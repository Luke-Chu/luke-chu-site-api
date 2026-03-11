package response

type TagItem struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	TagType string `json:"tagType"`
}

type TagListData struct {
	Items []TagItem `json:"items"`
}
