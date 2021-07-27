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
		context.JSON(http.StatusOK, gin.H{"success": true, "data": links.List, "root_url": link})
		// after responsing with json, run db update in the background
		urlModel := db.RootUrl{Name: &link}
		items := &db.Context{URL: urlModel, ProductList: links}
		db.SaveToDB(items)
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
		context.JSON(http.StatusOK, gin.H{"success": true, "data": links.List, "root_url": input.Link})
		// after responsing with json, run db update in the background
		urlModel := db.RootUrl{Name: &input.Link}
		items := &db.Context{URL: urlModel, ProductList: links}
		db.SaveToDB(items)
	}
}

func GetRootUrlModel(context *gin.Context) {
	// get all items from url model
	var urls []db.RootUrl
	var count int64

	db.DB.Scopes(db.Paginate(context)).Find(&urls).Count(&count)

	context.JSON(http.StatusOK, gin.H{"success": true, "root_url": urls, "total": count})
}

func GetRootUrlById(context *gin.Context) {
	var root_url db.RootUrl
	if err := db.DB.Where("id = ?", context.Param("id")).First(&root_url).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"success": true, "root_url": root_url})
}

func GetProductUrl(context *gin.Context) {
	var url_detail db.RootUrl
	var product_url []db.URLProductLink
	var count int64

	url := context.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("name = ?", url).First(&url_detail).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}

	db.DB.Scopes(db.Paginate(context)).Where("root_url_id = ?", url_detail.ID).Find(&product_url).Count(&count)

	context.JSON(http.StatusOK,
		gin.H{
			"success":  true,
			"root_url": url_detail,
			"product":  product_url,
			"total":    count,
		})
}

func DeleteUrlDetails(context *gin.Context) {
	var url_detail db.RootUrl
	url := context.Query("link")
	if url == "" {
		fmt.Println("Please provide a url in query")
		return
	}

	if err := db.DB.Where("name = ?", url).First(&url_detail).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "item deleted successfully"})

	items := &db.Context{URL: url_detail}
	db.DeleteID(items)
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

func SeachUrl(context *gin.Context) {
	var url_detail []db.RootUrl
	k := context.Query("k")
	var count int64
	if k == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", k)
	db.DB.Scopes(db.Paginate(context)).Where("name LIKE ?", search).Find(&url_detail).Count(&count)

	context.JSON(http.StatusOK,
		gin.H{
			"success":  true,
			"root_url": url_detail,
			"total":    count,
		})
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
