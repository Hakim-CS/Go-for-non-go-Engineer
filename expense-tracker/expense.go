package main

import "time"

// file to be saved
const expenseFile = "expense.json"

const (
	launch    = "Lauch"
	dinner    = "Dinner"
	breakfast = "Breakfast"
)

type Expense struct {
	Id          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
}
