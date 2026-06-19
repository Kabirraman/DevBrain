package database

import (
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/Kabirraman/DevBrain/internal/models"

)

var DB *gorm.DB

func Connect() error {

	var err error

	DB, err = gorm.Open(
		postgres.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{},
	)

	if err != nil {
		return err
	}

	err = DB.AutoMigrate(
	&models.User{},
	&models.Resource{},
)

	return err
}