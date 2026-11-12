package handler

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetSiteMapURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bare sitemap url",
			input: "https://example.com/sitemap.xml",
			want:  "https://example.com/sitemap.xml",
		},
		{
			name:  "nested sitemap path",
			input: "http://shop.example.com/maps/products.xml",
			want:  "http://shop.example.com/maps/products.xml",
		},
		{
			name:  "no sitemap present",
			input: "User-agent: *\nDisallow: /admin\n",
			want:  "",
		},
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSiteMapURL(tt.input); got != tt.want {
				t.Errorf("getSiteMapURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseSiteMapURL(t *testing.T) {
	body := `<urlset>
<url><loc>https://example.com/products/one</loc></url>
<url><loc>http://example.com/products/two</loc></url>
<url><loc>https://other.com/products/three</loc></url>
</urlset>`

	got := parseSiteMapURL(body, "example.com")
	want := []string{
		"https://example.com/products/one",
		"http://example.com/products/two",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseSiteMapURL() = %v, want %v", got, want)
	}

	if got := parseSiteMapURL(body, "missing.com"); len(got) != 0 {
		t.Errorf("parseSiteMapURL() for unknown host = %v, want no links", got)
	}
}

func TestNewConnect(t *testing.T) {
	logger := logrus.New()
	conn := NewConnect(logger, "https://example.com", "text/html", "test-agent")

	if conn.Logger != logger {
		t.Error("NewConnect did not keep the supplied logger")
	}
	if conn.urlLink != "https://example.com" {
		t.Errorf("urlLink = %q, want %q", conn.urlLink, "https://example.com")
	}
	if conn.accept != "text/html" {
		t.Errorf("accept = %q, want %q", conn.accept, "text/html")
	}
	if conn.userAgent != "test-agent" {
		t.Errorf("userAgent = %q, want %q", conn.userAgent, "test-agent")
	}
}

func TestMinInt(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{1, 2, 1},
		{2, 1, 1},
		{3, 3, 3},
		{-1, 0, -1},
	}
	for _, tt := range tests {
		if got := minInt(tt.a, tt.b); got != tt.want {
			t.Errorf("minInt(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
