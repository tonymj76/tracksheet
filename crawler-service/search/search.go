package search

import (
	"cache_service/crawler-service/db"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

// Run is the entry point of the app
// input url = [url string | xml.gz file], depth = int
// depth MaxDept limits the recursion depth of visited URLs.
func Run(url, user_agent, accept string, depth ...interface{}) ([]string, error) {
	url = strings.Trim(url, " ")
	productLinks := []string{}
	if depth != nil && depth[0].(int) > 1 {
		dp := depth[0].(int)
		productLinks = append(productLinks, links...)

		if err != nil {
			logrus.WithError(err).Fatalln("failed to initiate the db")
		}

	}

	database, err := db.CreateDBConnection()

	if err != nil {
		logrus.WithError(err).Fatalln("failed to initiate the db")
	}

	if len(productLinks) > 0 {
		urlLinks := db.URLProductLinks{}
		for _, link := range productLinks {
			urlLinks = append(urlLinks, db.URLProductLink{ProductURL: link})
		}
		err = database.Create(db.URLModel{Name: url, ProductURLs: urlLinks}).Error
	}
	return productLinks, err
}
