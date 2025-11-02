package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func JSONError(ctx echo.Context, httpStatus int, errorMessage string, err error) error {
	return ctx.JSON(httpStatus, echo.Map{"error": fmt.Sprintf("%v: %v", errorMessage, err)})
}

func CreateFundHandler(ctx echo.Context) error {
	createFund := CreateFund{}
	if err := ctx.Bind(&createFund); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("Failed to bind JSON to struct: %v", err),
		})
	}
	if err := createFund.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("Bad fund request: %v", err),
		})
	}
	fund := Fund{
		Name:          createFund.Name,
		VintageYear:   createFund.VintageYear,
		TargetSizeUsd: createFund.TargetSizeUsd,
		Status:        createFund.Status,
	}
	// add row to db funds table
	return ctx.JSON(http.StatusCreated, fund)
}

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("Environment variable HOST not set.")
	}
	e := echo.New()
	e.GET("/funds", func(ctx echo.Context) error {

		return ctx.JSON(http.StatusOK, echo.Map{})
	})
	e.Logger.Fatal(e.Start(host))
}
