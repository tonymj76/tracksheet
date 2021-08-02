package handler

import (
	"cache_service/crawler-service/config"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Connection connect http request and logger
type Connection struct {
	Logger                     *logrus.Logger
	Response                   *http.Response
	urlLink, accept, userAgent string
}

const (
	// url link to robots file
	RobotsURL = "%s://%s/robots.txt"
	// regex here map a url with a space or no space but must end with .xml
	XMLPattern = `(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.xml$)|(?P<path2>[^\s]*\.xml\s+$)`
	// get the links in html file
	HTMLLinkPattern = `(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^<">]+)`
)

func recoverMe() {
	if err := recover(); err != nil {
		logrus.Fatalln("recovering from panic ", err)
	}
}

// NewConnect create new Connection struct
func NewConnect(logger *logrus.Logger, urlPath, accept, userAgent string) *Connection {
	return &Connection{
		Logger:    logger,
		urlLink:   urlPath,
		accept:    accept,
		userAgent: userAgent,
	}
}

type ReturnList struct {
	// List of the product details
	List  []string
	Error error
}

// Search url and return result
func (connection *Connection) Search() ReturnList {
	parseURL, err := url.Parse(connection.urlLink)
	if err != nil {
		connection.Logger.WithError(err).Fatalln("make sure you entered a valid URL ")
	}
	defer recoverMe()

	switch filepath.Ext(connection.urlLink) {
	case ".gz":
		// check if the url is a .gz file
		result, err := connection.readDotGZFile(connection.urlLink, parseURL.Hostname())
		return ReturnList{
			List:  result,
			Error: err,
		}
	case ".xml":
		// check if the url is a xml file
		result, err := connection.retrieve(connection.urlLink, parseURL.Hostname())
		return ReturnList{
			List:  result,
			Error: err,
		}
	default:
		// normal urls and append robots to the url path
		response, err := GoClient(fmt.Sprintf(RobotsURL, parseURL.Scheme, parseURL.Hostname()), connection)
		if err != nil {
			connection.Logger.WithError(err).Fatalln("failed to get response")
		}
		byteReader, err := ioutil.ReadAll(response.Body)
		if err != nil {
			connection.Logger.WithError(err).Fatalln("failed to read byte")
		}
		siteMap := getSiteMapURL(string(byteReader))
		result, err := connection.retrieve(siteMap, parseURL.Hostname())
		return ReturnList{
			List:  result,
			Error: err,
		}
	}
}

// getSiteMapURL retrieve the xml sitemap
func getSiteMapURL(matchString string) (value string) {
	validString := regexp.MustCompile(XMLPattern)
	for _, link := range validString.FindAllString(matchString, -1) {
		link = strings.Trim(link, "\n")
		return link
	}
	return
}

// parseSiteMapURL retrieves all the url links from the html body
func parseSiteMapURL(matchString, host string) []string {
	validString := regexp.MustCompile(fmt.Sprintf(HTMLLinkPattern, host))
	return validString.FindAllString(matchString, -1)
}

// retrieve performs a HTTP Get request for the sitemap xml and decodes the results
func (connection *Connection) retrieve(siteMap, host string) ([]string, error) {
	resp, err := GoClient(siteMap, connection)
	if err != nil {
		connection.Logger.WithError(err).Fatalln("failed to get response")
	}
	byteReader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		connection.Logger.WithError(err).Fatalln("failed to read byte")
	}
	return parseSiteMapURL(string(byteReader), host), nil
}

// readDotGZFile performs a HTTP Get request for the zipped xml and decodes the results
func (connection *Connection) readDotGZFile(gzfile, host string) ([]string, error) {
	resp, err := GoClient(gzfile, connection)
	if err != nil {
		connection.Logger.WithError(err).Fatalln("failed to get response")
	}
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		connection.Logger.WithError(err).Fatalln("failed to unzip gz file")
	}
	byteReader, err := ioutil.ReadAll(reader)
	return parseSiteMapURL(string(byteReader), host), err
}

// GoClient to make http request and set brower header
func GoClient(urlPath string, connection *Connection) (*http.Response, error) {
	var client http.Client
	httpProxy := config.GetEnv("HTTP_PROXY")
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
	req.Header.Add("Accept", connection.accept)
	req.Header.Add("User-Agent", connection.userAgent)
	req.Header.Add("Accept-Charset", "utf8")
	response, err := client.Do(req)
	if response.StatusCode != 200 {
		connection.Logger.Println(response.Status)
	}
	return response, err
}
