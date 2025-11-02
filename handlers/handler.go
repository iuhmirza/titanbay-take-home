package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
			"error": fmt.Sprintf("Invalid JSON payload: %v", err),
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

func (h Handler) CreateInvestor(ctx echo.Context) error {
	ci := models.CreateInvestor{}
	if err := ctx.Bind(&ci); err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Invalid JSON payload", err)
	}
	if err := ci.Validate(); err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Validation failed", err)
	}

	fund, err := h.Db.CreateInvestor(ci)
	if err != nil {
		return JSONError(ctx, http.StatusInternalServerError, "Failed to write fund to database", err)
	}
	return ctx.JSON(http.StatusCreated, fund)
}

func (h Handler) ReadFundByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Invalid UUID provided in path parameter", err)
	}

	fund, err := h.Db.ReadFundByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return JSONError(ctx, http.StatusNotFound, "Fund not found", err)
		}
		return JSONError(ctx, http.StatusInternalServerError, "Failed to read fund from database", err)
	}

	return ctx.JSON(http.StatusOK, fund)
}

func (h Handler) UpdateFund(ctx echo.Context) error {
	var fund models.Fund
	if err := ctx.Bind(&fund); err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Invalid JSON payload", err)
	}

	updated, err := h.Db.UpdateFund(fund)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return JSONError(ctx, http.StatusNotFound, "Fund not found", err)
		}
		return JSONError(ctx, http.StatusInternalServerError, "Failed to update fund", err)
	}

	return ctx.JSON(http.StatusOK, updated)
}

func (h Handler) ReadInvestors(ctx echo.Context) error {
	investors, err := h.Db.ReadInvestors()
	if err != nil {
		return JSONError(ctx, http.StatusInternalServerError, "Failed to fetch investors", err)
	}
	return ctx.JSON(http.StatusOK, investors)
}

func (h Handler) CreateInvestment(ctx echo.Context) error {
	fundID, err := uuid.Parse(ctx.Param("fund_id"))
	if err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Invalid UUID for fund_id", err)
	}
	var ci models.CreateInvestment
	if err := ctx.Bind(&ci); err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Invalid JSON payload", err)
	}
	ci.FundID = fundID
	if err := ci.Validate(); err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Validation failed", err)
	}

	investment, err := h.Db.CreateInvestment(ci)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return JSONError(ctx, http.StatusNotFound, "Fund not found", err)
		}
		// handle non-existent fund/investor more gracefully
		return JSONError(ctx, http.StatusInternalServerError, "Failed to create investment", err)
	}

	return ctx.JSON(http.StatusCreated, investment)
}

func (h Handler) ReadInvestments(ctx echo.Context) error {
	fundID, err := uuid.Parse(ctx.Param("fund_id"))
	if err != nil {
		return JSONError(ctx, http.StatusBadRequest, "Invalid UUID for fund_id", err)
	}

	investments, err := h.Db.ReadInvestments(fundID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return JSONError(ctx, http.StatusNotFound, "Fund not found", err)
		}
		return JSONError(ctx, http.StatusInternalServerError, "Failed to fetch investments", err)
	}

	return ctx.JSON(http.StatusOK, investments)
}
