package controller

import (
	"cache_service/crawler-service/search"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"

func CrawlerGetController(c *gin.Context) {
	//accept query of user agent, accept, link and depth
	accept := c.DefaultQuery("accept", accept)
	user_agent := c.DefaultQuery("user_agent", userAgent)
	link := c.Query("link")
	depth := c.DefaultQuery("depth", "0")
	depth_val, err := strconv.Atoi(depth)
	if err != nil {
		log.Fatal("error parsing depth to int ", err)
	}

	if link == "" {
		log.Fatal("please provide a link in the query to crawl")
		return
	}

	links, err := search.Run(link, user_agent, accept, depth_val)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": "false", "error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "true", "links": links})
	}
}
