package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ANSI colour escape codes
const (
	Cyan  = "\033[36m"
	Red   = "\033[31m"
	Reset = "\033[0m"
)

// cli encapsulates the command line interface reading and writing
type cli struct {
	reader         *bufio.Reader
	loanRepository LoanRepository
}

// NewCLI creates a new instance of a cli
func NewCLI(loanRepository LoanRepository) *cli {
	return &cli{
		reader:         bufio.NewReader(os.Stdin),
		loanRepository: loanRepository,
	}
}

// DrawMenu draws the menu for the calculator
func (c *cli) DrawMenu() error {
	fmt.Println("Simple Daily Interest Loan Calculator")

	for {
		fmt.Println()
		input, err := c.requestString("action", "create, history, export, list, update, delete or exit", true)
		if err != nil {
			printErr(err)
			continue
		}

		switch input {
		case "create":
			err = c.handleCreate()
		case "history":
			err = c.handleHistory()
		case "export":
			err = c.handleExport()
		case "list":
			err = c.handleList()
		case "update":
			err = c.handleUpdate()
		case "delete":
			err = c.handleDelete()
		case "exit":
			return nil
		default:
			err = ErrInvalidInput
		}

		if err != nil {
			printErr(err)
		}
	}
}

// handleCreate handles creating a new loan
func (c *cli) handleCreate() error {
	id := randomString(8)

	loanDetails, err := c.requestLoanDetails(id)
	if err != nil {
		return err
	}

	loan := Loan{
		LoanDetails:   loanDetails,
		DailyInterest: CalculateDailySimpleInterest(loanDetails),
	}

	if err := c.loanRepository.Create(loan); err != nil {
		return err
	}

	fmt.Printf("\nCreated loan (%s) with following details\n", sprintColoured(loanDetails.ID, Cyan))
	printLoan(loan)

	return nil
}

// handleHistory handles a loan history
func (c *cli) handleHistory() error {
	id, err := c.requestString("Loan ID", "8 character ID", true)
	if err != nil {
		return err
	}
	loan, err := c.loanRepository.Read(id)
	if err != nil {
		return err
	}

	fmt.Printf("\nFetched history for loan (%s)\n", sprintColoured(loan.LoanDetails.ID, Cyan))
	printLoan(loan)

	return nil
}

// handleExport handles exporting a loan
func (c *cli) handleExport() error {
	id, err := c.requestString("Loan ID", "8 character ID", true)
	if err != nil {
		return err
	}
	loan, err := c.loanRepository.Read(id)
	if err != nil {
		return err
	}

	fmt.Printf("\nExported history for loan (%s) to JSON\n", sprintColoured(loan.LoanDetails.ID, Cyan))

	data, err := json.MarshalIndent(loan, "", "    ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", data)

	return nil
}

// handleList handles fetching a list of loans
func (c *cli) handleList() error {
	loans := c.loanRepository.List()

	if len(loans) == 0 {
		fmt.Println("\tThere are no loans to be listed")
		return nil
	}

	for id := range loans {
		fmt.Println("\t", id)
	}

	return nil
}

// handleUpdate handles a loan update
func (c *cli) handleUpdate() error {
	id, err := c.requestString("Loan ID", "8 character ID", true)
	if err != nil {
		return err
	}
	if _, err := c.loanRepository.Read(id); err != nil {
		return err
	}

	loanDetails, err := c.requestLoanDetails(id)
	if err != nil {
		return err
	}

	updatedLoan := Loan{
		LoanDetails:   loanDetails,
		DailyInterest: CalculateDailySimpleInterest(loanDetails),
	}

	if err := c.loanRepository.Update(updatedLoan); err != nil {
		return err
	}

	fmt.Printf("\nUpdated loan (%s) with following details\n", sprintColoured(loanDetails.ID, Cyan))
	printLoan(updatedLoan)

	return nil
}

