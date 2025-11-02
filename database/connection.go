package database

import (
	"log"
	"os"

	"github.com/iuhmirza/titanbay-take-home/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func ConnectToDB() (Db, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Environment variable DB_URL not set.")
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Fund{}, &models.Investor{}, &models.Investment{}); err != nil {
		log.Fatalf("Failed to automatically migrate database: %v", err)
	}
	log.Println("Migrated database successfully")
	return &PGDB{db}, nil
}
