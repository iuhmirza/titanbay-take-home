package main

import (
	"log"
	"net/http"
	"os"

	"github.com/iuhmirza/titanbay-take-home/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("Environment variable HOST not set.")
	}
	h := handlers.Handler{ }
	e := echo.New()
	e.GET("/funds", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{})
	})
	e.POST("/funds", h.CreateFund)
	e.Logger.Fatal(e.Start(host))
}

func connectToDB() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatal("Environment variable DB_HOST not set.")
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal("Environment variable DB_USER not set.")
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		log.Fatal("Environment variable DB_PASS not set.")
	}

}
