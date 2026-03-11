package ossutil

import "testing"

func TestInferRegionFromEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{name: "normal endpoint", endpoint: "oss-cn-hongkong.aliyuncs.com", want: "cn-hongkong"},
		{name: "empty endpoint", endpoint: "", want: ""},
		{name: "custom host", endpoint: "static.example.com", want: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := inferRegionFromEndpoint(tt.endpoint)
			if got != tt.want {
				t.Fatalf("inferRegionFromEndpoint() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestObjectKeyFromSourceURL(t *testing.T) {
	t.Parallel()

	signer := &PresignDownloadURLSigner{
		publicBaseHost: "luke-chu-site-photography.oss-cn-hongkong.aliyuncs.com",
		publicBasePath: "",
	}

	tests := []struct {
		name      string
		sourceURL string
		want      string
	}{
		{
			name:      "public url",
			sourceURL: "https://luke-chu-site-photography.oss-cn-hongkong.aliyuncs.com/photos/2026/03/a.jpg",
			want:      "photos/2026/03/a.jpg",
		},
		{
			name:      "raw key",
			sourceURL: "photos/2026/03/a.jpg",
			want:      "photos/2026/03/a.jpg",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := signer.objectKeyFromSourceURL(tt.sourceURL)
			if err != nil {
				t.Fatalf("objectKeyFromSourceURL() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("objectKeyFromSourceURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
