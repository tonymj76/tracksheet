package handler

import "sync"

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
