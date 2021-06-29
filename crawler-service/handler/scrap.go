package handler

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link string, depth int) ([]string, error) {
	links := []string{}
	base, err := url.Parse(link)
	if err != nil {
		logrus.WithError(err).Fatalln("Error parsing URL")
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
			links = append(links, r.URL.String())
		}
	})

	err = c.Visit(link)
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
