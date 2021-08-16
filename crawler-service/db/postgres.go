package db

import (
	"cache_service/crawler-service/config"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// URLModel is the url links to be save
type (
	URLProductLink struct {
		gorm.Model
		ProductURL *string `gorm:"unique;not null"`
	}
	ProductSelector struct {
		gorm.Model
		AttrName  *string `gorm:"not null"`
		AttrValue *string `gorm:"not null"`
		RootURL   *string `gorm:"not null"`
	}
)

var DB *gorm.DB

// CreateDBConnection sets up the database to be used
// the config details can be gotten from .env file
func CreateDBConnection() {
	uri := config.GetEnv("DATABASE_URI")
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		logrus.WithError(err).Fatalln("failed to initiate the db")
	}
	err = db.AutoMigrate(&URLProductLink{}, &ProductSelector{})
	if err != nil {
		logrus.WithError(err).Fatalln("error migrating")
	}
	DB = db
}

// pagination hook
func Paginate(context *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, err := strconv.Atoi(context.DefaultQuery("page", "0"))
		if err != nil {
			logrus.WithError(err).Fatalln("error parsing page to number switching to default value :", err)
			page = 0
		}
		if page == 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(context.DefaultQuery("limit", "1000"))
		if err != nil {
			logrus.WithError(err).Fatalln("error parsing page limit to number switching to default value :", err)
			pageSize = 1000
		}
		switch {
		case pageSize > 1000:
			pageSize = 1000
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
