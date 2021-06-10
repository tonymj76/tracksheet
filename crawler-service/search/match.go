package search

import (
	"fmt"
	"log"
)

type Matcher interface {
	Search(url string) ([]string, error)
}

func Match(matcher Matcher, url string) {
	searchResults, err := matcher.Search(url)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(searchResults)
}
