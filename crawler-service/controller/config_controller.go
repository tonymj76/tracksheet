package controller

import (
	"cache_service/crawler-service/db"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AttrDetails struct {
	AttrName     string `json:"attrname"`
	AttrSelector string `json:"attrselector"`
}

type SelectorConfig struct {
	RootURL     string        `json:"rooturl"`
	AttrDetails []AttrDetails `json:"attrdetails"`
}

//add class selector config to db
func AddSelectorConfig(context *gin.Context) {
	var input SelectorConfig
	var selector_config string

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	//root url value should not empty
	if input.RootURL == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "root url cannot be blank"})
		return
	}

	for _, selectors := range input.AttrDetails {
		//parse string input to be used as a class selector
		if strings.HasPrefix(selectors.AttrSelector, ".") {
			selector_config = selectors.AttrSelector
			selector_config = strings.ReplaceAll(selector_config, " ", ".")
		} else {
			selector_config = fmt.Sprintf(".%v", selectors.AttrSelector)
			selector_config = strings.ReplaceAll(selector_config, " ", ".")
		}
		//if a root url and attr name does not exist for a value, add it
		selector := db.ProductSelector{
			AttrName:     &selectors.AttrName,
			AttrSelector: &selector_config,
			RootURL:      &input.RootURL,
		}
		error := db.DB.Create(&selector).Error
		if error != nil {
			logrus.Warnln("ERROR WRITING TO DB ", error)
		}

	}

	context.JSON(http.StatusOK, gin.H{"success": true, "message": "selector config saved to DB"})
}

//get selector config for url
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

//delete selector config for url
func DeleteSelectorConfig(context *gin.Context) {
	var config db.ProductSelector
	var config2 []db.ProductSelector

	link := context.Query("url")
	attr_name := context.Query("attr_name")

	if link == "" {
		context.JSON(http.StatusOK, gin.H{"success": false, "message": "url query or attr name cannot be blank"})
		return
	}
	//remove whitespace
	attr_name = strings.TrimSpace(attr_name)

	//if attr name is empty delete all config with url else delete single attr
	if attr_name == "" {
		db.DB.Where("root_url = ?", link).Find(&config2)
		db.DB.Delete(&config)

	} else {
		if err := db.DB.Where("root_url = ?", link).Where("attr_name = ?", attr_name).First(&config).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Record not found!"})
			return
		}
		db.DB.Delete(&config)
	}
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "delete config successful"})
}
