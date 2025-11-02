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
	e.PUT("/funds", h.UpdateFund)
	e.GET("/funds/:fund_id", h.ReadFundByID)
	e.GET("/investors", h.ReadInvestors)
	e.POST("/investors", h.CreateInvestor)
	e.GET("/funds/:fund_id/investments", h.ReadInvestments)
	e.POST("/funds/:fund_id/investments", h.CreateInvestment)
	e.Logger.Fatal(e.Start(port))
}
