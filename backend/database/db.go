package database

import (
	"log"

	"github.com/silaswturner/GambleBankWebsite/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "host=localhost user=admin password=admin dbname=gamble_bank port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Automatically migrate all models
	modelsToMigrate := []interface{}{
		&models.User{},
		// Add other models here
	}

	for _, model := range modelsToMigrate {
		err = db.AutoMigrate(model)
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	}

	DB = db
}
