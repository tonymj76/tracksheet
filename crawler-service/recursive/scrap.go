package recursive

import (
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

var keywords = []string{
	"category", "cate", "sungless", "phone", "bag", "cloth", "shirt",
	"causal", "product", "pants", "fashion", "man", "woman", "watches",
	"women", "socks",
}

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link string, depth int) []*url.URL {
	validLink, err := url.Parse(link)
	if err != nil {
		logrus.WithError(err).Fatalf("%s is not a valid url link", link)
	}
	links := []*url.URL{}
	c := colly.NewCollector(
		colly.MaxDepth(depth),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if find(e.Attr("href"), keywords) {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		if r.URL.Host == validLink.Host {
			links = append(links, r.URL)
		}
	})

	c.Visit(link)
	return links
}

// find takes a slice and looks for an element in it. If found it will
// return boolean
func find(val string, slice []string) bool {
	for _, item := range slice {
		if strings.Contains(val, item) {
			return true
		}
	}
	return false
}
