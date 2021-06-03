package search

import (
	"fmt"

	colly "github.com/gocolly/colly/v2"
)

type ProductDetail struct {
	Name, Title, Description, Price, SKU, UPC string
}

// _
func Run(url string) map[string][]ProductDetail {
	// var product = make(map[string][]ProductDetail)
	var productsDetails []ProductDetail
	collyCollector := colly.NewCollector()
	collyCollector.OnHTML(".prd > .name", func(e *colly.HTMLElement) {
		productsDetails = append(productsDetails, ProductDetail{Name: e.Text})
	})
	if err := collyCollector.Visit(url); err != nil {
		fmt.Println(err)
	}
	// product[url] = productsDetails
	return map[string][]ProductDetail{
		url: productsDetails,
	}
}
