package structure

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/headzoo/surf.v1"
)

var robotURL = "%s://%s/robots.txt"

type ProductDescription struct{}

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
	siteMap := getSiteMapURL(bow.Body())
	result, err = prodDesc.retrieve(siteMap, parseURL.Hostname())
	return
}

func getSiteMapURL(matchString string) (value string) {
	//regex here map a url with a space or no space but must end with .xml
	validString := regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.xml$)|(?P<path2>[^\s]*\.xml\s+$)`)
	for _, link := range validString.FindAllString(matchString, -1) {
		link = strings.Trim(link, "\n")
		if strings.Contains(strings.ToLower(link), "product") {
			return link
		}
		return link
	}
	return
}

func parseSiteMapURL(matchString, host string) (validLink []string) {
	validString := regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^<">]+)`, host))
	for _, link := range validString.FindAllString(matchString, -1) {
		path, err := url.Parse(link)
		if err != nil {
			log.Fatalln(err)
			return
		}
		extension := strings.Split(path.Path, ".")
		if len(extension) > 1 {
			if strings.Contains(extension[1], "html") || strings.Contains(extension[1], "php") {
				validLink = append(validLink, link)
			}
		} else {
			validLink = append(validLink, link)
		}
	}
	return
}

// retrieve performs a HTTP Get request for the sitemap xml and decodes the results
func (prodDesc *ProductDescription) retrieve(siteMap, host string) ([]string, error) {
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")
	if err := bow.Open(siteMap); err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return parseSiteMapURL(bow.Body(), host), nil
}
