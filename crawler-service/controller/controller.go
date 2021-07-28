package controller

import (
	"cache_service/crawler-service/db"
	"cache_service/crawler-service/handler"
	"cache_service/crawler-service/search"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Input struct {
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
	Str3  string `json:"str3"`
	Depth int    `json:"depth"`
	Link  string `json:"link"`
}

func CrawlerGetController(context *gin.Context) {
	link := context.Query("link")
	depth := context.DefaultQuery("depth", "0")
	depth_val, err := strconv.Atoi(depth)
	if err != nil {
		log.Fatal("error parsing depth to int ", err)
	}

	if link == "" {
		log.Panicln("please provide a link in the query to crawl")
		return
	}

	links := search.Run(link, depth_val)
	if links.Error != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "error": err})
	} else {
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "job is complete"})
		// after responsing with json, run db update in the background
		// items := &db.Context{ProductList: links.List}
		// db.SaveToDB(items)
	}
}

func CrawlerPostController(context *gin.Context) {
	var input Input

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Link == "" {
		log.Println("please provide a link in the query to crawl")
		return
	}
	//get sample product link input
	sample_input := handler.SampleLink{
		Str1: input.Str1,
		Str2: input.Str2,
		Str3: input.Str3,
	}
	links := handler.GetProductUrl(input.Link, input.Depth, sample_input)
	if links.Error != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "error": links.Error})
	} else {
		//truncate list to avoid browser throwing error of too much payload
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "job is complete"})
		// after responsing with json, run db update in the background
		// items := &db.Context{ProductList: links.List}
		// db.SaveToDB(items)
	}
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
