package main

import (
	"cache_service/crawler-service/recursive"
	// "cache_service/crawler-service/search"
	"fmt"
	// "log"
)

var name = "https://www.obeezi.com"

// var name = " https://www.amazon.com/sitemaps.2307ea63773dfee.sitemap_detail_us_txt_9_2_41_20.xml.gz"

// var name = "https://www.google.com/"

func main() {
	// links, err := search.Run(name, 0)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(links)
	urls := recursive.GetUrl("https://www.obeezi.com", 3)
	fmt.Println(urls)
}
