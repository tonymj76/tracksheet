package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type ReturnList struct {
	List  []string
	Key   string
	Error error
}

// var fileString = "https://www.vitacost.com/vitacost-glucosamine-msm"

// var fileString = "https://www.vitacost.com/productsearch.aspx?rid=1301346.01+1307326.01+1307331.01+1307384.01+1024.05"

var fileString = "https://www.konga.com/"

const PricePattern = `[$]\s?[+-]?[0-9]{1,3}(?:(?:,?[0-9]{3}))*(?:\.[0-9]{1,2})?`

func main() {
	count := 0
	var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"
	c := colly.NewCollector(colly.UserAgent(userAgent))
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		fmt.Println(h.DOM.Contents().Text())
		count++
		fmt.Println(count)
		// h.ForEach("*", func(i int, h *colly.HTMLElement) {
		// 	// fmt.Println(h.Text)
		// })
	})

	c.OnResponse(func(r *colly.Response) {
		p := regexp.MustCompile(PricePattern)
		// fmt.Println(string(r.Body))
		fmt.Println(p.FindAllString(string(r.Body), -1))
	})
	c.Visit(fileString)

	// r := GetUrl("https://vitacost.com/", userAgent)
	// if r.Error != nil {
	// 	fmt.Println(r.Error)
	// } else {
	// 	fmt.Println(r.List)
	// }

	// bigString := strings.Split(fileString, " ")
	// holderMap := make(map[string]int)

	// for _, s := range bigString {
	// 	if _, ok := holderMap[s]; ok {
	// 		holderMap[s]++
	// 	} else {
	// 		holderMap[s] = 1
	// 	}
	// }

	// for key, value := range holderMap {
	// 	if value == 1 {
	// 		fmt.Println()
	// 	}
	// }

	// Test cases
	// https://www.vitacost.com/vitacost-glucosamine-msm detail page
	// https://www.vitacost.com/productsearch.aspx?rid=1301346.01+1307326.01+1307331.01+1307384.01+1024.05
}

func GetUrl(link, userAgent string) *ReturnList {
	links := []string{}
	base, err := url.Parse(link)

	if err != nil {
		logrus.WithError(err).Fatalln("Error parsing URL")
	}
	if !strings.Contains(base.Host, "www") {
		return &ReturnList{Error: fmt.Errorf("the URL is %v want %s://www.%s", base, base.Scheme, base.Host)}
	}
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.UserAgent(userAgent),
		colly.IgnoreRobotsTxt(),
		colly.DisallowedURLFilters(
			// don't craw any url with the following keywords
			regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>(?i)Career|deal|help|Sitemap|blog|policy|CheckOut|content|account|login|contact|faq|privacy-policy|about|customer[^\s]*)`, base.Hostname())),
			// don't craw image or png from the url
			regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.jpeg$)|(?P<path2>[^\s]*\.jpg$)|(?P<path2>[^\s]*\.png$)`),
		),
		colly.URLFilters(
			//resolve url if item domain is not the same as base domain, link is unrelated. skip
			regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^\s]*)`, base.Hostname())),
		),
		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		// add a job to clear this in every 24hrs......
		colly.CacheDir("./cache"),
	)
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 6, RandomDelay: 5 * time.Second})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// visit link found on page
		//Only those links are visited which are matched by the URLFilter regexps
		// 	v := e.ChildAttr("img", "src")
		// e.ForEach(".normal-price", func(i int, h *colly.HTMLElement) {
		// 	fmt.Println(h.ChildText(".price"))
		// })
		// fmt.Println(e.DOM.Html())

		// if v != "" {
		// 	err := c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		// 	if err != nil {
		// 		// fmt.Println(err)
		// 	}
		// }
		err := c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if err != nil {
			// fmt.Println(err)
		}
		// e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		// r.Ctx.Put("url", r.URL.String())
		fmt.Println(r.URL)
		// links = append(links, r.URL.String())
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Println(r.Ctx.Get("url"))
	// 	// fmt.Println(string(r.Body))
	// })

	err = c.Visit(link)
	c.Wait()
	fmt.Println("finishing >>>>>>>>>>>>>>.")
	//group links
	return &ReturnList{
		List:  links,
		Key:   link,
		Error: err,
	}
}
