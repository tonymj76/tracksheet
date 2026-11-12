package log

import "testing"

func TestFormatFilePath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/go/src/cache_service/crawler-service/handler/handler.go", "handler.go"},
		{"relative/dir/file.go", "file.go"},
		{"file.go", "file.go"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := formatFilePath(tt.path); got != tt.want {
			t.Errorf("formatFilePath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
