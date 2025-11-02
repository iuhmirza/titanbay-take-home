package database

import (
	"github.com/google/uuid"
	"github.com/iuhmirza/titanbay-take-home/models"
	"gorm.io/gorm"
)

type Db interface {
	CreateFund(models.CreateFund) (models.Fund, error)
	ReadFunds() ([]models.Fund, error)
	UpdateFund(models.Fund) (models.Fund, error)
	ReadFundByID(uuid.UUID) (models.Fund, error)
	CreateInvestor(models.CreateInvestor) (models.Investor, error)
	ReadInvestors() ([]models.Investor, error)
	CreateInvestment(models.CreateInvestment) (models.Investment, error)
	ReadInvestments(uuid.UUID) ([]models.Investment, error)
}

type PGDB struct {
	db *gorm.DB
}

func (pgdb *PGDB) CreateFund(createFund models.CreateFund) (models.Fund, error) {
	fund := models.Fund{
		Name:          createFund.Name,
		VintageYear:   createFund.VintageYear,
		TargetSizeUsd: createFund.TargetSizeUsd,
		Status:        createFund.Status,
	}
	return fund, pgdb.db.Create(&fund).Error
}

func (pgdb *PGDB) ReadFunds() ([]models.Fund, error) {
	funds := make([]models.Fund, 0, 32)
	return funds, pgdb.db.Find(&funds).Error
}

func (pgdb *PGDB) ReadFundByID(id uuid.UUID) (models.Fund, error) {
	var fund models.Fund
	return fund, pgdb.db.First(&fund, "id = ?", id).Error
}

func (pgdb *PGDB) UpdateFund(fund models.Fund) (models.Fund, error) {
	f := models.Fund{ID: fund.ID}
	if err := pgdb.db.First(&f).Error; err != nil {
		return models.Fund{}, err
	}

	return fund, pgdb.db.Save(&fund).Error
}

func (pgdb *PGDB) CreateInvestor(createInvestor models.CreateInvestor) (models.Investor, error) {
	investor := models.Investor{
		Name:         createInvestor.Name,
		InvestorType: createInvestor.InvestorType,
		Email:        createInvestor.Email,
	}

	return investor, pgdb.db.Create(&investor).Error
}

func (pgdb *PGDB) ReadInvestors() ([]models.Investor, error) {
	investors := make([]models.Investor, 0, 32)
	return investors, pgdb.db.Find(&investors).Error
}

func (pgdb *PGDB) CreateInvestment(createInvestment models.CreateInvestment) (models.Investment, error) {
	investment := models.Investment{
		InvestorID:     createInvestment.InvestorID,
		AmountUsd:      createInvestment.AmountUsd,
		InvestmentDate: createInvestment.InvestmentDate,
		FundID: createInvestment.FundID,
	}

	return investment, pgdb.db.Create(&investment).Error
}

func (pgdb *PGDB) ReadInvestments(fundID uuid.UUID) ([]models.Investment, error) {
	investments := make([]models.Investment, 0, 32)
	return investments, pgdb.db.Where("fund_id = ?", fundID).Find(&investments).Error
	
}