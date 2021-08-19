package handler

import (
	"cache_service/crawler-service/db"
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

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

	productMut *sync.RWMutex

	waitgroup *sync.WaitGroup

	visited     sync.Map
	allSiteLink sync.Map
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
		"javascript",
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
		"mailto",
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

func queryWithSelector(page *rod.Page, base *url.URL, crawlConfig crawlConfiguration, selectorConfig, otherSelector []db.ProductSelector) {
	allLinks := page.MustElements(linksSelector)

	for _, productSelector := range selectorConfig {
		// get all nodes that contain the parsed selector
		classSelector := page.MustElements(*productSelector.AttrSelector)
		// transverse the slice that has all the nodes
		for _, link := range classSelector {
			// extract href attr from parent node
			href := link.MustAttribute(`href`)

			// if no href is found in parent node
			// transverse child node for href attr
			if href == nil {
				children, err := link.Elements(childrenSelector)
				if err != nil {
					fmt.Println("error parsing children ", err)
					continue
				}

				for _, child := range children {
					childHref, _ := child.Attribute(`href`)
					// if href found, add to DB
					if childHref != nil {
						parsedURL, err := url.Parse(strings.TrimSpace(*childHref))
						if err != nil {
							continue
						}
						resolvedURL := base.ResolveReference(parsedURL)
						resolvedURLString := resolvedURL.String()

						productMut.Lock()
						db.SaveToDB(resolvedURLString)
						getOtherDetails(page, resolvedURLString, base, otherSelector)
						productMut.Unlock()

					}
				}

			} else {
				// if href found in parent node
				// add to DB
				parsedURL, err := url.Parse(strings.TrimSpace(*href))
				if err != nil {
					continue
				}
				resolvedURL := base.ResolveReference(parsedURL)
				resolvedURLString := resolvedURL.String()

				productMut.Lock()
				db.SaveToDB(resolvedURLString)
				getOtherDetails(page, resolvedURLString, base, otherSelector)
				productMut.Unlock()
			}
		}
	}

	// iterate all items in the link list
	for _, link := range allLinks {
		inSkipSet := false

		href, err := link.Attribute(`href`)
		if err != nil {
			logrus.Warnln("UNABLE TO GET HREF ", err)
			inSkipSet = true
		}
		parsedURL, err := url.Parse(strings.TrimSpace(*href))
		if err != nil {
			inSkipSet = true
		}
		resolvedURL := base.ResolveReference(parsedURL)
		resolvedURLString := resolvedURL.String()
		if base.Hostname() != resolvedURL.Hostname() {
			inSkipSet = true
		}

		for _, key := range skipSet {
			if strings.Contains(resolvedURL.RequestURI(), key) {
				inSkipSet = true
				break
			}
		}
		if !inSkipSet {
			allSiteLink.Store(resolvedURLString, crawlConfig.depth-1)
		}
	}

	fmt.Println("Getiing with product selector.....")
}

func queryRecursively(page *rod.Page, base *url.URL, crawlConfig crawlConfiguration, otherSelector []db.ProductSelector) {
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
			allSiteLink.Store(resolvedURLString, crawlConfig.depth-1)
			getOtherDetails(page, resolvedURLString, base, otherSelector)

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

func getOtherDetails(page *rod.Page, urlValue string, base *url.URL, otherSelector []db.ProductSelector) {
	parsed_url, err := url.Parse(urlValue)
	if err != nil {
		fmt.Println("error parsing url ", err)
	} else {
		if base.Hostname() == parsed_url.Hostname() {

			for _, selector_item := range otherSelector {
				classSelector := page.MustElements(*selector_item.AttrSelector)
				for _, node := range classSelector {
					text := node.MustText()
					text_trim := strings.TrimSpace(text)
					if text_trim != "" {
						db.SaveOtherAttr(urlValue, *selector_item.AttrName, text_trim)
					}
				}

			}
		}
	}

}
