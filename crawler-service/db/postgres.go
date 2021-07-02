package db

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// URLModel is the url links to be save
type (
	URLModel struct {
		gorm.Model
		Name string
	}

	URLProductLink struct {
		gorm.Model
		ProductURL string
		URLModelID uint
		URLModel   URLModel
	}
)
type URLProductLinks []URLProductLink

var DB *gorm.DB

// CreateDBConnection sets up the database to be used
// the config details can be gotten from .env file
func CreateDBConnection() {
	uri := os.Getenv("DATABASE_URI")
	// uri := "postgres://postgres:postgres@localhost:5432/cuvva_merchant?sslmode=disable"
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
