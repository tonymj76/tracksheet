package main

import (
	"cache_service/search"
	"fmt"
)

func main() {
	fmt.Println(search.Run("https://www.vitacost.com/"))
}
