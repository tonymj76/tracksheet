package handler

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"time"

	"github.com/sirupsen/logrus"
)

//Connection connect http request and logger
type Connection struct {
	Logger                     *logrus.Logger
	Response                   *http.Response
	urlLink, accept, userAgent string
}

var robotURL = "%s://%s/robots.txt"

func recoverMe() {
	if err := recover(); err != nil {
		logrus.Fatalln("recovering from panic ", err)
	}
}

//NewConnect create new Connection struct
func NewConnect(logger *logrus.Logger, urlPath, accept, userAgent string) *Connection {
	return &Connection{
		Logger:    logger,
		urlLink:   urlPath,
		accept:    accept,
		userAgent: userAgent,
	}
}

// Search url and return result
func (c *Connection) Search() ([]string, error) {
	parseURL, err := url.Parse(c.urlLink)
	if err != nil {
		c.Logger.WithError(err).Fatalln("make sure you entered a valid URL ")
	}
	defer recoverMe()
	// check if the url is a xml.gz file
	if filepath.Ext(c.urlLink) == ".gz" {
		resp, err := goClient(c.urlLink, c)
		if err != nil {
			c.Logger.WithError(err).Fatalln("failed to make request to the server")
		}
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		byteReader, err := ioutil.ReadAll(reader)
		if err != nil {
			c.Logger.WithError(err).Fatalln("failed to read byte")
		}
		result := parseSiteMapURL(string(byteReader), parseURL.Hostname())
		return result, nil
	}

	// normal urls and append robots to the url path
	resp, err := goClient(fmt.Sprintf(robotURL, parseURL.Scheme, parseURL.Hostname()), c)
	if err != nil {
		c.Logger.WithError(err).Fatalln("failed to get response")
	}
	byteReader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger.WithError(err).Fatalln("failed to read byte")
	}
	siteMap := getSiteMapURL(string(byteReader))
	result, err := c.retrieve(siteMap, parseURL.Hostname())
	return result, err
}

//getSiteMapURL retrieve the xml sitemap
func getSiteMapURL(matchString string) (value string) {
	//regex here map a url with a space or no space but must end with .xml
	validString := regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.xml$)|(?P<path2>[^\s]*\.xml\s+$)`)
	for _, link := range validString.FindAllString(matchString, -1) {
		link = strings.Trim(link, "\n")
		if strings.Contains(strings.ToLower(link), "product") {
			return link
		} else {
			return link
		}
	}
	return
}

//parseSiteMapURL retrieves all the url links from the html body
func parseSiteMapURL(matchString, host string) (validLink []string) {
	validString := regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^<">]+)`, host))
	for _, link := range validString.FindAllString(matchString, -1) {
		path, err := url.Parse(link)
		if err != nil {
			log.Fatalln("error while parsing link as url", err)
			return
		}
		// filter out static assets
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
func (c *Connection) retrieve(siteMap, host string) ([]string, error) {
	resp, err := goClient(siteMap, c)
	if err != nil {
		c.Logger.WithError(err).Fatalln("failed to get response")
	}
	byteReader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger.WithError(err).Fatalln("failed to read byte")
	}
	return parseSiteMapURL(string(byteReader), host), nil
}

// goClient to make http request and set brower header
func goClient(urlPath string, c *Connection) (*http.Response, error) {
	var (
		client http.Client
	)
	httpProxy := os.Getenv("HTTP_PROXY")
	if httpProxy != "" {
		proxyUrl, _ := url.Parse(httpProxy)
		client = http.Client{Timeout: 20 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = http.Client{Timeout: 20 * time.Second}
	}
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", c.accept)
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Accept-Charset", "utf8")
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		c.Logger.Println(resp.Status)
	}
	return resp, err
}

// RelatedURL returns similar product URL
func RelatedURL(urls []string) (duplicateURL []string) {
	paths := []string{}
	occurencePaths := []string{}

	for _, item := range urls {
		link, _ := url.Parse(item)
		paths = append(paths, link.Path)
	}

	for x := 0; x < len(paths); x++ {
		if !check(paths[x], occurencePaths) {
			occurence := 0
			for y := 0; y < len(paths); y++ {
				if compareDuplicate(paths[x], paths[y]) {
					occurence++
				}
			}
			if occurence >= 2 {
				occurencePaths = append(occurencePaths, paths[x])
			}
		}
	}

	for x := 0; x < len(urls); x++ {
		for u := 0; u < len(occurencePaths); u++ {
			if strings.Contains(urls[x], occurencePaths[u]) {
				duplicateURL = append(duplicateURL, urls[x])
			}
		}
	}
	return
}

// check checks if a string (str) exits in strs []string
func check(str string, strs []string) bool {
	for _, item := range strs {
		if compareDuplicate(item, s) {
			return true
		}
	}
	return false
}

// compareDuplicate takes two string and checks if the have
// 30% similar bytes
func compareDuplicate(x, y string) bool {
	byteX := []byte(x)
	byteY := []byte(y)
	var (
		minLen int
		count  int
	)
	if len(byteX) > len(byteY) {
		minLen = len(byteY)
	} else {
		minLen = len(byteX)
	}

	if minLen < 1 {
		return false
	}

	for i := 0; i < minLen; i++ {
		if byteX[i] == byteY[i] {
			count++
		}
	}
	percentage := (count * 100) / minLen
	return percentage > 30
}
