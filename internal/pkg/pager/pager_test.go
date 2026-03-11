package pager

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{
			name:         "默认值",
			page:         0,
			pageSize:     0,
			wantPage:     1,
			wantPageSize: 30,
		},
		{
			name:         "page下限",
			page:         -2,
			pageSize:     20,
			wantPage:     1,
			wantPageSize: 20,
		},
		{
			name:         "pageSize上限",
			page:         2,
			pageSize:     100,
			wantPage:     2,
			wantPageSize: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotSize := Normalize(tt.page, tt.pageSize)
			if gotPage != tt.wantPage || gotSize != tt.wantPageSize {
				t.Fatalf("Normalize() = (%d,%d), want (%d,%d)", gotPage, gotSize, tt.wantPage, tt.wantPageSize)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	got := Offset(3, 20)
	want := 40
	if got != want {
		t.Fatalf("Offset() = %d, want %d", got, want)
	}
}
