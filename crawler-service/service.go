package main

import (
	"cache_service/crawler-service/search"
	"fmt"

	"github.com/sirupsen/logrus"
)

var name = "https://www.obeezi.com"

// var name = " https://www.amazon.com/sitemaps.2307ea63773dfee.sitemap_detail_us_txt_9_2_41_20.xml.gz"

// var name = "https://www.google.com/"

func main() {
	links, err := search.Run(name)
	if err != nil {
		logrus.WithError(err).Fatalln(err)
	}
	fmt.Println(links)

}
