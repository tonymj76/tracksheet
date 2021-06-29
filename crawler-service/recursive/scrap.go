package recursive

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gocolly/colly"
)

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link string, depth int) []*url.URL {
	links := []*url.URL{}
	base, err := url.Parse(link)
	if err != nil {
		log.Fatal("Error parsing URL ", err)
	}
	//get base domain
	base_domain := fmt.Sprintf("%s://%s", base.Scheme, base.Host)

	c := colly.NewCollector(
		colly.MaxDepth(depth),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		//resolve url:
		//if item domain is not the same as base domain, link is unrelated. skip
		item_domain := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
		if item_domain != base_domain {
			//do nothing
		} else {
			links = append(links, r.URL)

		}
	})

	c.Visit(link)
	unique_url := removeDuplicateValues(links)
	return unique_url
}

//remove duplocate item from slice
func removeDuplicateValues(UrlSlice []*url.URL) []*url.URL {
	keys := make(map[*url.URL]bool)
	list := []*url.URL{}

	for _, entry := range UrlSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
