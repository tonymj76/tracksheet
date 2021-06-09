package structure

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/headzoo/surf.v1"
)

var robotURL = "%s://%s/robots.txt"

type (
	Product struct {
		XMLName xml.Name `xml:"urlset"`
		URL     []Feed   `xml:"url"`
	}

	Feed struct {
		XMLName xml.Name `xml:"url"`
		Loc     string   `xml:"loc"`
	}
	//Result of the search
	// Result struct {
	// 	URL     string
	// 	Content string
	// }
	// ProductDescription use for searching
	ProductDescription struct{}
)

// //Search _
// func (prodDesc *ProductDescription) Search(urlString string) (result []*Result, err error) {
// 	parseURL, err := url.Parse(urlString)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	c := colly.NewCollector()
// 	c.OnHTML("body", func(el *colly.HTMLElement) {
// 		fmt.Println(el.DOM)
// 		fmt.Println(el.Text)
// 	})

// 	err = c.Visit(fmt.Sprintf(robotURL, parseURL.Scheme, parseURL.Hostname()))
// 	return
// }
// Search
func (prodDesc *ProductDescription) Search(urlString string) (result []string, err error) {
	parseURL, err := url.Parse(urlString)
	if err != nil {
		log.Println(err)
		return
	}
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")
	if err = bow.Open(fmt.Sprintf(robotURL, parseURL.Scheme, parseURL.Hostname())); err != nil {
		log.Fatalln(err)
		return
	}
	siteMaps := getSiteMapURL(bow.Body())
	fmt.Println(siteMaps)
	result, err = prodDesc.retrieve(siteMaps[1])
	return
}

func getSiteMapURL(matchString string) (values []string) {
	validString := regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>.*)`)
	for _, siteMapUrl := range validString.FindAllString(matchString, -1) {
		if strings.Contains(strings.ToLower(siteMapUrl), "sitemap") {
			values = append(values, siteMapUrl)
		}
	}
	return
}

func parseSiteMapURL(matchString string) (vs []string) {
	validString := regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^<">]+)`)
	vs = validString.FindAllString(matchString, -1)
	return
}

// retrieve performs a HTTP Get request for the sitemap xml and decodes the results
func (prodDesc *ProductDescription) retrieve(siteMap string) ([]string, error) {
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")
	if err := bow.Open(siteMap); err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return parseSiteMapURL(bow.Body()), nil
}
