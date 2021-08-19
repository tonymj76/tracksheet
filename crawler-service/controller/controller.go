package controller

import (
	"cache_service/crawler-service/db"
	"cache_service/crawler-service/handler"
	"cache_service/crawler-service/search"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SelectorConfig struct {
	AttrName     string `json:"attrname"`
	AttrSelector string `json:"attrselector"`
	RootURL      string `json:"rooturl"`
}

func CrawlerGetController(context *gin.Context) {
	link := context.Query("link")
	depth := context.DefaultQuery("depth", "5")
	worker := context.DefaultQuery("worker", "5")
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

	domain := fmt.Sprintf("%s://%s", base.Scheme, base.Host)
	attr_name := "product_selector"
	db.DB.Where("root_url = ?", &domain).Where("attr_name = ?", attr_name).Find(&config)
	db.DB.Where("root_url = ?", &domain).Not("attr_name = ?", attr_name).Find(&other_config)

	handler.GetURL(link, depth_val, worker_val, config, other_config)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "job is complete"})

}

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
		//add links to db
		for _, product_url := range link_details.List {
			db.SaveToDB(product_url)
		}
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
	db.DB.Scopes(db.Paginate(context)).Where("product_url ILIKE ?", search).Find(&product_url).Count(&count)
	context.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_url,
			"total":   count,
		})
}

func AddSelectorConfig(context *gin.Context) {
	var input SelectorConfig
	var selector_config string

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if strings.HasPrefix(input.AttrSelector, ".") {
		selector_config = input.AttrSelector
		selector_config = strings.ReplaceAll(selector_config, " ", ".")
	} else {
		selector_config = fmt.Sprintf(".%v", input.AttrSelector)
		selector_config = strings.ReplaceAll(selector_config, " ", ".")
	}
	//root url value should not empty
	if input.RootURL == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "root url cannot be blank"})
		return
	}

	//if a root url and attr name does not exist for a value, add it
	selector := db.ProductSelector{
		AttrName:     &input.AttrName,
		AttrSelector: &selector_config,
		RootURL:      &input.RootURL,
	}
	error := db.DB.Create(&selector).Error

	if error != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "error saving to DB", "error": error})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "selector config saved to DB"})
		return
	}
}

func GetSelectorConfig(context *gin.Context) {
	var config []db.ProductSelector
	link := context.Query("url")

	if link == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "url query cannot be blank"})
		return
	}
	db.DB.Where("root_url = ?", link).Find(&config)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": config})

}

func DeleteSelectorConfig(context *gin.Context) {
	var config db.ProductSelector
	var config2 []db.ProductSelector

	link := context.Query("url")
	attr_name := context.Query("attr_name")

	if link == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "url query or attr name cannot be blank"})
		return
	}
	if attr_name == "" {
		if err := db.DB.Where("root_url = ?", link).Where("attr_name = ?", attr_name).First(&config).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
			return
		}
		db.DB.Delete(&config)
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "delete config successful"})
	} else {
		db.DB.Where("root_url = ?", link).Find(&config2)
		db.DB.Delete(&config)
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "delete config successful"})
	}
}

func GetProductDetails(context *gin.Context) {
	var product_details []db.ProductDetails
	var count int64
	search_query := context.Query("k")
	if search_query == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", search_query)
	db.DB.Scopes(db.Paginate(context)).Where("product_url ILIKE ?", search).Or("attr_value ILIKE ?", search).Find(&product_details).Count(&count)
	context.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_details,
			"total":   count,
		})
}

func DeleteProductDetails(context *gin.Context) {
	var product_details db.ProductDetails
	id := context.Param("id")
	if err := db.DB.Where("id = ?", id).First(&product_details).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Record not found!", "success": false})
		return
	}
	db.DB.Delete(&product_details)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "delete deltails successful"})
}
