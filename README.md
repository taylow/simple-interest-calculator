# Simple Loan Calculator

This is a simple loan calculator made as part of an engineering technical task for Oneiro Solutions Ltd.

The task was to make a console application that is capable of calculating daily simple interest given some input.

## 🧐 Simple Interest

- Simple interest is calculated by multiplying loan principal by the interest rate and then by the term of a loan.
  - `Daily Interest Rate` x `Principal Amount` x `Number of Days Elapsed`
  - Simple Interest = P × r × n
    - P = Principal
    - r = Interest rate
    - n = Term of loan
- Simple interest involves no calculation of compound interest.
  - Neither compounding interest nor calculation of the interest rate against a growing total balance is involved.
- Daily Simple Interest
  - Simple Interest is similar to Daily Simple Interest except that with the latter, interest accrues daily and is added to your account balance.
  - Also, while loan balances on simple interest debt are reduced on the payment due date, daily simple interest loan balances are reduced on the day payments are received.

## 🚀 Usage

The Makefile provided with this repo contains all you need to get started. Simply run `make help` for a list of available commands.

If you have Go installed, simply run `make run` to run the calculator.

If you have Docker installed, simply run `make docker` to run a containerised copy of the calculator.

Once running, the command line tool will guide you through the available routes.

From the root, you can choose:

- `create` - start a new loan
- `history` - see the history of an existing loan
- `export` - export the history of an existing loan as JSON
- `list` - list existing loan IDs
- `update` - update existing loan details
- `delete` - delete an existing loan

Each of the commands will enter into a sub menu, where a series of inputs will be requested. All inputs are sanitised and validated.

## 🧪 Testing & Vetting

To test the tool, simple run `make test`.

The CLI/navigational portion of this tool grew slightly out of scope, and is not currently tested, but all core parts of the tool are.

`make vet` runs the code through the `go vet` command to scan the codebase for suspicious constructs.

## 🤔 Some Uncertainties

- `Daily Interest Amount without Margin` and `Daily Interest Amount Accrued` are repeated per entry. This would in theory change if payments were made between the start and end, but as it stands this value never changes.
- Unsure whether day 1 should start on the same day the loan starts, or the day after, same for calculating total days.
- Floating point arithmetic and rounding was a bit of a pain to test, as always.
