package db

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// URLModel is the url links to be save
type (
	URLModel struct {
		gorm.Model
		Name *string `gorm:"unique;not null"`
	}

	URLProductLink struct {
		gorm.Model
		ProductURL *string `gorm:"unique_index;not null"`
		URLModelID uint
		URLModel   *URLModel
	}
)

var DB *gorm.DB

// CreateDBConnection sets up the database to be used
// the config details can be gotten from .env file
func CreateDBConnection() {
	uri := os.Getenv("DATABASE_URI")
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		logrus.WithError(err).Fatalln("failed to initiate the db")
	}
	err = db.AutoMigrate(&URLModel{}, &URLProductLink{})
	if err != nil {
		logrus.WithError(err).Fatalln("error migrating")

	}
	DB = db
}

//pagination hook
func Paginate(r *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.DefaultQuery("page", "0"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.DefaultQuery("limit", "100"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