// handleDelete handles deleting a loan
func (c *cli) handleDelete() error {
	id, err := c.requestString("Loan ID", "8 character ID", true)
	if err != nil {
		return err
	}

	if ok := c.requestConfirmation("Are you sure you want to continue?"); !ok {
		printColouredln("\tDelete was cancelled	", Red)
		return nil
	}

	err = c.loanRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// requestLoanDetails draws the loan details input form, validates the inputs, and outputs a LoanDetails struct
func (c *cli) requestLoanDetails(id string) (LoanDetails, error) {
	fmt.Println("\nInput the values for the loan")

	var (
		startDate        time.Time
		endDate          time.Time
		loanAmount       float64
		loanCurrency     Currency
		baseInterestRate float64
		margin           float64
		err              error
	)

	for {
		startDate, err = c.requestDate("Start Date", true)
		if err == nil {
			break
		}
		printErr(err)
	}

	for {
		endDate, err = c.requestDateAfter("End Date", startDate, true)
		if err == nil {
			break
		}
		printErr(err)
	}

	for {
		loanAmount, err = c.requestPositiveFloat64("Loan Amount", "principal amount being loaned", true)
		if err == nil {
			break
		}
		printErr(err)
	}

	for {
		loanCurrency, err = c.requestCurrency("Loan Currency", AllowedCurrencies, true)
		if err == nil {
			break
		}
		printErr(err)
	}

	for {
		baseInterestRate, err = c.requestPositiveFloat64("Base Interest Rate", "percentage", true)
		if err == nil {
			break
		}
		printErr(err)
	}

	for {
		margin, err = c.requestPositiveFloat64("Margin", "percentage", true)
		if err == nil {
			break
		}
		printErr(err)
	}

	return LoanDetails{
		ID:               id,
		StartDate:        startDate,
		EndDate:          endDate,
		Currency:         loanCurrency,
		PrincipalAmount:  loanAmount,
		BaseInterestRate: baseInterestRate,
		Margin:           margin,
	}, nil
}

// requestString requests a string input from the user
func (c *cli) requestString(name, hint string, required bool) (string, error) {
	if len(hint) > 0 {
		fmt.Printf("%s%s (%s): %s", Cyan, name, hint, Reset)
	} else {
		fmt.Printf("%s%s: %s", Cyan, name, Reset)
	}

	input, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSuffix(input, "\n")

	if required && len(input) == 0 {
		return "", ErrInvalidInput
	}

	return input, nil
}

// requestFloat64 requests a float input from the user
func (c *cli) requestFloat64(name, hint string, required bool) (float64, error) {
	val, err := c.requestString(name, hint, required)
	if err != nil {
		return 0, err
	}

	if err := validateDecimalPlaces(val, 2); err != nil {
		return 0, err
	}

	floatVal, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, err
	}

	floatVal = math.Round(floatVal*100) / 100

	return floatVal, nil
}

// requestPositiveFloat64 requests a float input from the user that is >= 0
func (c *cli) requestPositiveFloat64(name, hint string, required bool) (float64, error) {
	val, err := c.requestFloat64(name, hint, required)
	if err != nil {
		return 0, err
	}

	if val < 0 {
		return 0, errors.Wrap(ErrInvalidInput, "value must be greater than 0")
	}

	return val, nil
}

