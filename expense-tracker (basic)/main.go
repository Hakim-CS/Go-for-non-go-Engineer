package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("nothing written ")
		return
	}

	command := os.Args[1] // add,delete...
	args := os.Args[2:]

	var id int
	for _, arg := range args {
		id, _ = strconv.Atoi(arg)
	}

	switch command {
	case "add":
		cmdAdd(args) //arg is a list [dinner 10]
	case "list":
		cmdList()
	case "summary":
		cmdSummary()
	case "delete":
		cmdDelete(id)
	case "summary-m":
		cmdSummaryByMonth(id)
	case "update":
		cmdUpdate(args)

	}

	// what is the data type of data varialble : answer: []Expense, it means a list of Expense struct

}
