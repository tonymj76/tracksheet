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
