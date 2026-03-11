package response

type TagItem struct {
	Name    string `json:"name"`
	TagType string `json:"tagType"`
}

type TagListData struct {
	Items []TagItem `json:"items"`
}
