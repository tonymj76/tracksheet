package main

import (
	"cache_service/crawler-service/config"
	"cache_service/crawler-service/controller"
	"cache_service/crawler-service/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// load env var
	config.LoadConfigFromJSON()

	db.CreateDBConnection()
	// rest api config
	ginConfig := gin.Default()

	crawler := ginConfig.Group("/crawler")

	{
		crawler.GET("/scrap", controller.CrawlerGetController)
		crawler.GET("/productdelete", controller.DeleteProductUrl)
		crawler.GET("/searchproduct", controller.SearchProduct)
		crawler.GET("/scrapsitemap", controller.CrawlerGetSitemapController)
		crawler.POST("/addconfig", controller.AddSelectorConfig)
		crawler.GET("/getconfig", controller.GetSelectorConfig)
		crawler.GET("/deleteconfig", controller.DeleteSelectorConfig)
		crawler.GET("/productdetails", controller.GetProductDetails)
		crawler.GET("/deleteproductdetail/:id", controller.DeleteProductDetails)
	}

	ginConfig.Run(config.GetEnv("PORT"))
}
