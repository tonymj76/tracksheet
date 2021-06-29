package main

import (
	"cache_service/crawler-service/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/crawler", controller.CrawlerGetController)
	r.Run(":3001")
}
