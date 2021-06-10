package search

import (
	"cache_service/crawler-service/structure"
	"regexp"
	"strings"
)

// _
func Run(url string) {
	url = strings.Trim(url, " ")
	validString := regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.xml$)|(?P<path2>[^\s]*\.gz$)`)
	if validString.MatchString(url) {
		var sitemapProject *structure.GzProduct
		Match(sitemapProject, url)
	} else {
		var product *structure.ProductDescription
		Match(product, url)
	}
}
