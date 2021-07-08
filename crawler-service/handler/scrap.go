package handler

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type ReturnList struct {
	List  []string
	Key   string
	Error error
}

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link, accept, user_agent string, depth int) ReturnList {
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
	c.UserAgent = user_agent

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		//resolve url if item domain is not the same as base domain, link is unrelated. skip
		item_domain := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
		if item_domain != base_domain {
			//do nothing
		} else {
			links = append(links, r.URL.String())
		}
	})

	err = c.Visit(link)
	list, key := Builder(removeDuplicateValues(links), accept, user_agent)
	resolved_list := resolveUrl(link, list)
	return ReturnList{
		List:  resolved_list,
		Key:   key,
		Error: err,
	}
}

func resolveUrl(base_url string, urls []string) []string {
	base, _ := url.Parse(base_url)
	return_slice := []string{}
	for _, val := range urls {
		u, err := url.Parse(val)
		if err != nil {
			fmt.Println("error resolving url ", err)
			continue
		}
		return_slice = append(return_slice, base.ResolveReference(u).String())
	}
	return return_slice
}

//remove duplocate item from slice
func removeDuplicateValues(UrlSlice []string) []string {
	fmt.Println("resolving url...")
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
