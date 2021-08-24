package controller

import (
	"cache_service/crawler-service/db"
	"cache_service/crawler-service/handler"
	"cache_service/crawler-service/search"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// controller to crawl url input
func CrawlerGetController(ctx *gin.Context) {
	link := ctx.Query("link")
	// depth value to crawl
	depth := ctx.DefaultQuery("depth", "5")
	// number of go workers
	worker := ctx.DefaultQuery("worker", "5")
	var config []db.ProductSelector
	var other_config []db.ProductSelector

	depth_val, err := strconv.Atoi(depth)
	if err != nil {
		log.Fatal("error parsing depth to int ", err)
		return
	}
	worker_val, err := strconv.Atoi(worker)
	if err != nil {
		log.Fatal("error parsing worker to int ", err)
		return
	}

	if link == "" {
		log.Panicln("please provide a link in the query to crawl")
		return
	}
	base, err := url.Parse(link)
	if err != nil {
		log.Panicln("error parsing string to url ", err)
		return
	}
	// construct base url name
	domain := fmt.Sprintf("%s://%s", base.Scheme, base.Host)
	attr_name := "product_selector"
	// get all product selector config
	db.DB.Where("root_url = ?", &domain).Where("attr_name = ?", attr_name).Find(&config)
	// get other config that is not product selector congig
	db.DB.Where("root_url = ?", &domain).Not("attr_name = ?", attr_name).Find(&other_config)

	go func() {
		<-ctx.Request.Context().Done()
		fmt.Println("request terminated")
		os.Exit(0)
	}()

	handler.GetURL(link, depth_val, worker_val, config, other_config)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "job is complete"})
}

// controller to parse sitemap
func CrawlerGetSitemapController(context *gin.Context) {
	link := context.Query("link")
	if link == "" {
		log.Panicln("please provide a link in the query to crawl")
		return
	}

	link_details := search.Run(link)
	if link_details.Error != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "error performing job", "error": link_details.Error})
	} else {
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "job is complete"})
		// add links to db
		for _, product_url := range link_details.List {
			db.SaveToDB(product_url)
		}
	}
}

// delete product url
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

// search product url
func SearchProduct(context *gin.Context) {
	var product_url []db.URLProductLink
	var count int64
	search_query := context.Query("k")
	if search_query == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", search_query)
	// make serach case insensitive
	db.DB.Scopes(db.Paginate(context)).Where("product_url ILIKE ?", search).Find(&product_url).Count(&count)
	context.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_url,
			"total":   count,
		})
}
