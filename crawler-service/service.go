package main

import (
	"cache_service/crawler-service/controller"
	"cache_service/crawler-service/db"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func main() {
	//load env var
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}
	db.CreateDBConnection()
	cron_job()
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
		crawler.POST("/scrap", controller.CrawlerPostController)
	}

	r.Run(":3001")
}

func cron_job() {
	//cron job to clean cache every 24 hours
	c := cron.New()
	c.AddFunc("@daily", func() {
		log.Println("running daily job")
		os.RemoveAll("cache/")
	})
	c.Start()
}
