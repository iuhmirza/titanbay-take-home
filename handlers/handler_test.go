package handlers

import (
	"bytes"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/database"
	"github.com/iuhmirza/titanbay-take-home/models"
	"gorm.io/gorm"
)

type notFoundDb struct {
	delegate database.Db
}

func (n notFoundDb) CreateFund(cf models.CreateFund) (models.Fund, error) {
	return n.delegate.CreateFund(cf)
}

func (n notFoundDb) ReadFunds() ([]models.Fund, error) { return n.delegate.ReadFunds() }

func (n notFoundDb) UpdateFund(models.Fund) (models.Fund, error) {
	return models.Fund{}, gorm.ErrRecordNotFound
}

func (n notFoundDb) ReadFundByID(uuid.UUID) (models.Fund, error) {
	return models.Fund{}, gorm.ErrRecordNotFound
}

func (n notFoundDb) CreateInvestor(ci models.CreateInvestor) (models.Investor, error) {
	return n.delegate.CreateInvestor(ci)
}

func (n notFoundDb) ReadInvestors() ([]models.Investor, error) { return n.delegate.ReadInvestors() }

func (n notFoundDb) CreateInvestment(ci models.CreateInvestment) (models.Investment, error) {
	return n.delegate.CreateInvestment(ci)
}

func (n notFoundDb) ReadInvestments(uuid.UUID) ([]models.Investment, error) {
	return n.delegate.ReadInvestments(uuid.UUID{})
}

func marshal(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

// func TestFail(t *testing.T) {
// 	t.FailNow()
// }
