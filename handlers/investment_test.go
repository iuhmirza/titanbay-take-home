package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvestment_InvalidFundUUID(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	body := map[string]any{
		"investor_id":     uuid.New(),
		"amount_usd":      decimal.NewFromInt(100),
		"investment_date": "2023-01-01",
	}
	req := httptest.NewRequest(http.MethodPost, "/funds/not-a-uuid/investments", marshal(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("fund_id")
	ctx.SetParamValues("not-a-uuid")

	err := h.CreateInvestment(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateInvestment_ValidationError(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	// seed required fund + investor to avoid FK error—let validation fail on amount
	f, _ := db.CreateFund(models.CreateFund{
		Name:          "Fund Y",
		VintageYear:   2020,
		TargetSizeUsd: decimal.NewFromInt(1_000_000),
		Status:        "Investing",
	})
	inv, _ := db.CreateInvestor(models.CreateInvestor{
		Name:         "Carl",
		InvestorType: "Family Office",
		Email:        "carl@family.com",
	})

	body := map[string]any{
		"investor_id":     inv.ID,
		"amount_usd":      decimal.NewFromInt(0), // invalid per validation
		"investment_date": "2024-06-01",
	}
	req := httptest.NewRequest(http.MethodPost, "/funds/"+f.ID.String()+"/investments", marshal(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("fund_id")
	ctx.SetParamValues(f.ID.String())

	err := h.CreateInvestment(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReadInvestments_Success(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	// seed fund + investor + two investments (only one for our fund)
	f1, _ := db.CreateFund(models.CreateFund{
		Name:          "Fund One",
		VintageYear:   2018,
		TargetSizeUsd: decimal.NewFromInt(1_000_000),
		Status:        "Investing",
	})
	f2, _ := db.CreateFund(models.CreateFund{
		Name:          "Fund Two",
		VintageYear:   2019,
		TargetSizeUsd: decimal.NewFromInt(1_000_000),
		Status:        "Investing",
	})
	inv, _ := db.CreateInvestor(models.CreateInvestor{
		Name:         "Dana",
		InvestorType: "Individual",
		Email:        "dana@example.com",
	})
	_, _ = db.CreateInvestment(models.CreateInvestment{
		InvestorID:     inv.ID,
		FundID:         f1.ID,
		AmountUsd:      decimal.NewFromInt(1234),
		InvestmentDate: "2024-05-05",
	})
	_, _ = db.CreateInvestment(models.CreateInvestment{
		InvestorID:     inv.ID,
		FundID:         f2.ID,
		AmountUsd:      decimal.NewFromInt(5678),
		InvestmentDate: "2024-06-06",
	})

	req := httptest.NewRequest(http.MethodGet, "/funds/"+f1.ID.String()+"/investments", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("fund_id")
	ctx.SetParamValues(f1.ID.String())

	err := h.ReadInvestments(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var list []models.Investment
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	assert.Len(t, list, 1)
	assert.Equal(t, f1.ID, list[0].FundID)
	assert.Equal(t, "2024-05-05", list[0].InvestmentDate)
}

func TestReadInvestments_InvalidFundUUID(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	req := httptest.NewRequest(http.MethodGet, "/funds/not-a-uuid/investments", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("fund_id")
	ctx.SetParamValues("not-a-uuid")

	err := h.ReadInvestments(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
