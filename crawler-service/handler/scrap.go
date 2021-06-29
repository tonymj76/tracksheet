package handler

import (
	"github.com/gocolly/colly"
)

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link string, depth int) ([]string, error) {
	links := []string{}
	c := colly.NewCollector(
		colly.MaxDepth(depth),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		links = append(links, r.URL.String())
	})

	err := c.Visit(link)
	return RelatedURL(removeDuplicateValues(links)), err
}

//remove duplocate item from slice
func removeDuplicateValues(UrlSlice []string) []string {
	keys := make(map[string]struct{})
	list := []string{}

	for _, entry := range UrlSlice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
