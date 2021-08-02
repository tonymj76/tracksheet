package controller

import (
	"cache_service/crawler-service/db"
	"cache_service/crawler-service/handler"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CrawlerGetController(context *gin.Context) {
	link := context.Query("link")
	depth := context.DefaultQuery("depth", "0")
	worker := context.DefaultQuery("worker", "5")
	depth_val, err := strconv.Atoi(depth)
	if err != nil {
		log.Fatal("error parsing depth to int ", err)
	}
	worker_val, err := strconv.Atoi(worker)
	if err != nil {
		log.Fatal("error parsing worker to int ", err)
	}

	if link == "" {
		log.Panicln("please provide a link in the query to crawl")
		return
	}

	handler.GetUrl(link, depth_val, worker_val)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "job is complete"})

}

func DeleteProductUrl(context *gin.Context) {
	var product_url db.URLProductLink

	url := context.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("product_url = ?", url).First(&product_url).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}

	db.DB.Delete(&product_url)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "item deleted successfully"})
}

func SearchProduct(context *gin.Context) {
	var product_url []db.URLProductLink
	var count int64
	search_query := context.Query("k")
	if search_query == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", search_query)
	db.DB.Scopes(db.Paginate(context)).Where("product_url LIKE ?", search).Find(&product_url).Count(&count)
	context.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_url,
			"total":   count,
		})
}
