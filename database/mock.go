package database

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/models"
)

var (
	errFundNotFound     = errors.New("fund not found")
	errInvestorNotFound = errors.New("investor not found")
)

type MockDb struct {
	funds       map[uuid.UUID]models.Fund
	investors   map[uuid.UUID]models.Investor
	investments map[uuid.UUID]models.Investment
	investorEmailIndex map[string]uuid.UUID
	mu                 sync.RWMutex
}

func NewMockDb() *MockDb {
	return &MockDb{
		funds:              make(map[uuid.UUID]models.Fund),
		investors:          make(map[uuid.UUID]models.Investor),
		investments:        make(map[uuid.UUID]models.Investment),
		investorEmailIndex: make(map[string]uuid.UUID),
	}
}

func (db *MockDb) CreateFund(createFund models.CreateFund) (models.Fund, error) {
	// Simulate DB constraints by validating DTO
	if err := createFund.Validate(); err != nil {
		return models.Fund{}, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	id := uuid.New()
	now := time.Now().UTC()

	fund := models.Fund{
		ID:            id,
		Name:          createFund.Name,
		VintageYear:   createFund.VintageYear,
		TargetSizeUsd: createFund.TargetSizeUsd,
		Status:        createFund.Status,
		CreatedAt:     now,
	}

	db.funds[id] = fund
	return fund, nil
}

func (db *MockDb) ReadFunds() ([]models.Fund, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	funds := make([]models.Fund, 0, len(db.funds))
	for _, f := range db.funds {
		funds = append(funds, f)
	}
	return funds, nil
}

func (db *MockDb) ReadFundByID(id uuid.UUID) (models.Fund, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	fund, ok := db.funds[id]
	if !ok {
		return models.Fund{}, errFundNotFound
	}
	return fund, nil
}

func (db *MockDb) UpdateFund(fund models.Fund) (models.Fund, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	existing, ok := db.funds[fund.ID]
	if !ok {
		return models.Fund{}, errFundNotFound
	}

	if fund.CreatedAt.IsZero() {
		fund.CreatedAt = existing.CreatedAt
	}

	db.funds[fund.ID] = fund
	return fund, nil
}

func (db *MockDb) CreateInvestor(createInvestor models.CreateInvestor) (models.Investor, error) {
	if err := createInvestor.Validate(); err != nil {
		return models.Investor{}, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, taken := db.investorEmailIndex[createInvestor.Email]; taken {
		return models.Investor{}, errors.New("email already exists")
	}

	id := uuid.New()
	now := time.Now().UTC()

	investor := models.Investor{
		ID:           id,
		Name:         createInvestor.Name,
		InvestorType: createInvestor.InvestorType,
		Email:        createInvestor.Email,
		CreatedAt:    now,
	}

	db.investors[id] = investor
	db.investorEmailIndex[investor.Email] = id

	return investor, nil
}

func (db *MockDb) ReadInvestors() ([]models.Investor, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	investors := make([]models.Investor, 0, len(db.investors))
	for _, inv := range db.investors {
		investors = append(investors, inv)
	}
	return investors, nil
}

func (db *MockDb) CreateInvestment(createInvestment models.CreateInvestment) (models.Investment, error) {
	if err := createInvestment.Validate(); err != nil {
		return models.Investment{}, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.investors[createInvestment.InvestorID]; !ok {
		return models.Investment{}, errInvestorNotFound
	}
	if _, ok := db.funds[createInvestment.FundID]; !ok {
		return models.Investment{}, errFundNotFound
	}

	id := uuid.New()
	now := time.Now().UTC()

	investment := models.Investment{
		ID:             id,
		InvestorID:     createInvestment.InvestorID,
		FundID:         createInvestment.FundID,
		AmountUsd:      createInvestment.AmountUsd,
		InvestmentDate: createInvestment.InvestmentDate,
		CreatedAt:      now,
	}

	db.investments[id] = investment
	return investment, nil
}

func (db *MockDb) ReadInvestments(fundID uuid.UUID) ([]models.Investment, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	investments := make([]models.Investment, 0)
	for _, inv := range db.investments {
		if inv.FundID == fundID {
			investments = append(investments, inv)
		}
	}
	return investments, nil
}
