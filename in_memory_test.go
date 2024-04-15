package main

import (
	"testing"
)

func TestInMemoryLoanRepository(t *testing.T) {
	repo := NewInMemoryLoanRepository()

	loan1 := Loan{
		LoanDetails: LoanDetails{
			ID:       "1",
			Currency: CurrencyEUR,
		},
		DailyInterest: []Interest{},
	}

	loan2 := Loan{
		LoanDetails:   LoanDetails{ID: "2"},
		DailyInterest: []Interest{},
	}

	loan3 := Loan{
		LoanDetails:   LoanDetails{ID: "3"},
		DailyInterest: []Interest{},
	}

	// create
	err := repo.Create(loan1)
	if err != nil {
		t.Errorf("Unexpected error in Create: %v", err)
	}

	// create (duplicate)
	err = repo.Create(loan1)
	if err == nil {
		t.Errorf("Expected an error in Create when creating duplicate entry, but got none")
	}

	// create (another)
	err = repo.Create(loan2)
	if err != nil {
		t.Errorf("Unexpected error in Create: %v", err)
	}

	// read
	readLoan, err := repo.Read(loan1.LoanDetails.ID)
	if err != nil {
		t.Errorf("Unexpected error in Read: %v", err)
	}
	if readLoan.LoanDetails.ID != loan1.LoanDetails.ID {
		t.Errorf("Read got wrong loan. Got %v, want %v", readLoan.LoanDetails.ID, loan1.LoanDetails.ID)
	}

	// list
	loans := repo.List()
	if len(loans) != 2 {
		t.Errorf("List got wrong number of loans. Got %v, want %v", len(loans), 1)
	}

	// update
	loan1.LoanDetails.Currency = CurrencyUSD
	err = repo.Update(loan1)
	if err != nil {
		t.Errorf("Unexpected error in Update: %v", err)
	}

	// read updated loan
	updatedLoan, err := repo.Read(loan1.LoanDetails.ID)
	if err != nil {
		t.Errorf("Unexpected error in Read: %v", err)
	}
	if updatedLoan.LoanDetails.ID != loan1.LoanDetails.ID {
		t.Errorf("Read got wrong loan. Got %v, want %v", updatedLoan.LoanDetails.ID, loan1.LoanDetails.ID)
	}
	if updatedLoan.LoanDetails.Currency != loan1.LoanDetails.Currency {
		t.Errorf("Updated loan details were not saved. Got %v, want %v", updatedLoan.LoanDetails.Currency, loan1.LoanDetails.Currency)
	}

	// update a non-existing loan
	err = repo.Update(loan3)
	if err == nil {
		t.Errorf("Expected an error when updating a non-existing loan but got none")
	}

	// delete
	err = repo.Delete(loan1.LoanDetails.ID)
	if err != nil {
		t.Errorf("Unexpected error in Delete: %v", err)
	}

	// delete (again)
	err = repo.Delete(loan1.LoanDetails.ID)
	if err == nil {
		t.Errorf("Expected an error when deleting a non-existing loan but got none")
	}

	// read deleted loan
	_, err = repo.Read(loan1.LoanDetails.ID)
	if err == nil {
		t.Errorf("Expected an error when reading a deleted loan but got none")
	}
}
