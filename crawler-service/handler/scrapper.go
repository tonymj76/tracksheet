package handler

import (
	"cache_service/crawler-service/config"
	"cache_service/crawler-service/db"
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
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

func setupPage(page string) *rod.Page {
	userAgentConfig := proto.NetworkSetUserAgentOverride{
		UserAgent: config.GetEnv("userAgent"),
	}

	// browser exists?
	if browser != nil {
		return browser.MustPage(page).MustSetUserAgent(&userAgentConfig).Timeout(1 * time.Minute)
	}

	browser = rod.New().MustConnect()

	return browser.MustPage(page).MustSetUserAgent(&userAgentConfig).Timeout(1 * time.Minute)
}

func getLinks(base *url.URL, crawlConfig crawlConfiguration, selector_config []db.ProductSelector) {
	// create a rod instance with user agent
	page := setupPage(base.String())
	defer page.Close()

	page.MustWaitLoad()

	switch {
	case len(selector_config) > 0:
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
	default:
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
}

// get all links on page
func getAllLinksOnPage(crawlConfig crawlConfiguration, selector_config []db.ProductSelector) {
	if err := sem.Acquire(semCtx, 1); err != nil {
		logrus.Warnf("failed to acquire semaphore %v\n", err)
	}

	crawlConfig.website = strings.TrimSpace(crawlConfig.website)

	// visited?
	visitMut.RLock()
	_, ok := visited[crawlConfig.website]
	visitMut.RUnlock()

	if !ok && crawlConfig.depth >= 0 {
		// mark visited
		visitMut.Lock()
		visited[crawlConfig.website] = none
		visitMut.Unlock()

		base, err := url.Parse(crawlConfig.website)
		if err != nil {
			logrus.Warnf("unable to parse website: %s\n%v", crawlConfig.website, err)
			return
		}

		tryError := rod.Try(func() {
			getLinks(base, crawlConfig, selector_config)
		})

		if errors.Is(tryError, context.DeadlineExceeded) {
			// code for timeout error
			logrus.Errorf("TIMEOUT ON %s", crawlConfig.website, tryError)
		} else if tryError != nil {
			// code for other types of error
			logrus.Errorf("ERROR ON %v %v", crawlConfig.website, tryError)
		}
	}

	sem.Release(1)
	waitgroup.Done()
}

func GetUrl(website string, depth, worker int, selector_config []db.ProductSelector) {
	// number of workers
	numWorkers := worker

	waitgroup = &sync.WaitGroup{}

	visitMut = &sync.RWMutex{}
	productMut = &sync.RWMutex{}
	allSiteMut = &sync.RWMutex{}

	config := crawlConfiguration{website: website, depth: depth}

	allSiteLink = map[string]int{}
	visited = map[string]struct{}{}

	// get root urls
	waitgroup.Add(1)
	go getAllLinksOnPage(config, selector_config)
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

			newConfig := crawlConfiguration{website: keys[i].key, depth: keys[i].depth - 1}
			waitgroup.Add(1)
			go getAllLinksOnPage(newConfig, selector_config)
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
