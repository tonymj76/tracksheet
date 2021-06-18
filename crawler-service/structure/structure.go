package structure

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/headzoo/surf.v1"
)

var robotURL = "%s://%s/robots.txt"

// ProductDescription implements the Matcher interface.
type ProductDescription struct{}

// GzProduct implements the Matcher interface.
type GzProduct struct{}

// Search url with xml.zip
func (gzProduct *GzProduct) Search(urlString string) (result []string, err error) {
	parseURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatalln("make sure you entered a valid URL ", err)
	}
	reader, err := gzip.NewReader(getReader(urlString))
	if err != nil {
		log.Fatalln(err)
	}
	byteReader, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return
	}
	result = parseSiteMapURL(string(byteReader), parseURL.Hostname())
	return
}

// Search normal urls and append robots to the url path
func (prodDesc *ProductDescription) Search(urlString string) (result []string, err error) {
	parseURL, err := url.Parse(urlString)
	if err != nil {
		log.Fatalln("make sure you entered a valid URL ", err)
	}
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")
	if err = bow.Open(fmt.Sprintf(robotURL, parseURL.Scheme, parseURL.Hostname())); err != nil {
		log.Fatalln("error while opening link", err)
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
			log.Fatalln("error while parsing link as url", err)
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

func getReader(url_path string) io.ReadCloser {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url_path, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error getting reader %s", err)
	}
	return resp.Body
}
