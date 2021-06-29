package search

import (
	"cache_service/crawler-service/handler"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"

// Run is the entry point of the app
// input url = [url string | xml.gz file], depth = int
func Run(url string, depth ...interface{}) ([]string, error) {
	url = strings.Trim(url, " ")
	if depth != nil {
		dp := depth[0].(int)
		if dp > 1 {
			return handler.GetUrl(url, dp)
		}
	}
	conf := handler.NewConnect(logrus.New(), url, accept, userAgent)
	return Match(conf)
}
