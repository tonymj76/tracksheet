package search

import "cache_service/structure"

// _
func Run(url string) {
	var product *structure.ProductDescription
	Match(product, url)
}
