package search

import (
	"cache_service/crawler-service/handler"
)

type Matcher interface {
	Search() handler.ReturnList
}

func Match(matcher Matcher) handler.ReturnList {
	return matcher.Search()
}
