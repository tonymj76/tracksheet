package search

type ProductDetail struct {
	Name,
	Title,
	Description,
	Price,
	SKU,
	UPC string
}

// _
func Run(url string) map[string][]ProductDetail {
	var product = make(map[string][]ProductDetail)
	return product
}
