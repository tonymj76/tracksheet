package recursive

import (
	"net/url"

	"github.com/gocolly/colly"
)

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link string, depth int) []*url.URL {
	links := []*url.URL{}
	c := colly.NewCollector(
		colly.MaxDepth(depth),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		links = append(links, r.URL)
	})

	c.Visit(link)
	return links
}
