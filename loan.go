package main

import (
	"slices"
	"time"
)

const (
	CurrencyEUR = "EUR"
	CurrencyGBP = "GBP"
	CurrencyUSD = "USD"
)

var (
	AllowedCurrencies = []Currency{
		CurrencyEUR,
		CurrencyGBP,
		CurrencyUSD,
	}
)

// Currency holds a 3 letter ISO 4217 currency code
type Currency string

// String stringifies the currency
func (c Currency) String() string {
	return string(c)
}

// Symbol returns the corresponding symbol to the currency
func (c Currency) Symbol() string {
	switch c {
	case CurrencyEUR:
		return "€"
	case CurrencyGBP:
		return "£"
	case CurrencyUSD:
		return "$"
	default:
		return ""
	}
}

// Validate validates whether the ISO 4217 currency is supported
func (c Currency) Validate() error {
	if ok := slices.Contains(AllowedCurrencies, c); !ok {
		return ErrInvalidCurrency
	}

	return nil
}

// Loan represents a loan and the accompanying daily accrued interest
type Loan struct {
	LoanDetails   LoanDetails `json:"loan_details"`   // LoanDetails contains all details of the loan
	DailyInterest []Interest  `json:"daily_interest"` // DailyInterest contains interest data for each day of the loan period
}

// LoanDetails holds details of a loan
type LoanDetails struct {
	ID               string    `json:"id"`                 // ID is the unique identifier for the loan
	StartDate        time.Time `json:"start_date"`         // StartDate is the the start of the loan period
	EndDate          time.Time `json:"end_date"`           // EndDate is the end of the loan period
	Currency         Currency  `json:"currency"`           // Currency is an ISO 4217 3-letter currency code
	PrincipalAmount  float64   `json:"principal_amount"`   // PrincipalAmount is the initial loan amount
	BaseInterestRate float64   `json:"base_interest_rate"` // BaseInterestRate represents a percentage for the base interest rate
	Margin           float64   `json:"margin"`             // Margin is the additional interest on top of the base interest rate
}

// Interest holds information about daily accrued interest from the loan
type Interest struct {
	AccrualDate                time.Time `json:"accrual_date"`                  // AccrualDate is the date the interest was accrued
	DaysElapsed                int       `json:"days_elapsed"`                  // DaysElapsed is the number of days elapsed since the start date of the loan
	DailyInterestWithoutMargin float64   `json:"daily_interest_without_margin"` // DailyInterestWithoutMargin is the daily interest accrued without the margin
	DailyInterestAccrued       float64   `json:"daily_interest_accrued"`        // DailyInterestAccrued is the total daily interest accrued
	TotalInterest              float64   `json:"total_interest"`                // TotalInterest is the total accrued interest calculated over the given period
}

// LoanRepository is an abstraction on the storage of loans
type LoanRepository interface {
	// Create creates a new loan
	Create(loan Loan) error
	// Loan reads a loan from the store
	Read(id string) (Loan, error)
	// List lists all available loans
	List() map[string]Loan
	// Update updates the details of an existing loan
	Update(loan Loan) error
	// Delete deletes an existing loan
	Delete(id string) error
}

// CalculateDailySimpleInterest calculates the daily accrued interest using the daily simple interest formula
func CalculateDailySimpleInterest(loan LoanDetails) []Interest {
	totalDays := int(loan.EndDate.Sub(loan.StartDate).Hours() / 24)
	dailyInterestRateWithoutMargin := dailyInterestRate(loan.BaseInterestRate)
	dailyInterestRateWithMargin := dailyInterestRate(loan.BaseInterestRate + loan.Margin)
	dailyInterest := make([]Interest, totalDays)
	totalInterest := 0.0

	for i := 0; i < totalDays; i++ {
		dailyInterestWithoutMargin := dailyInterestRateWithoutMargin * loan.PrincipalAmount
		dailyInterestWithMargin := dailyInterestRateWithMargin * loan.PrincipalAmount
		totalInterest += dailyInterestWithMargin

		interest := Interest{
			AccrualDate:                loan.StartDate.Add(time.Duration(i) * 24 * time.Hour),
			DaysElapsed:                i + 1,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              totalInterest,
		}
		dailyInterest[i] = interest
	}

	return dailyInterest
}

// dailyInterestRate divides the annual interest rate into a daily amount
func dailyInterestRate(baseRate float64) float64 {
	return baseRate / 100 / 365
}
