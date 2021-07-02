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

	links, err := search.Run(link, user_agent, accept, depth_val)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "links": links})
		//after responsing with json, run db update in the background
		//todo: make this run as a background job
		url_model := db.URLModel{Name: link}
		db.DB.Create(&url_model)
		for _, value := range links {
			link_model := db.URLProductLink{ProductURL: value, URLModelID: url_model.ID}
			db.DB.Create(&link_model)
		}

	}
}

func GetUrlModel(c *gin.Context) {
	//get all items from url model
	var urls []db.URLModel
	var count int64

	db.DB.Scopes(db.Paginate(c)).Find(&urls).Count(&count)

	c.JSON(http.StatusOK, gin.H{"success": true, "url": urls, "total": count})

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
			"success": "true",
			"url":     url_detail,
			"product": product_url,
			"total":   count,
		})

}

func DeleteUrlDetails(c *gin.Context) {
	var url_detail db.URLModel
	var product_url db.URLProductLink

	url := c.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("name = ?", url).First(&url_detail).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}
	//todo: make delete  run as a background job

	db.DB.Where("url_model_id = ?", url_detail.ID).Delete(&product_url)
	db.DB.Delete(&url_detail)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "item deleted successfully"})
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
