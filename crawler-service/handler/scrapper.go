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
	"golang.org/x/sync/semaphore"
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

func getLinks(base *url.URL, crawlConfig crawlConfiguration, selectorConfig []db.ProductSelector) {
	// create a rod instance with user agent
	page := setupPage(base.String())
	defer page.Close()

	page.MustWaitLoad()

	switch {
	case len(selectorConfig) > 0:
		queryWithSelector(page, base, crawlConfig, selectorConfig)
	default:
		queryRecursively(page, base, crawlConfig, selectorConfig)

	}
}

// get all links on page
func getAllLinksOnPage(crawlConfig crawlConfiguration, selectorConfig []db.ProductSelector) {
	if err := sem.Acquire(semCtx, 1); err != nil {
		logrus.Warnf("failed to acquire semaphore %v\n", err)
	}

	crawlConfig.website = strings.TrimSpace(crawlConfig.website)

	// visited?
	_, ok := visited.Load(crawlConfig.website)

	if !ok && crawlConfig.depth >= 0 {
		// mark visited
		visited.Store(crawlConfig.website, none)

		base, err := url.Parse(crawlConfig.website)
		if err != nil {
			logrus.Warnf("unable to parse website: %s\n%v", crawlConfig.website, err)
			return
		}

		tryError := rod.Try(func() {
			getLinks(base, crawlConfig, selectorConfig)
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

// GetURL gets the links in a URL
func GetURL(website string, depth, worker int, selectorConfig []db.ProductSelector) {
	// number of workers
	numWorkers := worker

	waitgroup = &sync.WaitGroup{}

	productMut = &sync.RWMutex{}

	semCtx = context.TODO()
	sem = semaphore.NewWeighted(int64(numWorkers))

	config := crawlConfiguration{website: website, depth: depth}

	// get root urls
	waitgroup.Add(1)
	go getAllLinksOnPage(config, selectorConfig)
	waitgroup.Wait()

	type kvPair struct {
		key   string
		depth int
	}

	keys := []kvPair{}

	allSiteLink.Range(func(k, v interface{}) bool {
		if _, ok := visited.Load(k); !ok {
			keys = append(keys, kvPair{k.(string), v.(int)})
		}
		return true
	})

	for len(keys) > 0 {
		size := minInt(len(keys), numWorkers)

		for i := 0; i < size; i++ {
			depth := keys[i].depth
			fmt.Println("DEPTH::::::::", depth)

			allSiteLink.Delete(keys[i].key)

			newConfig := crawlConfiguration{website: keys[i].key, depth: keys[i].depth - 1}
			waitgroup.Add(1)
			go getAllLinksOnPage(newConfig, selectorConfig)
		}

		waitgroup.Wait()

		keys = keys[size:]

		allSiteLink.Range(func(k, v interface{}) bool {
			if _, ok := visited.Load(k); !ok {
				keys = append(keys, kvPair{key: k.(string), depth: v.(int) - 1})
				allSiteLink.Delete(k)
			}
			return true
		})
	}

	fmt.Println("---", len(keys))
}

func minInt(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
