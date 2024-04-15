package main

import "errors"

var (
	ErrInvalidCurrency   = errors.New("invalid currency")
	ErrLoanAlreadyExists = errors.New("loan already exists")
	ErrLoanDoesNotExists = errors.New("loan does not exists")
)
