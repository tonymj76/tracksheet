package main

import (
	"cache_service/crawler-service/recursive"
	"cache_service/crawler-service/search"
	"fmt"
	"log"
)

// var name = "https://www.obeezi.com"
var name = " https://www.amazon.com/sitemaps.2307ea63773dfee.sitemap_detail_us_txt_9_2_41_20.xml.gz"

func main() {
	data := recursive.GetUrl("https://www.google.com/", 2)
	fmt.Println(data)
	links, err := search.Run(name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(links)
}
