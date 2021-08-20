package controller

import (
	"cache_service/crawler-service/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//get product details
func GetProductDetails(context *gin.Context) {
	var product_details []db.ProductDetails
	var count int64
	search_query := context.Query("k")
	if search_query == "" {
		fmt.Println("Please provide a search string")
		return
	}
	search := fmt.Sprintf("%%%s%%", search_query)
	//make search case insensitive
	db.DB.Scopes(db.Paginate(context)).Where("product_url ILIKE ?", search).Or("attr_value ILIKE ?", search).Find(&product_details).Count(&count)
	context.JSON(http.StatusOK,
		gin.H{
			"success": true,
			"product": product_details,
			"total":   count,
		})
}

//delete product details
func DeleteProductDetails(context *gin.Context) {
	var product_details db.ProductDetails
	id := context.Param("id")
	if err := db.DB.Where("id = ?", id).First(&product_details).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Record not found!", "success": false})
		return
	}
	db.DB.Delete(&product_details)
	context.JSON(http.StatusOK, gin.H{"success": true, "message": "delete deltails successful"})
}
