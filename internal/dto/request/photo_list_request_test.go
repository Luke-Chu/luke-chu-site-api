package request

import (
	"reflect"
	"testing"
)

func TestPhotoListRequestNormalize(t *testing.T) {
	req := &PhotoListRequest{
		Page:        0,
		PageSize:    100,
		Sort:        "unknown",
		Order:       "bad",
		TagMode:     "",
		Orientation: "bad",
		Year:        3000,
		Month:       13,
	}

	req.Normalize()

	if req.Page != 1 {
		t.Fatalf("page = %d, want 1", req.Page)
	}
	if req.PageSize != 60 {
		t.Fatalf("pageSize = %d, want 60", req.PageSize)
	}
	if req.Sort != "shot_time" {
		t.Fatalf("sort = %s, want shot_time", req.Sort)
	}
	if req.Order != "desc" {
		t.Fatalf("order = %s, want desc", req.Order)
	}
	if req.TagMode != "any" {
		t.Fatalf("tagMode = %s, want any", req.TagMode)
	}
	if req.Orientation != "" {
		t.Fatalf("orientation = %s, want empty", req.Orientation)
	}
	if req.Year != 0 || req.Month != 0 {
		t.Fatalf("year/month = %d/%d, want 0/0", req.Year, req.Month)
	}
}

func TestPhotoListRequestTagList(t *testing.T) {
	req := PhotoListRequest{Tags: "风光, 夕阳，风光、 海边"}
	got := req.TagList()
	want := []string{"风光", "夕阳", "海边"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TagList() = %#v, want %#v", got, want)
	}
}
