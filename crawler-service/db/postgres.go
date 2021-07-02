package db

import (
	"fmt"

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
		ProductURLs []URLProductLink
	}

	URLProductLink struct {
		gorm.Model
		ProductURL string
		URLModel   uint
	}
)

// CreateDBConnection sets up the database to be used
// the config details can be gotten from .env file
func CreateDBConnection(host, port, user, dbname, password, sslmode string) (*gorm.DB, error) {
	const format = "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	uri := fmt.Sprintf(format, user, password, host, port, dbname, sslmode)
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		logrus.WithError(err).Fatalln("failed to initiate the db")
	}
	err = db.AutoMigrate(&URLModel{})
	return db, err
}
