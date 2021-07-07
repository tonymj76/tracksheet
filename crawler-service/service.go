package main

import (
	"cache_service/crawler-service/controller"
	"cache_service/crawler-service/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//load env var
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}
	db.CreateDBConnection()
	//rest api config
	r := gin.Default()

	crawler := r.Group("/crawler")

	{
		crawler.GET("/scrap", controller.CrawlerGetController)
		crawler.GET("/url", controller.GetUrlModel)
		crawler.GET("/product", controller.GetProductUrl)
		crawler.GET("/productdelete", controller.DeleteProductUrl)
		crawler.GET("/urldelete", controller.DeleteUrlDetails)
		crawler.GET("/searchurl", controller.SeachUrl)
		crawler.GET("/searchproduct", controller.SearchProduct)
		crawler.GET("/getpaths", controller.GetPaths)
	}

	r.Run(":3001")
}
