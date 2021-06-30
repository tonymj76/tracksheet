package search

import (
	"cache_service/crawler-service/db"
	"cache_service/crawler-service/handler"
	_ "cache_service/crawler-service/log"
	"strings"

	"github.com/sirupsen/logrus"
)

// Run is the entry point of the app
// input url = [url string | xml.gz file], depth = int
// depth MaxDept limits the recursion depth of visited URLs.
var (
	host     string
	port     string
	user     string
	dbname   string
	password string
	sslmode  = "disable"
)

func Run(url, user_agent, accept string, depth ...interface{}) ([]string, error) {
	url = strings.Trim(url, " ")
	productLinks := []string{}
	if depth != nil && depth[0].(int) > 1 {
		dp := depth[0].(int)
		links, err := handler.GetUrl(url, dp)
		productLinks = append(productLinks, links...)

		if err != nil {
			logrus.WithError(err).Fatalln("failed to recursively fetch urls")
		}
	} else {
		conf := handler.NewConnect(logrus.New(), url, accept, user_agent)
		links, err := Match(conf)
		productLinks = append(productLinks, links...)
		if err != nil {
			logrus.WithError(err).Fatalln("failed to initiate the db")
		}

	}

	database, err := db.CreateDBConnection(host, port, user, dbname, password, sslmode)

	if err != nil {
		logrus.WithError(err).Fatalln("failed to initiate the db")
	}

	err = database.Create(db.URLModel{Name: url, ProductURLs: productLinks}).Error
	return productLinks, err
}
