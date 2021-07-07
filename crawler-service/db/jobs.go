package db

import (
	"cache_service/crawler-service/handler"
	"fmt"
)

// RedisServer create a custom settings for redis
type RedisServerContext struct {
	URL         *URLModel
	ProductList handler.ReturnList
}

func SaveToDB(r *RedisServerContext) error {
	fmt.Println("saving to db...")
	err := DB.Create(r.URL).Error
	if r.ProductList.Key != "" {
		pathModel := PathMap{PathField: r.ProductList.Key, URLModelID: r.URL.ID, URLModel: r.URL}
		err = DB.Create(&pathModel).Error
	}
	for _, value := range r.ProductList.List {
		linkModel := URLProductLink{ProductURL: value, URLModelID: r.URL.ID, URLModel: r.URL}
		err = DB.Create(&linkModel).Error
	}
	return err
}

func DeleteID(r *RedisServerContext) error {
	var err error
	var productUrl URLProductLink
	if err := DB.Where("url_model_id = ?", r.URL.ID).Delete(&productUrl).Error; err != nil {
		return err
	}
	err = DB.Delete(&r.URL).Error
	return err
}
