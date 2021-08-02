package main

//// visit with same base

// line 73 keywords
// css, jpg, png, .ico, .pdf
// css, js, png, jpeg, jpg, tiff, icon, ico, js, pdf, docx, doc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
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

	wg *sync.WaitGroup

	allSiteLink map[string]int
	visited     map[string]struct{}
	productLink map[string]struct{}
)

type crawlconfig struct {
	website string
	depth   int
}

type linksConfig struct {
	crawlconfig
	base *url.URL
}

var (
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
		"contact":  none,
		"faq":      none,
		"privacy":  none,
		"about":    none,
		"customer": none,
		"icon":     none,

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

func isPossibleProductURL(el *rod.Element, href string) bool {
	// must not contain search in href
	if strings.Contains(strings.ToLower(href), search) {
		return false
	}

	// contains product?
	if strings.Contains(strings.ToLower(href), product) {
		return true
	}

	// contains keyword text?
	text, err := el.Text()
	if err != nil {
		// is text in keyword
		_, ok := keywordSet[text]

		if ok {
			return true
		}
	}

	children, err := el.Elements(childrenSelector)

	if err == nil && len(children) > 0 {
		// has image tag?
		if _, err := children.First().Element(imageQuery); err == nil {
			return true
		}

		if len(children) > 0 {
			for _, child := range children {
				hr, err := child.Attribute(`href`)
				if err != nil {
					logrus.Warnln("unable to get href", err)
					continue
				}

				if hr != nil {
					if isPossibleProductURL(child, *hr) {
						return true
					}
				}
			}
		}
	}
	return false
}

func getLinks(base *url.URL, config crawlconfig) {
	page := rod.New().MustConnect().MustPage(base.String())
	defer page.Close()

	// timeout
	page.Timeout(1 * time.Minute)
	page.MustWaitLoad()

	allLinks := page.MustElements(linksSelector)

	for _, link := range allLinks {

		hr, err := link.Attribute(`href`)
		if err != nil {
			logrus.Warnln("UNABLE TO GET HREF ", err)
			continue
		}

		l, err := url.Parse(strings.TrimSpace(*hr))
		if err != nil {
			logrus.Warnf("unable to parse link:%s\n%v", *hr, err)
			continue
		}

		resolvedURL := base.ResolveReference(l)
		resolvedURLString := resolvedURL.String()

		// skip if not same hostname
		if base.Hostname() != resolvedURL.Hostname() {
			fmt.Println("not equal", resolvedURLString)
			continue
		}

		// skip if in skipset
		inSkipSet := false
		for key := range skipSet {
			if strings.Contains(resolvedURL.RequestURI(), key) {
				fmt.Println("skipping", resolvedURLString)
				inSkipSet = true
				break
			}
		}

		if !inSkipSet {

			// write
			allSiteMut.Lock()
			allSiteLink[resolvedURLString] = config.depth - 1
			allSiteMut.Unlock()

			// check if it's product URL
			possible := isPossibleProductURL(link, resolvedURLString)

			if possible {
				// add it to link
				productMut.RLock()
				_, ok := productLink[resolvedURLString]
				productMut.RUnlock()

				if !ok {
					logrus.Infoln("product link: ", resolvedURL, len(productLink))
					productMut.Lock()
					productLink[resolvedURLString] = none
					productMut.Unlock()
				}
			}
		}

	}
}

func getAllLinksOnPage(config crawlconfig) {
	// logrus.Printf("getting: %s\n", config.website)

	config.website = strings.TrimSpace(config.website)

	// visited?
	visitMut.RLock()
	_, ok := visited[config.website]
	visitMut.RUnlock()

	if !ok && config.depth >= 0 {
		// mark visited
		visitMut.Lock()
		visited[config.website] = none
		visitMut.Unlock()

		base, err := url.Parse(config.website)
		if err != nil {
			logrus.Warnf("unable to parse website: %s\n%v", config.website, err)
			return
		}

		tryError := rod.Try(func() {
			getLinks(base, config)
		})

		if errors.Is(tryError, context.DeadlineExceeded) {
			// code for timeout error
			logrus.Errorf("TIMEOUT ON %s", config.website, tryError)
		} else if tryError != nil {
			// code for other types of error
			logrus.Errorf("ERROR ON %s", config.website, tryError)
		}
	}

	wg.Done()
}

func main() {
	website := "https://www.vitacost.com/"
	// website := "https://www.cdw.ca/"

	const numWorkers = 10

	wg = &sync.WaitGroup{}

	visitMut = &sync.RWMutex{}
	productMut = &sync.RWMutex{}
	allSiteMut = &sync.RWMutex{}

	config := crawlconfig{website: website, depth: 1}

	allSiteLink = map[string]int{}
	visited = map[string]struct{}{}
	productLink = map[string]struct{}{}

	// get root urls
	wg.Add(1)
	go getAllLinksOnPage(config)
	wg.Wait()

	fmt.Println("links #", len(allSiteLink))
	fmt.Println("visited #", len(visited))

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
		size := minInt(numWorkers, len(keys))

		fmt.Println("all links #", len(keys))

		for i := 0; i < size; i++ {
			depth := keys[i].depth
			fmt.Println("DEPTH::::::::", depth)

			// visited[keys[i].key] = none
			delete(allSiteLink, keys[i].key)

			newConfig := crawlconfig{website: keys[i].key, depth: keys[i].depth}
			wg.Add(1)
			go getAllLinksOnPage(newConfig)
		}

		wg.Wait()

		keys = keys[size:]

		// add new keys
		for key, value := range allSiteLink {
			if _, ok := visited[key]; !ok {
				fmt.Println("ADDING............", key)
				keys = append(keys, kvPair{key: key, depth: value})
				delete(allSiteLink, key)
			}
		}
	}

	fmt.Println("---", len(keys))

	fmt.Println("visited #", len(visited))
	fmt.Println("product #", len(productLink))
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
