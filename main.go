package main

import (
	"log"
	"net/http"
	"os"

	"github.com/iuhmirza/titanbay-take-home/handlers"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Environment variable PORT not set.")
	}
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	h := handlers.Handler{Db: db}
	e := echo.New()
	e.GET("/funds", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{})
	})
	e.POST("/funds", h.CreateFund)
	e.Logger.Fatal(e.Start(port))
}

type PGDB struct {
	db *gorm.DB
}

func (pgdb *PGDB) AddFund(createFund models.CreateFund) (models.Fund, error) {
	return models.Fund{}, nil
}

func connectToDB() (handlers.Database, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Environment variable DB_HOST not set.")
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PGDB{db}, nil
}
