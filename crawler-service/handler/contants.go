package handler

import (
	"context"
	"sync"

	"cache_service/crawler-service/db"
	"fmt"
	"net/url"

	"strings"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

const (
	childrenSelector = "*"
	href             = "href"
	search           = "search"
	linksSelector    = "[href]"
	product          = "product"
)

var (
	none = struct{}{}

	browser *rod.Browser

	semCtx context.Context
	sem    *semaphore.Weighted

	visitMut   *sync.RWMutex
	productMut *sync.RWMutex
	allSiteMut *sync.RWMutex

	waitgroup *sync.WaitGroup

	allSiteLink map[string]int
	visited     map[string]struct{}
)

type crawlConfiguration struct {
	website string
	depth   int
}

var (
	// key words to identify a product link
	keywordSet = map[string]struct{}{
		"price":       none,
		"add to cart": none,
		"quantity":    none,
		"rating":      none,
		"save":        none,
		"buy it now":  none,
		"buy now":     none,
		"cart":        none,
		"buy":         none,
	}

	// keywords on url not to visit
	skipSet = [...]string{
		"career",
		"coupons",
		"deal",
		"help",
		"sitemap",
		"blog",
		"policy",
		"checkout",
		"content",
		"account",
		"login",
		"log-in",
		"contact",
		"faq",
		"privacy",
		"about",
		"customer",
		"icon",
		"sign-up",
		"signup",
		".css",
		".js",
		".png",
		".jpeg",
		".jpg",
		".tiff",
		".ico",
		".pdf",
		".docx",
		".doc",
		".woff2",
		".woff",
	}
)

func queryWithSelector(page *rod.Page, base *url.URL, crawlConfig crawlConfiguration, selector_config []db.ProductSelector) {
	allLinks := page.MustElements(linksSelector)

	for _, product_selector := range selector_config {
		//get all nodes that contain the parsed selector
		classSelector := page.MustElements(*product_selector.AttrValue)
		//transverse the slice that has all the nodes
		for _, link := range classSelector {
			//extract href attr from parent node
			href, _ := link.Attribute(`href`)

			//if no href is found in parent node
			//transverse child node for href attr
			if href == nil {
				children, err := link.Elements(childrenSelector)
				if err != nil {
					fmt.Println("error parsing children ", err)
					continue
				}

				for _, child := range children {
					child_href, _ := child.Attribute(`href`)
					//if href found, add to DB
					if child_href != nil {
						parsed_url, _ := url.Parse(strings.TrimSpace(*child_href))
						resolvedURL := base.ResolveReference(parsed_url)
						resolvedURLString := resolvedURL.String()

						productMut.Lock()
						db.SaveToDB(resolvedURLString)
						productMut.Unlock()

					}
				}

			} else {
				//if href found in parent node
				//add to DB
				parsed_url, _ := url.Parse(strings.TrimSpace(*href))
				resolvedURL := base.ResolveReference(parsed_url)
				resolvedURLString := resolvedURL.String()

				productMut.Lock()
				db.SaveToDB(resolvedURLString)
				productMut.Unlock()
			}
		}
	}

	// iterate all items in the link list
	for _, link := range allLinks {

		href, err := link.Attribute(`href`)
		if err != nil {
			logrus.Warnln("UNABLE TO GET HREF ", err)
			continue
		}
		parsedURL, _ := url.Parse(strings.TrimSpace(*href))
		resolvedURL := base.ResolveReference(parsedURL)
		resolvedURLString := resolvedURL.String()
		if base.Hostname() != resolvedURL.Hostname() {
			continue
		}

		inSkipSet := false
		for _, key := range skipSet {
			if strings.Contains(resolvedURL.RequestURI(), key) {
				inSkipSet = true
				break
			}
		}
		if !inSkipSet {

			allSiteMut.Lock()
			allSiteLink[resolvedURLString] = crawlConfig.depth - 1
			allSiteMut.Unlock()
		}
	}

	fmt.Println("Getiing with product selector.....")
}

func queryRecursively(page *rod.Page, base *url.URL, crawlConfig crawlConfiguration, selector_config []db.ProductSelector) {
	fmt.Println("Getting via recursion.......")
	allLinks := page.MustElements(linksSelector)
	// iterate all items in the link list
	for _, link := range allLinks {

		href, err := link.Attribute(`href`)
		if err != nil {
			logrus.Warnln("UNABLE TO GET HREF ", err)
			continue
		}

		parsedURL, err := url.Parse(strings.TrimSpace(*href))
		if err != nil {
			logrus.Warnf("unable to parse link:%s\n%v", *href, err)
			continue
		}

		resolvedURL := base.ResolveReference(parsedURL)
		resolvedURLString := resolvedURL.String()

		// skip if not same hostname
		if base.Hostname() != resolvedURL.Hostname() {
			continue
		}

		// skip if in skipset
		inSkipSet := false
		for _, key := range skipSet {
			if strings.Contains(resolvedURL.RequestURI(), key) {
				inSkipSet = true
				break
			}
		}

		if !inSkipSet {

			// write
			allSiteMut.Lock()
			allSiteLink[resolvedURLString] = crawlConfig.depth - 1
			allSiteMut.Unlock()
			// check if it's product URL
			possible := isPossibleProductURL(link, resolvedURLString)
			if possible {
				productMut.Lock()
				db.SaveToDB(resolvedURLString)
				productMut.Unlock()
			}
		}
	}
}
