package handler

import (
	"cache_service/crawler-service/config"
	"cache_service/crawler-service/db"
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
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

	visitMut   *sync.RWMutex
	productMut *sync.RWMutex
	allSiteMut *sync.RWMutex

	waitgroup *sync.WaitGroup

	allSiteLink map[string]int
	visited     map[string]struct{}
)

type crawlconfig struct {
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

// identify possible url
func isPossibleProductURL(element *rod.Element, href string) bool {
	// must not contain search in href
	if strings.Contains(strings.ToLower(href), search) {
		return false
	}

	// contains product?
	if strings.Contains(strings.ToLower(href), product) {
		return true
	}

	// contains keyword text?
	text, err := element.Text()
	if err != nil {
		// is text in keyword
		_, ok := keywordSet[text]

		if ok {
			return true
		}
	}

	// get element children
	children, err := element.Elements(childrenSelector)

	if err == nil && len(children) > 0 {
		// has image tag?
		for _, val := range children {
			if strings.Contains(val.String(), "img") {
				return true
			}

		}

		// get the href in the children node and recursively visit those too
		if len(children) > 0 {
			for _, child := range children {
				href, err := child.Attribute(`href`)
				if err != nil {
					logrus.Warnln("unable to get href", err)
					continue
				}

				if href != nil {
					if value := isPossibleProductURL(child, *href); value {
						return true
					}
				}
			}
		}
	}
	return false
}

func getLinks(base *url.URL, crawl_config crawlconfig) {
	// create a rod instance with user agent
	page := rod.New().MustConnect().MustPage(base.String())
	defer page.Close()

	userAgentHeaderOverride := proto.NetworkSetUserAgentOverride{
		UserAgent: config.GetEnv("userAgent"),
	}
	page.SetUserAgent(&userAgentHeaderOverride)
	// timeout
	page.Timeout(1 * time.Minute)
	// wait for page load
	page.MustWaitLoad()
	selectorConfic, err := getSelectorConfig(base.String(), "product_selector")
	if err != nil {
		fmt.Println("error getting selector ", err)
	}
	switch {
	case selectorConfic.AttrValue != nil:
		classSelector := page.MustElements(*selectorConfic.AttrValue)

		domChildren := []string{}
		for _, item := range classSelector {
			domChildren = append(domChildren, parseURLFromSelector(item.MustHTML())...)
		}
		for _, child := range domChildren {
			parsed_url, err := url.Parse(strings.TrimSpace(child))
			if err != nil {
				continue
			}
			resolvedURL := base.ResolveReference(parsed_url)
			resolvedURLString := resolvedURL.String()

			// skip if not same hostname
			if base.Hostname() != resolvedURL.Hostname() {
				continue
			}

			for _, key := range skipSet {
				if strings.Contains(resolvedURL.RequestURI(), key) {
					continue
				}
			}

			allSiteMut.Lock()
			allSiteLink[resolvedURLString] = crawl_config.depth - 1
			allSiteMut.Unlock()

			productMut.Lock()
			db.SaveToDB(resolvedURLString)
			productMut.Unlock()
		}
		fmt.Println("am in my >>>>>>>>>>>>>>>>>>>>>>>>>>...")
	default:
		fmt.Println("am in default >>>>>>>>>>>>>>>>>>>>>>>>>>...")
		allLinks := page.MustElements(linksSelector)
		// iterate all items in the link list
		for _, link := range allLinks {

			href, err := link.Attribute(`href`)
			if err != nil {
				logrus.Warnln("UNABLE TO GET HREF ", err)
				continue
			}

			parsed_url, err := url.Parse(strings.TrimSpace(*href))
			if err != nil {
				logrus.Warnf("unable to parse link:%s\n%v", *href, err)
				continue
			}

			resolvedURL := base.ResolveReference(parsed_url)
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
				allSiteLink[resolvedURLString] = crawl_config.depth - 1
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
}

// get all links on page
func getAllLinksOnPage(crawl_config crawlconfig) {
	crawl_config.website = strings.TrimSpace(crawl_config.website)

	// visited?
	visitMut.RLock()
	_, ok := visited[crawl_config.website]
	visitMut.RUnlock()

	if !ok && crawl_config.depth >= 0 {
		// mark visited
		visitMut.Lock()
		visited[crawl_config.website] = none
		visitMut.Unlock()

		base, err := url.Parse(crawl_config.website)
		if err != nil {
			logrus.Warnf("unable to parse website: %s\n%v", crawl_config.website, err)
			return
		}

		tryError := rod.Try(func() {
			getLinks(base, crawl_config)
		})

		if errors.Is(tryError, context.DeadlineExceeded) {
			// code for timeout error
			logrus.Errorf("TIMEOUT ON %s", crawl_config.website, tryError)
		} else if tryError != nil {
			// code for other types of error
			logrus.Errorf("ERROR ON %v %v", crawl_config.website, tryError)
		}
	}

	waitgroup.Done()
}

func GetUrl(website string, depth, worker int) {
	// number of workers
	numWorkers := worker

	waitgroup = &sync.WaitGroup{}

	visitMut = &sync.RWMutex{}
	productMut = &sync.RWMutex{}
	allSiteMut = &sync.RWMutex{}

	config := crawlconfig{website: website, depth: depth}

	allSiteLink = map[string]int{}
	visited = map[string]struct{}{}

	// get root urls
	waitgroup.Add(1)
	go getAllLinksOnPage(config)
	waitgroup.Wait()

	type kvPair struct {
		key   string
		depth int
	}

	keys := []kvPair{}

	for k, v := range allSiteLink {
		if _, ok := visited[k]; !ok {
			keys = append(keys, kvPair{k, v})
		}
	}

	for len(keys) > 0 {
		size := minInt(len(keys), numWorkers)

		for i := 0; i < size; i++ {
			depth := keys[i].depth
			fmt.Println("DEPTH::::::::", depth)

			// visited[keys[i].key] = none
			delete(allSiteLink, keys[i].key)

			newConfig := crawlconfig{website: keys[i].key, depth: keys[i].depth - 1}
			waitgroup.Add(1)
			go getAllLinksOnPage(newConfig)
		}

		waitgroup.Wait()

		keys = keys[size:]

		// add new keys
		for key, value := range allSiteLink {
			if _, ok := visited[key]; !ok {
				keys = append(keys, kvPair{key: key, depth: value - 1})
				delete(allSiteLink, key)
			}
		}
	}

	fmt.Println("---", len(keys))
}

func minInt(a, b int) int {
	if a >= b {
		return b
	}
	return a
}

func getSelectorConfig(link, attr string) (*db.ProductSelector, error) {
	product_selector, err := GetUrlConfig(link, attr)
	return &product_selector, err
}

// parseURLFromSelector retrieves all the url links from the html body
func parseURLFromSelector(matchString string) []string {
	niceStr := []string{}
	HTMLLinkPattern := `href=[^\s]+`
	re := regexp.MustCompile(`(href=")`)
	validString := regexp.MustCompile(HTMLLinkPattern)

	for _, str := range validString.FindAllString(matchString, -1) {
		niceStr = append(niceStr, strings.Split(re.ReplaceAllString(str, ""), `"`)[0])
	}
	return niceStr
}

func GetUrlConfig(url, attr_name string) (db.ProductSelector, error) {
	var config db.ProductSelector
	if err := db.DB.Where("root_url = ?", &url).Where("attr_name = ?", &attr_name).First(&config).Error; err != nil {
		return config, err
	}
	return config, nil
}
