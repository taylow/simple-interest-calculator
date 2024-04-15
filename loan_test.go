package main

import (
	"math"
	"testing"
	"time"
)

func TestCurrencySymbol(t *testing.T) {
	currencies := map[Currency]string{
		CurrencyEUR: "€",
		CurrencyGBP: "£",
		CurrencyUSD: "$",
	}

	for currency, symbol := range currencies {
		if currency.Symbol() != symbol {
			t.Errorf("Incorrect currency->symbol mapping for %s. Got %s, want %s", currency, currency.Symbol(), symbol)
		}

		if err := currency.Validate(); err != nil {
			t.Errorf("Unexpected error while validating currency %s: %v", currency, err)
		}
	}

	invalidCurrency := Currency("someRandomCurrency")
	err := invalidCurrency.Validate()
	if err == nil {
		t.Errorf("Expected error while validating an invalid currency but got none")
	}

	if symbol := invalidCurrency.Symbol(); symbol != "" {
		t.Errorf("Expected empty symbol from an invalid currency but got %s", symbol)
	}
}

func TestCalculateDailySimpleInterest(t *testing.T) {
	const tolerance = 1e-9

	startDate, err := time.Parse("2006-01-02", "2024-01-01")
	if err != nil {
		t.Errorf("Unexpected error parsing time: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", "2024-01-11")
	if err != nil {
		t.Errorf("Unexpected error parsing time: %v", err)
	}

	loan := LoanDetails{
		StartDate:        startDate,
		EndDate:          endDate,
		Currency:         CurrencyEUR,
		PrincipalAmount:  1000,
		BaseInterestRate: 10,
		Margin:           1,
	}

	dailyInterestWithoutMargin := loan.BaseInterestRate / 100 / 365 * loan.PrincipalAmount
	dailyInterestWithMargin := (loan.BaseInterestRate + loan.Margin) / 100 / 365 * loan.PrincipalAmount

	expectedDailyInterest := []Interest{
		{
			AccrualDate:                startDate,
			DaysElapsed:                1,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin,
		},
		{
			AccrualDate:                startDate.Add(1 * 24 * time.Hour),
			DaysElapsed:                2,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 2,
		},
		{
			AccrualDate:                startDate.Add(2 * 24 * time.Hour),
			DaysElapsed:                3,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 3,
		},
		{
			AccrualDate:                startDate.Add(3 * 24 * time.Hour),
			DaysElapsed:                4,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 4,
		},
		{
			AccrualDate:                startDate.Add(4 * 24 * time.Hour),
			DaysElapsed:                5,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 5,
		},
		{
			AccrualDate:                startDate.Add(5 * 24 * time.Hour),
			DaysElapsed:                6,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 6,
		},
		{
			AccrualDate:                startDate.Add(6 * 24 * time.Hour),
			DaysElapsed:                7,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 7,
		},
		{
			AccrualDate:                startDate.Add(7 * 24 * time.Hour),
			DaysElapsed:                8,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 8,
		},
		{
			AccrualDate:                startDate.Add(8 * 24 * time.Hour),
			DaysElapsed:                9,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 9,
		},
		{
			AccrualDate:                startDate.Add(9 * 24 * time.Hour),
			DaysElapsed:                10,
			DailyInterestWithoutMargin: dailyInterestWithoutMargin,
			DailyInterestAccrued:       dailyInterestWithMargin,
			TotalInterest:              dailyInterestWithMargin * 10,
		},
	}

	dailyInterest := CalculateDailySimpleInterest(loan)

	if len(dailyInterest) != len(expectedDailyInterest) {
		t.Errorf("Daily interest returned more entries than expected. got %d, want %d", len(dailyInterest), len(expectedDailyInterest))
	}

	for i, interest := range dailyInterest {
		expected := expectedDailyInterest[i]

		if interest.AccrualDate != expected.AccrualDate {
			t.Errorf("Unexpected `accrual date` in daily interest. got %v, expected %v", interest.AccrualDate, expected.AccrualDate)
		}
		if interest.DaysElapsed != expected.DaysElapsed {
			t.Errorf("Unexpected `days elapsed` in daily interest. got %v, expected %v", interest.DaysElapsed, expected.DaysElapsed)
		}
		if math.Abs(interest.DailyInterestWithoutMargin-expected.DailyInterestWithoutMargin) > tolerance {
			t.Errorf("Unexpected `daily interest without margin` in daily interest out of tolerance. got %v, expected %v", interest.DailyInterestWithoutMargin, expected.DailyInterestWithoutMargin)
		}
		if math.Abs(interest.DailyInterestAccrued-expected.DailyInterestAccrued) > tolerance {
			t.Errorf("Unexpected `daily interest accrued` in daily interest out of tolerance. got %v, expected %v", interest.DailyInterestAccrued, expected.DailyInterestAccrued)
		}
		if math.Abs(interest.TotalInterest-expected.TotalInterest) > tolerance {
			t.Errorf("Unexpected `total interest` in daily interest out of tolerance. got %v, expected %v", interest.TotalInterest, expected.TotalInterest)
		}
	}
}
