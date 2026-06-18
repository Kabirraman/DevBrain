package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error

	DB, err = gorm.Open(
		postgres.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{},
	)

	return err
}