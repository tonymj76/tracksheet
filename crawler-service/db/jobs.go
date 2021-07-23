package db

import (
	"cache_service/crawler-service/handler"
	"fmt"
)

// RedisServer create a custom settings for redis
type Context struct {
	URL         RootUrl
	ProductList handler.ReturnList
}

//save to db
func SaveToDB(r *Context) {
	fmt.Println("saving to db...")
	var root_url RootUrl
	error := DB.Create(&r.URL).Error
	if error != nil {
		//if url already exist, use it to ref other product url
		if err := DB.Where("name = ?", r.URL.Name).First(&root_url).Error; err != nil {
			return
		}
		for _, value := range r.ProductList.List {
			linkModel := URLProductLink{ProductURL: &value, RootUrlId: root_url.ID}
			err := DB.Create(&linkModel).Error
			if err != nil {
				fmt.Println("error writing product link to db ", value, err)
			}
		}
	} else {
		//if url already exist, use it to ref other product url
		if err := DB.Where("name = ?", r.URL.Name).First(&root_url).Error; err != nil {
			return
		}
		for _, value := range r.ProductList.List {
			linkModel := URLProductLink{ProductURL: &value, RootUrlId: root_url.ID}
			err := DB.Create(&linkModel).Error
			if err != nil {
				fmt.Println("error writing product link to db ", value, err)
			}
		}
	}

}

func DeleteID(r *Context) error {
	var err error
	var productUrl URLProductLink
	if err := DB.Where("root_url_id = ?", r.URL.ID).Delete(&productUrl).Error; err != nil {
		return err
	}
	err = DB.Delete(&r.URL).Error
	return err
}
