package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/models"
	"github.com/labstack/echo/v4"
)

type MockDb struct {
	db map[uint]models.Fund
	mu sync.Mutex
}

func (mockDb *MockDb) CreateFund(createFund models.CreateFund) (models.Fund, error) {
	fund := models.Fund{
		Name:          createFund.Name,
		VintageYear:   createFund.VintageYear,
		TargetSizeUsd: createFund.TargetSizeUsd,
		Status:        createFund.Status,
	}
	fund.ID = uuid.New()
	fund.CreatedAt = time.Now()
	mockDb.mu.Lock()
	mockDb.db[uint(len(mockDb.db))] = fund
	mockDb.mu.Unlock()
	return fund, nil
}

func (mockDb *MockDb) ReadFunds() ([]models.Fund, error) {
	funds := make([]models.Fund, 0, len(mockDb.db))
	mockDb.mu.Lock()
	for _, v := range mockDb.db {
		funds = append(funds, v)
	}
	mockDb.mu.Unlock()
	return funds, nil
}

var fundJSON = `
{
  "name": "Titanbay Growth Fund II",
  "vintage_year": 2025,
  "target_size_usd": 50000000000,
  "status": "Fundraising"
}
`

func TestCreateFund(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/funds", strings.NewReader(fundJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	h := &Handler{Db: &MockDb{db: make(map[uint]models.Fund)}}

	if err := h.CreateFund(ctx); err != nil {
		t.Fatalf("CreateFund returned error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %v, want %v; body = %v ", rec.Code, http.StatusCreated, rec.Body.String())
	}
}
