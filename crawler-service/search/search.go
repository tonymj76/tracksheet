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
func Run(url string) handler.ReturnList {
	url = strings.Trim(url, " ")

	conf := handler.NewConnect(logrus.New(), url, config.GetEnv("accept"), config.GetEnv("userAgent"))
	return Match(conf)
}
