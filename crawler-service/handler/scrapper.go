package handler

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"cache_service/crawler-service/db"

	"cache_service/crawler-service/config"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
)

const (
	childrenSelector = "*"
	imageQuery       = "img"
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

	allSiteLink   map[string]int
	visited       map[string]struct{}
	count         = 0
	ahrefSelector map[string]int
)

type crawlconfig struct {
	website string
	depth   int
}

var (
	//key words to identify a product link
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

	//keywords on url not to visit
	skipSet = map[string]struct{}{
		"career":   none,
		"coupons":  none,
		"deal":     none,
		"help":     none,
		"sitemap":  none,
		"blog":     none,
		"policy":   none,
		"checkout": none,
		"content":  none,
		"account":  none,
		"login":    none,
		"log-in":   none,
		"contact":  none,
		"faq":      none,
		"privacy":  none,
		"about":    none,
		"customer": none,
		"icon":     none,
		"sign-up":  none,
		"signup":   none,

		".css":  none,
		".js":   none,
		".png":  none,
		".jpeg": none,
		".jpg":  none,
		".tiff": none,
		".ico":  none,
		".pdf":  none,
		".docx": none,
		".doc":  none,
	}
)

//identify possible url
func isPossibleProductURL(element *rod.Element, href string) (bool, string) {

	class := element.MustAttribute("class")
	if class == nil {
		return false, ""
	}
	class_details := strings.ReplaceAll(*class, " ", ",")

	// must not contain search in href
	if strings.Contains(strings.ToLower(href), search) {
		return false, ""
	}

	// contains product?
	if strings.Contains(strings.ToLower(href), product) {
		return true, class_details
	}

	// contains keyword text?
	text, err := element.Text()
	if err != nil {
		// is text in keyword
		_, ok := keywordSet[text]

		if ok {
			return true, class_details
		}
	}

	//get element children
	children, err := element.Elements(childrenSelector)

	if err == nil && len(children) > 0 {

		// has image tag?
		if _, err := children.First().Element(imageQuery); err == nil {
			return true, class_details
		}

		//get the href in the children node and recursively visit those too
		if len(children) > 0 {
			for _, child := range children {
				href, err := child.Attribute(`href`)
				if err != nil {
					logrus.Warnln("unable to get href", err)
					continue
				}

				if href != nil {
					if value, _ := isPossibleProductURL(child, *href); value {
						return true, class_details
					}
				}
			}
		}
	}
	return false, ""
}

func getLinks(base *url.URL, crawl_config crawlconfig) {
	//create a rod instance with user agent
	page := rod.New().MustConnect().MustPage(base.String())
	defer page.Close()

	userAgentHeaderOverride := proto.NetworkSetUserAgentOverride{
		UserAgent: config.GetEnv("userAgent"),
	}
	page.SetUserAgent(&userAgentHeaderOverride)
	// timeout
	page.Timeout(1 * time.Minute)
	//wait for page load
	page.MustWaitLoad()

	allLinks := page.MustElements(linksSelector)
	//iterate all items in the link list
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
		for key := range skipSet {
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
			possible, class := isPossibleProductURL(link, resolvedURLString)
			if possible {
				count++
				if count < 1001 {
					if class_count, ok := ahrefSelector[class]; ok {
						ahrefSelector[class] = class_count + 1
					} else {
						ahrefSelector[class] = 1
					}
					if count == 1000 {
						selector := getMaxclass(ahrefSelector)
						db.SaveSelector(selector, base.String())
					}
				}

				productMut.Lock()
				db.SaveToDB(resolvedURLString)
				productMut.Unlock()
			}
		}

	}
}

func getMaxclass(hredSelector map[string]int) string {
	count := 0
	return_key := ""
	for key, value := range hredSelector {
		if value > count {
			count = value
			return_key = key
		}
	}
	return return_key
}

//get all links on page
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
			logrus.Errorf("ERROR ON %s", crawl_config.website, tryError)
		}
	}

	waitgroup.Done()
}

func GetUrl(website string, depth, worker int) {
	//number of workers
	numWorkers := worker

	waitgroup = &sync.WaitGroup{}

	visitMut = &sync.RWMutex{}
	productMut = &sync.RWMutex{}
	allSiteMut = &sync.RWMutex{}

	config := crawlconfig{website: website, depth: depth}

	allSiteLink = map[string]int{}
	visited = map[string]struct{}{}
	ahrefSelector = map[string]int{}

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
		size := minInt(numWorkers)

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

		//add new keys
		for key, value := range allSiteLink {
			if _, ok := visited[key]; !ok {
				keys = append(keys, kvPair{key: key, depth: value - 1})
				delete(allSiteLink, key)
			}
		}
	}

	fmt.Println("---", len(keys))

}

func minInt(a int) int {
	if a >= 50 {
		return 50
	}
	return a
}
