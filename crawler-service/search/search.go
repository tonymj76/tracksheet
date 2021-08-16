package search

import (
	"cache_service/crawler-service/config"
	"cache_service/crawler-service/handler"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

// Run is the entry point of the app
// input url = [.xml | .xml.gz file],
func Run(url string) handler.ReturnList {
	url = strings.TrimSpace(url)
	conf := handler.NewConnect(logrus.New(), url, config.GetEnv("accept"), config.GetEnv("userAgent"))
	return Match(conf)
}
