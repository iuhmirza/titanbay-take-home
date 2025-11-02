package database

import (
	"github.com/iuhmirza/titanbay-take-home/models"
	"gorm.io/gorm"
)

type Db interface {
	CreateFund(models.CreateFund) (models.Fund, error)
	ReadFunds() ([]models.Fund, error)
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