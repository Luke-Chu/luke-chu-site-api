package response

type OrientationOption struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type FilterData struct {
	Years        []int                `json:"years"`
	Categories   []string             `json:"categories"`
	Orientations []OrientationOption  `json:"orientations"`
	TagTypes     []string             `json:"tagTypes"`
	Tags         map[string][]TagItem `json:"tags"`
}
