package handler

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type ReturnList struct {
	// List of the product details
	List  []string
	Key   string
	Error error
}

const PricePattern = `\d{1,3}(?:[.,]\d{3})*(?:[.,]\d{2})`

// checkImgPrice check for price and if the anchor tag has img
func checkImgPrice(a string, b *goquery.Selection) bool {
	re := regexp.MustCompile(PricePattern)
	_, found := b.Find("img").Attr("src")
	html, _ := b.Html()
	if re.MatchString(html) && found {
		return true
	}
	return false
}

//urlResolver add www. if the url don't have it
func urlResolver(link string) string {
	base, _ := url.Parse(link)
	if !strings.Contains(base.Host, "www") {
		return fmt.Sprintf("%s://www.%s%s?%s", base.Scheme, base.Host, base.Path, base.RawQuery)
	}
	return link
}

// keywords check for certain keywords in the dom eg like `buy now`
func keywords(a string, b *goquery.Selection) bool {
	kwords := []string{"price", "add to cart", "quantity", "rating", "save", "buy it now", "buy now", "cart", "buy"}
	html, _ := b.Html()
	for _, keys := range kwords {
		if strings.Contains(strings.ToLower(html), keys) {
			return true
		}
	}
	return false
}

// GetUrl given the url and depth, recursively get all related url of the site
func GetUrl(link, accept, userAgent string, depth int) ReturnList {
	link = urlResolver(link)
	links := []string{}
	base, err := url.Parse(link)
	if err != nil {
		logrus.WithError(err).Fatalln("Error parsing URL")
	}
	//get base domain
	base_domain := fmt.Sprintf("%s://%s", base.Scheme, base.Host)

	c := colly.NewCollector(
		colly.MaxDepth(depth),
		colly.Async(true),
		colly.UserAgent(userAgent),
		colly.IgnoreRobotsTxt(),
		colly.DisallowedURLFilters(
			// don't scap any url with the following keywords
			regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>(?i)Career|deal|help|Sitemap|blog|policy|CheckOut|content|account|login|contact|faq|privacy-policy|about|customer|pages|cart|search[^\s]*)`, base.Hostname())),
			// don't scap image or png from the url
			regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.jpeg$)|(?P<path2>[^\s]*\.jpg$)|(?P<path2>[^\s]*\.png$)`),
		),
		colly.URLFilters(
			//resolve url if item domain is not the same as base domain, link is unrelated. skip
			regexp.MustCompile(fmt.Sprintf(`(?P<head>http[s]?:\/\/)(?P<host>%s\/)(?P<path>[^\s]*)`, base.Hostname())),
		),
		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./cache"),
	)
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 6, RandomDelay: 10 * time.Second})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		dom := e.DOM
		if keywords(e.Attr("href"), dom) || checkImgPrice(e.Attr("href"), dom) {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		//resolve url if item domain is not the same as base domain, link is unrelated. skip
		item_domain := fmt.Sprintf("%s://%s", r.URL.Scheme, r.URL.Host)
		if item_domain == base_domain {
			links = append(links, r.URL.String())
		}
	})

	err = c.Visit(link)
	c.Wait()
	//group links
	list := removeDuplicateValues(links)
	resolved_list := resolveUrl(link, list)
	return ReturnList{
		List:  resolved_list,
		Key:   "",
		Error: err,
	}
}

func resolveUrl(base_url string, urls []string) []string {
	base, _ := url.Parse(base_url)
	return_slice := []string{}
	for _, val := range urls {
		u, err := url.Parse(val)
		if err != nil {
			fmt.Println("error resolving url ", err)
			continue
		}
		return_slice = append(return_slice, base.ResolveReference(u).String())
	}
	return return_slice
}

//remove duplocate item from slice
func removeDuplicateValues(UrlSlice []string) []string {
	fmt.Println("resolving url...")
	keys := make(map[string]struct{})
	list := []string{}

	for _, entry := range UrlSlice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
