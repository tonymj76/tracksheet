package main

import (
	"cache_service/crawler-service/search"
)

// var name = "https://www.vitacost.com"

// var name = "https://www.konga.com"

var name = "https://www.obeezi.com/robots.txt"

func main() {
	search.Run(name)
}
