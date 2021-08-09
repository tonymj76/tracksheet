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

type SelectorConfig struct {
	ProductSelector      string `json:"productselector"`
	ReatilerNameSelector string `json:"retailnameselector"`
	TitleSelector        string `json:"titleselector"`
	DescriptionSelector  string `json:"descriptionselector"`
	ImageSelector        string `json:"imageselector"`
	UpcSelector          string `json:"upcselector"`
	SkuSelector          string `json:"skuselector"`
	BrandSelector        string `json:"brandselector"`
	CategorySelector     string `json:"categoryselector"`
	RootURL              string `json:"rooturl"`
}

func CrawlerGetController(context *gin.Context) {
	link := context.Query("link")
	depth := context.DefaultQuery("depth", "5")
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
	db.DB.Scopes(db.Paginate(context)).Where("product_url LIKE ?", search).Find(&product_url).Count(&count)
	context.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_url,
			"total":   count,
		})
}

func AddSelectorConfig(context *gin.Context) {
	var config SelectorConfig
	if err := context.ShouldBindJSON(&config); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if config.RootURL == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "root url cannot be blank"})
		return
	}
	selector := db.ProductSelector{
		ProductSelector:      &config.ProductSelector,
		ReatilerNameSelector: &config.ReatilerNameSelector,
		TitleSelector:        &config.TitleSelector,
		DescriptionSelector:  &config.DescriptionSelector,
		ImageSelector:        &config.ImageSelector,
		UpcSelector:          &config.UpcSelector,
		SkuSelector:          &config.SkuSelector,
		BrandSelector:        &config.BrandSelector,
		CategorySelector:     &config.CategorySelector,
		RootURL:              &config.RootURL,
	}
	error := db.DB.Create(&selector).Error

	if error != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "error saving to DB", "error": error})
	} else {
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "selector config saved to DB"})
	}
}

func UpdateSelectorConfig(context *gin.Context) {
	var config db.ProductSelector
	var input SelectorConfig
	link := context.Query("url")

	if link == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "url query cannot be blank"})
		return
	}

	if err := db.DB.Where("root_url = ?", link).First(&config).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	selector := db.ProductSelector{
		ProductSelector:      &input.ProductSelector,
		ReatilerNameSelector: &input.ReatilerNameSelector,
		TitleSelector:        &input.TitleSelector,
		DescriptionSelector:  &input.DescriptionSelector,
		ImageSelector:        &input.ImageSelector,
		UpcSelector:          &input.UpcSelector,
		SkuSelector:          &input.SkuSelector,
		BrandSelector:        &input.BrandSelector,
		CategorySelector:     &input.CategorySelector,
	}
	error := db.DB.Model(&config).Updates(selector).Error

	if error != nil {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "error updating field", "error": error})
	} else {
		context.JSON(http.StatusOK, gin.H{"success": true, "message": "selector config updated"})
	}
}

func GetSelectorConfig(context *gin.Context) {
	var config db.ProductSelector
	link := context.Query("url")

	if link == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "url query cannot be blank"})
		return
	}
	if err := db.DB.Where("root_url = ?", link).First(&config).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"success": true, "message": config})

}

func DeleteSelectorConfig(context *gin.Context) {
	var config db.ProductSelector
	link := context.Query("url")

	if link == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "url query cannot be blank"})
		return
	}
	if err := db.DB.Where("root_url = ?", link).First(&config).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Record not found!"})
		return
	}
	db.DB.Delete(&config)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "delete config successful"})
}
