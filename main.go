package main

import (
	"log"
	"os"

	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	log.Println("Starting server..")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Environment variable PORT not set.")
	}
	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	h := handlers.Handler{Db: db}
	e := echo.New()
	e.GET("/funds", h.ReadFunds)
	e.POST("/funds", h.CreateFund)
	e.Logger.Fatal(e.Start(port))
}
