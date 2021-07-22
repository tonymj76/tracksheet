package handler

import (
	"cache_service/crawler-service/config"
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

func getallurls(link string, depth int) ([]string, error) {
	links := []string{}
	base, err := url.Parse(link)
	if err != nil {
		logrus.WithError(err).Fatalln("Error parsing URL")
	}

	c := colly.NewCollector(
		colly.MaxDepth(depth),
		colly.Async(true),
		colly.UserAgent(config.GetEnv("userAgent")),
		colly.URLFilters(
			//resolve url if item domain is not the same as base domain, link is unrelated. skip
			regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^\s]*)`, base.Hostname())),
		),
		// Cache responses to prevent multiple download of pages
		colly.CacheDir("./cache"),
	)
	c.IgnoreRobotsTxt = true
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 6, RandomDelay: 10 * time.Second})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("error encoutered ", err)
	})
	c.OnRequest(func(r *colly.Request) {
		links = append(links, r.URL.String())
	})

	err = c.Visit(link)

	c.Wait()

	list := removeDuplicateValues(links)
	return list, err
}

//remove extra char from string and find common item
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

//filter out url based on pattern
func GetProductUrl(link string, depth int, sample SampleLink) ReturnList {
	list, err := getallurls(link, depth)
	common := stringCompare(sample.Str1, sample.Str2, sample.Str3)
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
		Error: err,
	}
}
