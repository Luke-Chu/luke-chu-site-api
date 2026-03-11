package search

import (
	"reflect"
	"testing"
)

func TestParseKeywords(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "空字符串",
			in:   "",
			want: []string{},
		},
		{
			name: "单个关键词",
			in:   "天空",
			want: []string{"天空"},
		},
		{
			name: "空格分隔",
			in:   "天空 风筝",
			want: []string{"天空", "风筝"},
		},
		{
			name: "英文逗号分隔",
			in:   "天空,风筝",
			want: []string{"天空", "风筝"},
		},
		{
			name: "中文逗号分隔",
			in:   "天空，风筝",
			want: []string{"天空", "风筝"},
		},
		{
			name: "顿号分隔",
			in:   "天空、风筝",
			want: []string{"天空", "风筝"},
		},
		{
			name: "连续分隔符",
			in:   "天空,,， 、风筝   海边",
			want: []string{"天空", "风筝", "海边"},
		},
		{
			name: "去重",
			in:   "天空,风筝,天空,风筝",
			want: []string{"天空", "风筝"},
		},
		{
			name: "最多保留5个关键词",
			in:   "a,b,c,d,e,f,g",
			want: []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseKeywords(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseKeywords() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
