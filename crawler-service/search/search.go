package search

import (
	"cache_service/crawler-service/handler"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

// Run is the entry point of the app
// input url = [url string | xml.gz file], depth = int
// depth MaxDept limits the recursion depth of visited URLs.
func Run(url, user_agent, accept string, depth ...interface{}) handler.ReturnList {
	url = strings.Trim(url, " ")
	if depth != nil && depth[0].(int) >= 1 {
		return handler.GetUrl(url, accept, user_agent, depth[0].(int))
	}
	conf := handler.NewConnect(logrus.New(), url, accept, user_agent)
	return Match(conf)
}
