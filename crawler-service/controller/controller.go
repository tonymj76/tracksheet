package controller

import (
	"cache_service/crawler-service/db"
	"cache_service/crawler-service/search"
	"fmt"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"

func CrawlerGetController(c *gin.Context) {
	//accept query of user agent, accept, link and depth
	accept := c.DefaultQuery("accept", accept)
	user_agent := c.DefaultQuery("user_agent", userAgent)
	link := c.Query("link")
	depth := c.DefaultQuery("depth", "0")
	depth_val, err := strconv.Atoi(depth)
	if err != nil {
		log.Fatal("error parsing depth to int ", err)
	}

	if link == "" {
		log.Fatal("please provide a link in the query to crawl")
		return
	}

	links := search.Run(link, user_agent, accept, depth_val)
	if links.Error != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "links": links.List})
		//after responsing with json, run db update in the background
		urlModel := db.URLModel{Name: &link}
		items := &db.RedisServerContext{URL: &urlModel, ProductList: links}
		db.SaveToDB(items)
	}
}

func GetUrlModel(c *gin.Context) {
	//get all items from url model
	var urls []db.URLModel
	var count int64

	db.DB.Scopes(db.Paginate(c)).Find(&urls).Count(&count)

	c.JSON(http.StatusOK, gin.H{"success": true, "url": urls, "total": count})

}

func GetPaths(c *gin.Context) {
	var paths []db.PathMap
	var count int64
	db.DB.Scopes(db.Paginate(c)).Find(&paths).Count(&count)
	c.JSON(http.StatusOK, gin.H{"success": true, "path": paths, "total": count})

}

func GetProductUrl(c *gin.Context) {
	var url_detail db.URLModel
	var product_url []db.URLProductLink
	var count int64

	url := c.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("name = ?", url).First(&url_detail).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}

	db.DB.Scopes(db.Paginate(c)).Where("url_model_id = ?", url_detail.ID).Find(&product_url).Count(&count)

	c.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"url":     url_detail,
			"product": product_url,
			"total":   count,
		})

}

func DeleteUrlDetails(c *gin.Context) {
	var url_detail db.URLModel
	url := c.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("name = ?", url).First(&url_detail).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "item deleted successfully"})

	items := &db.RedisServerContext{URL: &url_detail}
	db.DeleteID(items)
}

func DeleteProductUrl(c *gin.Context) {
	var product_url db.URLProductLink

	url := c.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("product_url = ?", url).First(&product_url).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}

	db.DB.Delete(&product_url)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "item deleted successfully"})
}

func SeachUrl(c *gin.Context) {
	var url_detail []db.URLModel
	k := c.Query("k")
	var count int64
	if k == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", k)
	db.DB.Scopes(db.Paginate(c)).Where("name LIKE ?", search).Find(&url_detail).Count(&count)

	c.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"url":     url_detail,
			"total":   count,
		})
}

func SearchProduct(c *gin.Context) {
	var product_url []db.URLProductLink
	var count int64
	k := c.Query("k")
	if k == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", k)
	db.DB.Scopes(db.Paginate(c)).Where("product_url LIKE ?", search).Find(&product_url).Count(&count)
	c.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_url,
			"total":   count,
		})
}
