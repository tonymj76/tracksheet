package db

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// RedisServer create a custom settings for redis
type Context struct {
	ProductList []string
}

//save to db
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

func SaveSelector(selector, root_url string) {
	attr_name := "product_selector"
	err := DB.Create(&ProductSelector{AttrSelector: &selector, AttrName: &attr_name, RootURL: &root_url}).Error
	if err != nil {
		logrus.Warnln("ERROR WRITING TO DB  ", err)
	}
}

func SaveOtherAttr(product, attr_name, attr_value string) {
	var product_detail []ProductDetails

	if err := DB.Where("product_url = ?", &product).Where("attr_name = ?", attr_name).Where("attr_value = ?", attr_value).First(&product_detail); err != nil {
		errs := DB.Create(&ProductDetails{ProductUrl: &product, AttrName: &attr_name, AttrValue: &attr_value}).Error
		if errs != nil {
			logrus.Warnln("ERROR WRITING TO DB ", err.Error)
		}
	}

}
