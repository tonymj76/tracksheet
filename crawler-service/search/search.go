package search

import (
	_ "cache_service/crawler-service/log"
)

// Run is the entry point of the app
// input url = [url string | xml.gz file], depth = int
// depth MaxDept limits the recursion depth of visited URLs.
func Run(url string, depth ...interface{}) {
	// url = strings.Trim(url, " ")
	// if depth != nil && depth[0].(int) >= 1 {
	// 	return handler.GetUrl(url, depth[0].(int))
	// }
	// conf := handler.NewConnect(logrus.New(), url, config.GetEnv("accept"), config.GetEnv("userAgent"))
	// return Match(conf)
}
