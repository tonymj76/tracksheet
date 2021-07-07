package handler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

//build a hash map of parent node mapping to a slice of all related links
func Builder(urls []string, accept, user_agent string) ([]string, string) {
	m := make(map[string][]string)
	return_slice := []string{}
	return_key := ""

	for _, val := range urls {
		base, _ := url.Parse(val)
		res, err := GoClient(val, NewConnect(logrus.New(), val, accept, user_agent))
		if err != nil {
			fmt.Println("error getting response ", err)
		} else {
			// parsse web page and get all links
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				fmt.Println("error parsing html ", err)
			}
			//get all links in page and group
			doc.Find("a").Each(func(i int, link *goquery.Selection) {
				href, ok := link.Attr("href")
				//resolve relative url to absolute url
				u, _ := url.Parse(href)
				if ok {
					class, ok := link.ParentsUntil("~").Attr("class")
					if ok {
						class_details := strings.ReplaceAll(class, " ", ",")
						if m[class_details] == nil {
							m[class_details] = []string{base.ResolveReference(u).String()}
						} else {
							links := m[class_details]
							links = append(links, base.ResolveReference(u).String())
							m[class_details] = links
						}

					}
				}
			})

		}
	}
	//return the group with the longest item
	for key, value := range m {
		if len(value) > len(return_slice) {
			return_slice = value
			return_key = key
		}
	}
	return return_slice, return_key
}
