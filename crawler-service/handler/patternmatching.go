package handler

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type SampleLink struct {
	Str1, Str2, Str3 string
}

func getallurls(link, accept, user_agent string, depth int, noAsync bool) ([]string, error) {
	links := []string{}
	base, err := url.Parse(link)
	if err != nil {
		logrus.WithError(err).Fatalln("Error parsing URL")
	}

	c := colly.NewCollector(
		colly.MaxDepth(depth),
		colly.Async(!noAsync),
		colly.UserAgent(user_agent),
		colly.URLFilters(
			//resolve url if item domain is not the same as base domain, link is unrelated. skip
			regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^\s]*)`, base.Hostname())),
		),
		// Cache responses to prevent multiple download of pages
		colly.CacheDir("./cache"),
	)
	c.IgnoreRobotsTxt = true
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 6, Delay: 10 * time.Second})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("error encoutered ", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
		links = append(links, r.URL.String())
	})

	err = c.Visit(link)

	c.Wait()

	list := removeDuplicateValues(links)
	return list, err
}

func stringCompare(str1, str2, str3 string) []string {
	array1 := regexp.MustCompile(`[\/\&\:\=\?\s]+`).Split(str1, -1)
	array2 := regexp.MustCompile(`[\/\&\:\=\?\s]+`).Split(str2, -1)
	array3 := regexp.MustCompile(`[\/\&\:\=\?\s]+`).Split(str3, -1)

	common_item := []string{}

	for _, val1 := range array1 {
		for _, val2 := range array2 {
			for _, val3 := range array3 {
				if val1 == val2 && val2 == val3 {
					common_item = append(common_item, val2)
				}
			}
		}
	}
	return common_item
}

func GetProductUrl(link, accept, user_agent string, depth int, sample SampleLink, noAsync bool) ReturnList {
	list, err := getallurls(link, accept, user_agent, depth, noAsync)
	common := stringCompare(sample.Str1, sample.Str2, sample.Str3)
	fmt.Println(common)
	final := []string{}
	fmt.Println("filtering product link...")
	for _, val := range list {
		num := 0
		for _, val2 := range common {
			if strings.Contains(val, val2) {
				num++
				if num == len(common) {
					final = append(final, val)
				}
			}  
		}
	}
	return ReturnList{
		List:  final,
		Key:   "",
		Error: err,
	}
}
