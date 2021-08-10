package db

import (
	"github.com/sirupsen/logrus"
)

// RedisServer create a custom settings for redis
type Context struct {
	ProductList []string
}

//save to db
func SaveToDB(product string) {
	error := DB.Create(&URLProductLink{ProductURL: &product}).Error
	if error != nil {
		logrus.Warnln("ERROR WRITING TO DB  ", error)
	}

}

func SaveSelector(selector, root_url string) {
	attr_name := "product_selector"
	error := DB.Create(&ProductSelector{AttrValue: &selector, AttrName: &attr_name, RootURL: &root_url}).Error
	if error != nil {
		logrus.Warnln("ERROR WRITING TO DB  ", error)
	}
}
