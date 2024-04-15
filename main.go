package main

// main is the entrypoint to the program
func main() {
	loanRepository := NewInMemoryLoanRepository()

	cli := NewCLI(loanRepository)
	if err := cli.DrawMenu(); err != nil {
		panic(err)
	}
}
