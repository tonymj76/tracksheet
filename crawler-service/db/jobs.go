package db

import (
	"fmt"
)

// RedisServer create a custom settings for redis
type Context struct {
	ProductList []string
}

//save to db
func SaveToDB(product string) {
	error := DB.Create(&URLProductLink{ProductURL: &product}).Error
	if error != nil {
		fmt.Println("error writing to db ", error)
	}

}
