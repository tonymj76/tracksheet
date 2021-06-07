package structure

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/tebeka/selenium"
)

var robotURL = "%s://%s/robots.txt"

type (
	product struct {
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

//Search _
const (
	geckoDriverPath = "../vendor/geckodriver.exe"
	// geckoDriverPath = "../vendor/geckodriver-v0.29.1-linux64.tar"
	port         = 8080
	seleniumPath = "../vendor/selenium-server-standalone-3.141.59.jar"
)

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
		log.Println(err, "am here")
		return
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
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
	fmt.Println(elem.Text())
	return
}