// requestDate requests a date input from the user in the format YYYY-MM-DD
func (c *cli) requestDate(name string, required bool) (time.Time, error) {
	val, err := c.requestString(name, "YYYY-MM-DD", required)
	if err != nil {
		return time.Time{}, err
	}

	date, err := time.Parse("2006-01-02", val)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// requestDateAfter requests a date input from the user in the format YYYY-MM-DD after a specific date
func (c *cli) requestDateAfter(name string, date time.Time, required bool) (time.Time, error) {
	input, err := c.requestDate(name, required)
	if err != nil {
		return time.Time{}, err
	}

	if !input.After(date) {
		return time.Time{}, errors.Wrap(ErrInvalidInput, "end date needs to be after start date")
	}

	return input, nil
}

// requestCurrency requests a currency input from the user in the ISO 4217 format
func (c *cli) requestCurrency(name string, allowedCurrencies []Currency, required bool) (Currency, error) {
	currencies := ""
	for _, c := range allowedCurrencies {
		currencies += c.String() + ", "
	}
	currencies = currencies[:len(currencies)-2]

	input, err := c.requestString(name, currencies, required)
	if err != nil {
		return "", err
	}
	input = strings.ToUpper(input)

	currency := Currency(input)
	if err := currency.Validate(); err != nil {
		return "", err
	}

	return currency, nil
}

// requestConfirmation requests a yes/y/no/n confirmation
func (c *cli) requestConfirmation(msg string) bool {
	for {
		val, err := c.requestString(msg, "", true)
		if err != nil {
			return false
		}

		val = strings.ToLower(val)

		switch val {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			continue
		}
	}
}

// validateDecimalPlaces validates whether an input string's decimal places is within the allowed amount
func validateDecimalPlaces(input string, places int) error {
	decimalIndex := strings.Index(input, ".")
	if decimalIndex == -1 {
		return nil
	}
	decimalPlaces := len(input) - decimalIndex - 1
	if decimalPlaces > places {
		return errors.Wrapf(ErrInvalidDecimalPlaces, "max decimal places should not exceed %d", places)
	}
	return nil
}

// printErr prints an error with ANSI escape codes to colourise the output in red
func printErr(err error) {
	errStr := err.Error()
	errStr = strings.ToUpper(string(errStr[0])) + errStr[1:]
	fmt.Printf("\t%s\n", sprintColoured(errStr, Red))
}

// sprintColoured returns a string wrapped in colourised ANSI escape codes
func sprintColoured(value, colour string) string {
	return fmt.Sprintf("%s%s%s", colour, value, Reset)
}

// printColouredln returns prints a colourised string to the console with a new line
func printColouredln(value, colour string) {
	fmt.Println(sprintColoured(value, colour))
}

// printValf prints a value with a colourised name
func printValf(prefix, name, fmtStr string, args ...any) {
	fmt.Printf("%s%s: %s", prefix, sprintColoured(name, Cyan), fmt.Sprintf(fmtStr, args...))
}

// printLoan prints out the loan details in a stylised way
func printLoan(loan Loan) {
	printValf("", "Loan ID", "%s\n", loan.LoanDetails.ID)
	printValf("", "Start Date", "%s\n", loan.LoanDetails.StartDate.Format("2006-01-02"))
	printValf("", "End Date", "%s\n", loan.LoanDetails.EndDate.Format("2006-01-02"))
	printValf("", "Loan Amount", " %s%.2f\n", loan.LoanDetails.Currency.Symbol(), loan.LoanDetails.PrincipalAmount)
	printValf("", "Loan Currency", "%s\n", loan.LoanDetails.Currency)
	printValf("", "Base Interest Rate", " %v%%\n", loan.LoanDetails.BaseInterestRate)
	printValf("", "Margin", "%v%%\n", loan.LoanDetails.Margin)

	for _, interest := range loan.DailyInterest {
		printValf("\t- ", "Accrual Date", "%s\n", interest.AccrualDate.Format("2006-01-02"))
		printValf("\t  ", "Days Elapsed", "%d\n", interest.DaysElapsed)
		printValf("\t  ", "Daily Interest Amount without Margin", " %s%f\n", loan.LoanDetails.Currency.Symbol(), interest.DailyInterestWithoutMargin)
		printValf("\t  ", "Daily Interest Amount Accrued", " %s%f\n", loan.LoanDetails.Currency.Symbol(), interest.DailyInterestAccrued)
		printValf("\t  ", "Total Interest", " %s%f\n", loan.LoanDetails.Currency.Symbol(), interest.TotalInterest)
	}

	printValf("\n", "Loan ID", "%s\n", loan.LoanDetails.ID)
}
