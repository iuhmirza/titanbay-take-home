package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateFund_Success(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	req := httptest.NewRequest(http.MethodPost, "/funds", marshal(models.CreateFund{
		Name:          "Fund I",
		VintageYear:   2020,
		TargetSizeUsd: decimal.NewFromInt(5_000_000),
		Status:        "Fundraising",
	}))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := h.CreateFund(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var got models.Fund
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "Fund I", got.Name)
	assert.False(t, got.ID == uuid.Nil)
	assert.WithinDuration(t, time.Now(), got.CreatedAt, 2*time.Second)
}

func TestCreateFund_BadJSON(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	req := httptest.NewRequest(http.MethodPost, "/funds", bytes.NewBufferString(`{"name":`)) // broken json
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := h.CreateFund(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "Invalid JSON payload")
}

func TestCreateFund_ValidationError(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	body := models.CreateFund{
		Name:          "",
		VintageYear:   2020,
		TargetSizeUsd: decimal.NewFromInt(10),
		Status:        "Fundraising",
	}

	req := httptest.NewRequest(http.MethodPost, "/funds", marshal(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.CreateFund(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "Bad fund request")
}

func TestReadFunds_EmptyOK(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	req := httptest.NewRequest(http.MethodGet, "/funds", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.ReadFunds(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var funds []models.Fund
	_ = json.Unmarshal(rec.Body.Bytes(), &funds)
	assert.Len(t, funds, 0)
}

func TestReadFunds_WithData(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	_, _ = db.CreateFund(models.CreateFund{
		Name:          "Fund A",
		VintageYear:   2019,
		TargetSizeUsd: decimal.NewFromInt(1_000_000),
		Status:        "Investing",
	})

	req := httptest.NewRequest(http.MethodGet, "/funds", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.ReadFunds(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var funds []models.Fund
	_ = json.Unmarshal(rec.Body.Bytes(), &funds)
	assert.Len(t, funds, 1)
	assert.Equal(t, "Fund A", funds[0].Name)
}

func TestReadFundByID_Success(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	f, _ := db.CreateFund(models.CreateFund{
		Name:          "FindMe",
		VintageYear:   2021,
		TargetSizeUsd: decimal.NewFromInt(2_000_000),
		Status:        "Fundraising",
	})

	req := httptest.NewRequest(http.MethodGet, "/funds/"+f.ID.String(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(f.ID.String())

	err := h.ReadFundByID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var got models.Fund
	_ = json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Equal(t, f.ID, got.ID)
	assert.Equal(t, "FindMe", got.Name)
}

func TestReadFundByID_InvalidUUID(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	req := httptest.NewRequest(http.MethodGet, "/funds/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("not-a-uuid")

	err := h.ReadFundByID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReadFundByID_NotFoundMapsTo404(t *testing.T) {
	e := echo.New()
	ndb := notFoundDb{delegate: database.NewMockDb()}
	h := Handler{Db: ndb}

	id := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/funds/"+id.String(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())

	err := h.ReadFundByID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateFund_Success(t *testing.T) {
	e := echo.New()
	db := database.NewMockDb()
	h := Handler{Db: db}

	seed, _ := db.CreateFund(models.CreateFund{
		Name:          "Old",
		VintageYear:   2010,
		TargetSizeUsd: decimal.NewFromInt(1_500_000),
		Status:        "Investing",
	})

	update := models.Fund{
		ID:            seed.ID,
		Name:          "New",
		VintageYear:   seed.VintageYear,
		TargetSizeUsd: seed.TargetSizeUsd,
		Status:        seed.Status,
		CreatedAt:     seed.CreatedAt,
	}

	req := httptest.NewRequest(http.MethodPut, "/funds", marshal(update))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.UpdateFund(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var got models.Fund
	_ = json.Unmarshal(rec.Body.Bytes(), &got)
	assert.Equal(t, "New", got.Name)
}

func TestUpdateFund_NotFoundMapsTo404(t *testing.T) {
	e := echo.New()
	h := Handler{Db: notFoundDb{delegate: database.NewMockDb()}}

	update := models.Fund{
		ID:            uuid.New(),
		Name:          "DoesNotExist",
		VintageYear:   2000,
		TargetSizeUsd: decimal.NewFromInt(1_000_000),
		Status:        "Closed",
	}

	req := httptest.NewRequest(http.MethodPut, "/funds", marshal(update))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.UpdateFund(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateFund_BadJSON(t *testing.T) {
	e := echo.New()
	h := Handler{Db: database.NewMockDb()}

	req := httptest.NewRequest(http.MethodPut, "/funds", bytes.NewBufferString(`{"id":`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.UpdateFund(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}