package search

import (
	"cache_service/crawler-service/handler"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

var accept string = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
var userAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"

// Run is the entry point of the app
func Run(url string) ([]string, error) {
	url = strings.Trim(url, " ")
	conf := handler.NewConnect(logrus.New(), url, accept, userAgent)
	return Match(conf)
}
