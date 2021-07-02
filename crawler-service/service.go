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
	//create db connection
	db.CreateDBConnection()
	//rest api config
	r := gin.Default()

	r.GET("/crawler/scrap", controller.CrawlerGetController)
	r.GET("/crawler/url", controller.GetUrlModel)
	r.GET("/crawler/product", controller.GetProductUrl)
	r.GET("/crawler/productdelete", controller.DeleteProductUrl)
	r.GET("/crawler/urldelete", controller.DeleteUrlDetails)

	r.Run(":3001")
}
