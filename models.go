package main

import "errors"

// Separate struct for creating funds allows for more flexibility at a later date.
type CreateFund struct {
	Name          string
	VintageYear   uint
	TargetSizeUsd uint
	Status        string
}

// consider validating all errors instead of returning on first error
func (fund *CreateFund) Validate() error {
	if fund.Name == "" {
		return errors.New("name is required")
	}
	if fund.VintageYear < 1900 {
		return errors.New("vintage_year must be greater than 1900")
	}
	if fund.TargetSizeUsd > 1_000_000_00 {
		return errors.New("target_size_usd must be greater than 10,000,000 cents (use cents instead of dollars)")
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
	Id            string
	Name          string
	VintageYear   uint
	TargetSizeUsd uint
	Status        string
	Created_At    string
}
