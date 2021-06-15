package search

import (
	"cache_service/crawler-service/structure"
	"regexp"
	"strings"
)

// _
func Run(url string) ([]string, error) {
	url = strings.Trim(url, " ")
	validString := regexp.MustCompile(`(?P<head>http[s]?:\/\/)(?P<host>[^\/\s]+\/)(?P<path>[^\s]*\.gz$)`)
	if validString.MatchString(url) {
		var sitemapProject *structure.GzProduct
		return Match(sitemapProject, url)
	} else {
		var product *structure.ProductDescription
		return Match(product, url)
	}
}
