package db

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// RedisServer create a custom settings for redis
type Context struct {
	ProductList []string
}

//save product url to db
func SaveToDB(product string) {
	err := DB.Create(&URLProductLink{ProductURL: &product}).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			//do nothing
		} else {
			logrus.Warnln("ERROR WRITING TO DB ", err)
		}
	}
}

//save selector config to db
func SaveSelector(selector, root_url string) {
	attr_name := "product_selector"
	err := DB.Create(&ProductSelector{AttrSelector: &selector, AttrName: &attr_name, RootURL: &root_url}).Error
	if err != nil {
		logrus.Warnln("ERROR WRITING TO DB  ", err)
	}
}

//save other attr value to db
func SaveOtherAttr(product, attr_name, attr_value string) {
	err := DB.Create(&ProductDetails{ProductUrl: &product, AttrName: &attr_name, AttrValue: &attr_value}).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			//do nothing
		} else {
			logrus.Warnln("ERROR WRITING TO DB ", err)
		}
	}
}
