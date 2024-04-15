package main

import "sync"

var _ (LoanRepository) = (*inMemoryLoanRepository)(nil)

// inMemoryLoanRepository is an in-memory implementation of LoanRepository
type inMemoryLoanRepository struct {
	loans map[string]Loan
	mx    sync.RWMutex
}

// NewInMemoryLoanRepository creates a new in-memory LoanRepository
func NewInMemoryLoanRepository() *inMemoryLoanRepository {
	return &inMemoryLoanRepository{
		loans: map[string]Loan{},
		mx:    sync.RWMutex{},
	}
}

// Create implements LoanRepository
func (i *inMemoryLoanRepository) Create(loan Loan) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	if _, ok := i.loans[loan.LoanDetails.ID]; ok {
		return ErrLoanAlreadyExists
	}

	i.loans[loan.LoanDetails.ID] = loan
	return nil
}

// Read implements LoanRepository
func (i *inMemoryLoanRepository) Read(id string) (Loan, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()

	loan, ok := i.loans[id]
	if !ok {
		return Loan{}, ErrLoanDoesNotExists
	}

	return loan, nil
}

// List implements LoanRepository
func (i *inMemoryLoanRepository) List() map[string]Loan {
	i.mx.RLock()
	defer i.mx.RUnlock()

	return i.loans
}

// Update implements LoanRepository
func (i *inMemoryLoanRepository) Update(loan Loan) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	if _, ok := i.loans[loan.LoanDetails.ID]; !ok {
		return ErrLoanDoesNotExists
	}

	i.loans[loan.LoanDetails.ID] = loan
	return nil
}

// Delete implements LoanRepository
func (i *inMemoryLoanRepository) Delete(id string) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	if _, ok := i.loans[id]; !ok {
		return ErrLoanDoesNotExists
	}

	delete(i.loans, id)
	return nil
}
