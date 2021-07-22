package search

import (
	"cache_service/crawler-service/config"
	"cache_service/crawler-service/handler"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

// Run is the entry point of the app
// input url = [url string | xml.gz file], depth = int
// depth MaxDept limits the recursion depth of visited URLs.
func Run(url string, depth ...interface{}) handler.ReturnList {
	url = strings.Trim(url, " ")
	if depth != nil && depth[0].(int) >= 1 {
		return handler.GetUrl(url, depth[0].(int))
	}
	conf := handler.NewConnect(logrus.New(), url, config.GetEnv("accept"), config.GetEnv("userAgent"))
	return Match(conf)
}
