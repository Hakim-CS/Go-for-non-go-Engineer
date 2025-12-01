package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func cmdAdd(arg []string) {
	if len(arg) < 1 {
		fmt.Println("provide amount !")
		return
	}

	expense, _ := loadFile()

	var description string
	switch arg[0] {
	case "breakfast":
		description = breakfast
	case "launch":
		description = launch
	case "dinner":
		description = dinner
	default:
		fmt.Println("unknown male time")
		return
	}

	now := time.Now
	amountString := arg[1]
	amountInt, _ := strconv.Atoi(amountString)

	newExpense := Expense{
		Id:          getNextId(expense),
		Date:        now(),
		Description: description,
		Amount:      amountInt,
	}

	expense = append(expense, newExpense)
	if err := saveFile(&expense); err != nil {
		fmt.Println("error ocurred while saving : ", err)
	}

	fmt.Println("Expense added successfully ")

}

// list the Expenses
// ID    Description	Amount	Date
func cmdList() {
	data, _ := loadFile()
	fmt.Println("Lis of Expenses: ")

	fmt.Println("ID  	 Desc   Amount		Date")

	for i := 0; i < len(data); i++ {
		fmt.Printf("%d \t%s \t%d \t\t%s\n",
			data[i].Id,
			data[i].Description,
			data[i].Amount,
			data[i].Date.Format("2006-01-02"),
		)
	}

}

// summary
//input: go run . summary
//output: total expense 50

// get the json file
//it return a list ,
//sum all data[i].amout and print

func cmdSummary() {
	data, _ := loadFile()

	total := 0
	for i := 0; i < len(data); i++ {
		total += data[i].Amount
	}

	fmt.Printf("total expense: $%d", total)

}

// input: go run . delete 2
// output: expense 2 deleted successfully

func cmdDelete(index int) {
	data, _ := loadFile()

	for i := 0; i < len(data); i++ {
		if data[i].Id == index {
			data = append(data[:i], data[i+1:]...)
			break
		}
	}
	if err := saveFile(&data); err != nil {
		fmt.Println("Error while saving the file: ", err)
		return
	}
	fmt.Println("file deleted successfully ")
	// how to remove an element from a list in go
	// answer : use append function
	// example : a = append(a[:i], a[i+1:]...)
	// explanation : a[:i] means all elements before index i
	// a[i+1:] means all elements after index i
	// ... means unpack the slice into individual elements
}

// summaryByMonth
func cmdSummaryByMonth(month int) {
	if month > 12 || month < 1 {
		fmt.Println("invalid month ")
		return
	}
	data, _ := loadFile()

	total := 0
	m := 0
	for i := 0; i < len(data); i++ {
		if int(data[i].Date.Month()) == month {
			m = i
			total += data[i].Amount

		}
	}

	fmt.Printf("total expense for %s is %d", data[m].Date.Month(), total)

}

// update a male time by its id
// instance: go run . update 1 breakfast
func cmdUpdate(args []string) {
	//args = [1 , breakfst]
	if len(args) < 2 {
		fmt.Println("invalid input ")
		return
	}
	id, _ := strconv.Atoi(args[0])
	male := args[1:] // breakfst
	// is there any ready function in go to join a list of string into a single string
	// answer : yes , strings.Join
	desc := strings.Join(male, " ")

	expenses, _ := loadFile()
	// find that instance with that id
	for i := range expenses {
		if id == expenses[i].Id {
			expenses[i].Description = desc
			break
		}
	}
	if err := saveFile(&expenses); err != nil {
		fmt.Println("error while saving the file , ", err)
		return
	}
	fmt.Println("Expense description updated successfully ")

}
