package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvestor_Success(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	body := models.CreateInvestor{
		Name:         "Alice",
		InvestorType: "Individual",
		Email:        "alice@example.com",
	}

	req := httptest.NewRequest(http.MethodPost, "/investors", marshal(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.CreateInvestor(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var inv models.Investor
	_ = json.Unmarshal(rec.Body.Bytes(), &inv)
	assert.Equal(t, "Alice", inv.Name)
	assert.Equal(t, "alice@example.com", inv.Email)
}

func TestCreateInvestor_ValidationError(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	body := models.CreateInvestor{
		Name:         "Bob",
		InvestorType: "Individual",
		Email:        "not-an-email",
	}

	req := httptest.NewRequest(http.MethodPost, "/investors", marshal(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.CreateInvestor(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReadInvestors_EmptyOK(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	req := httptest.NewRequest(http.MethodGet, "/investors", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.ReadInvestors(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var list []models.Investor
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	assert.Len(t, list, 0)
}

func TestCreateInvestment_Success(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	// seed fund + investor
	f, _ := db.CreateFund(models.CreateFund{
		Name:          "Fund X",
		VintageYear:   2022,
		TargetSizeUsd: decimal.NewFromInt(2_000_000),
		Status:        "Fundraising",
	})
	inv, _ := db.CreateInvestor(models.CreateInvestor{
		Name:         "Eve",
		InvestorType: "Institution",
		Email:        "eve@inst.com",
	})

	body := map[string]any{
		"investor_id":     inv.ID,
		"amount_usd":      decimal.NewFromInt(250000),
		"investment_date": "2024-02-01",
	}
	req := httptest.NewRequest(http.MethodPost, "/funds/"+f.ID.String()+"/investments", marshal(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("fund_id")
	ctx.SetParamValues(f.ID.String())

	err := h.CreateInvestment(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var invst models.Investment
	_ = json.Unmarshal(rec.Body.Bytes(), &invst)
	assert.Equal(t, f.ID, invst.FundID)
	assert.Equal(t, inv.ID, invst.InvestorID)
	assert.Equal(t, "2024-02-01", invst.InvestmentDate)
}
