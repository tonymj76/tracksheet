package main

import (
	"cache_service/crawler-service/controller"
	"cache_service/crawler-service/db"
	"log"

	"github.com/gocraft/work"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//load env var
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}
	redisPool := db.NewRedisPool(5, 5, true, "6379")
	db.Pool = work.NewWorkerPool(db.RedisServerContext{}, 10, "my_app", redisPool)
	// Start processing jobs
	db.Pool.Start()
	//create db connection
	db.CreateDBConnection()
	//rest api config
	r := gin.Default()

	r.GET("/crawler/scrap", controller.CrawlerGetController)
	r.GET("/crawler/url", controller.GetUrlModel)
	r.GET("/crawler/product", controller.GetProductUrl)
	r.GET("/crawler/productdelete", controller.DeleteProductUrl)
	r.GET("/crawler/urldelete", controller.DeleteUrlDetails)
	r.GET("/crawler/searchurl", controller.SeachUrl)
	r.GET("/crawler/searchproduct", controller.SearchProduct)
	r.GET("/crawler/getpaths", controller.GetPaths)

	r.Run(":3001")
}
