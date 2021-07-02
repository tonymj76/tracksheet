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
		Name        string
		ProductURLs URLProductLinks
	}

	URLProductLink struct {
		gorm.Model
		ProductURL string
		URLModel   uint
	}
)
type URLProductLinks []URLProductLink

// CreateDBConnection sets up the database to be used
// the config details can be gotten from .env file
func CreateDBConnection() (*gorm.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	// uri := "postgres://postgres:postgres@localhost:5432/cuvva_merchant?sslmode=disable"
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		logrus.WithError(err).Fatalln("failed to initiate the db")
	}
	err = db.AutoMigrate(&URLModel{}, &URLProductLink{})
	return db, err
}
