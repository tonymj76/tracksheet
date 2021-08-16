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
	var product_item URLProductLink
	if err := DB.Where("product_url = ?", product).First(&product_item).Error; err != nil {
		error := DB.Create(&URLProductLink{ProductURL: &product}).Error
		if error != nil {
			logrus.Warnln("ERROR WRITING TO DB  ", error)
		}
	}
}

func SaveSelector(selector, root_url string) {
	attr_name := "product_selector"
	error := DB.Create(&ProductSelector{AttrValue: &selector, AttrName: &attr_name, RootURL: &root_url}).Error
	if error != nil {
		logrus.Warnln("ERROR WRITING TO DB  ", error)
	}
}
