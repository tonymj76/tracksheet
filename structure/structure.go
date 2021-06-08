package structure

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/tebeka/selenium"
)

var robotURL = "%s://%s/robots.txt"

type (
	Product struct {
		XMLName xml.Name `xml:"url"`
		Loc     string   `xml:"loc"`
	}

	feed struct {
	}
	//Result of the search
	Result struct {
		URL     string
		Content string
	}
	// ProductDescription use for searching
	ProductDescription struct{}
)

// //Search _
// func (prodDesc *ProductDescription) Search(urlString string) (retult []*Result, err error) {
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

const (
	geckoDriverPath = "vendor/geckodriver"
	port            = 8080
	seleniumPath    = "vendor/selenium-server-standalone-3.141.59.jar"
)

func headlessBrowser() (selenium.WebDriver, error) {

	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug((true))
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	return selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
}

func (prodDesc *ProductDescription) Search(urlString string) (retult []*Result, err error) {
	parseURL, err := url.Parse(urlString)
	if err != nil {
		log.Println(err)
		return
	}
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug((true))
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	// wd, err := headlessBrowser()
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	if err = wd.Get(fmt.Sprintf(robotURL, parseURL.Scheme, parseURL.Hostname())); err != nil {
		log.Println(err)
		return
	}
	// Get a reference to the text box containing code.
	elem, err := wd.FindElement(selenium.ByCSSSelector, "pre")
	text, err := elem.Text()
	p, err := prodDesc.retrieve(getSiteMapURL(text)[1])
	fmt.Println(&p)
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

// retrieve performs a HTTP Get request for the sitemap xml and decodes the results
func (prodDesc ProductDescription) retrieve(matchString string) (*Product, error) {
	if matchString == "" {
		return nil, fmt.Errorf("should not be empty")
	}
	resp, err := http.Get(matchString)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	// Check the status code for a 200 so we know we have received a
	// proper response.
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
	}
	var file Product
	err = xml.NewDecoder(resp.Body).Decode(&file)
	return &file, err
}
