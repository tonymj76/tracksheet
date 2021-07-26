package main

import (
	"cache_service/crawler-service/config"
	"cache_service/crawler-service/controller"
	"cache_service/crawler-service/db"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func main() {
	// load env var
	config.LoadConfigFromJSON()

	db.CreateDBConnection()
	cron_job()
	// rest api config
	ginConfig := gin.Default()

	crawler := ginConfig.Group("/crawler")

	{
		crawler.GET("/scrap", controller.CrawlerGetController)
		crawler.GET("/domains", controller.GetRootUrlModel)
		crawler.GET("/product", controller.GetProductUrl)
		crawler.GET("/productdelete", controller.DeleteProductUrl)
		crawler.GET("/urldelete", controller.DeleteUrlDetails)
		crawler.GET("/searchdomains", controller.SeachUrl)
		crawler.GET("/searchproduct", controller.SearchProduct)
		crawler.POST("/scrap", controller.CrawlerPostController)
		crawler.GET("/domain/:id", controller.GetRootUrlById)
	}

	ginConfig.Run(config.GetEnv("PORT"))
}

func cron_job() {
	// cron job to clean cache every 24 hours
	cron_instance := cron.New()
	cron_instance.AddFunc("@daily", func() {
		log.Println("running daily job")
		os.RemoveAll("cache/")
	})
	cron_instance.Start()
}
