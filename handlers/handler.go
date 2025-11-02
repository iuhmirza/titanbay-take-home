package handlers

import (
	"fmt"
	"net/http"

	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Db database.Db
}

func JSONError(ctx echo.Context, httpStatus int, errorMessage string, err error) error {
	return ctx.JSON(httpStatus, echo.Map{"error": fmt.Sprintf("%v: %v", errorMessage, err)})
}

func (h Handler) CreateFund(ctx echo.Context) error {
	createFund := models.CreateFund{}
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
	// add row to db funds table using gorm
	fund, err := h.Db.CreateFund(createFund)
	if err != nil {
		// consider not exposing db error
		return JSONError(ctx, http.StatusInternalServerError, "Failed to write fund to database", err)
	}
	return ctx.JSON(http.StatusCreated, fund)
}

func (h Handler) ReadFunds(ctx echo.Context) error {
	funds, err := h.Db.ReadFunds()
	if err != nil {
		return JSONError(ctx, http.StatusInternalServerError, "Failed to read funds from database", err)
	}
	return ctx.JSON(http.StatusOK, funds)
}
