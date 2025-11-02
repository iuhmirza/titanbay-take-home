package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Separate struct for creating funds allows for more flexibility at a later date.
type CreateFund struct {
	Name          string `json:"name"`
	VintageYear   int   `json:"vintage_year"`
	TargetSizeUsd decimal.Decimal   `json:"target_size_usd"`
	Status        string `json:"status"`
}

// consider validating all errors instead of returning on first error
func (fund *CreateFund) Validate() error {
	if fund.Name == "" {
		return errors.New("name is required")
	}
	if fund.VintageYear < 1900 {
		return errors.New("vintage_year must be greater than 1900")
	}
	if fund.VintageYear >= 2100 {
		return errors.New("vintage_year must be less than 2100")
	}
	if fund.TargetSizeUsd.LessThan(decimal.NewFromInt(1_000_000)) {
		return errors.New("target_size_usd must be greater than 1_000_000")
	}

	if fund.Status == "" {
		return errors.New("status is required")
	}

	if fund.Status != "Fundraising" && fund.Status != "Investing" && fund.Status != "Closed" {
		return errors.New("status must be either 'Fundraising', 'Investing', or 'Closed'")
	}

	return nil
}

type Fund struct {
	ID            uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name          string          `json:"name" gorm:"not null"`
	VintageYear   int             `json:"vintage_year" gorm:"not null;check:vintage_year_range,vintage_year >= 1900 AND vintage_year <= 2100"`
	TargetSizeUsd decimal.Decimal `json:"target_size_usd" gorm:"type:numeric(20,2);not null;default:0"`
	Status        string          `json:"status" gorm:"type:text;not null;check:fund_status_chk,status IN ('Fundraising','Investing','Closed')"`
	CreatedAt     time.Time       `json:"created_at"`
	Investments   []Investment    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type Investor struct {
	ID           uuid.UUID    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string       `json:"name" gorm:"not null"`
	InvestorType string       `json:"investor_type" gorm:"type:text;not null;check:investor_type_chk,investor_type IN ('Individual','Institution','Family Office')"`
	Email        string       `json:"email" gorm:"not null;uniqueIndex;size:320"`
	CreatedAt    time.Time    `json:"created_at"`
	Investments  []Investment `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type Investment struct {
	ID             uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	InvestorID     uuid.UUID       `json:"investor_id" gorm:"type:uuid;not null;index"`
	FundID         uuid.UUID       `json:"fund_id" gorm:"type:uuid;not null;index"`
	AmountUsd      decimal.Decimal `json:"amount_usd" gorm:"type:numeric(20,2);not null;check:amount_usd_nonneg,amount_usd >= 0"`
	CreatedAt      time.Time       `json:"created_at"`
	InvestmentDate time.Time       `json:"investment_date" gorm:"type:date;not null"`
	Fund           Fund            `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Investor       Investor        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
