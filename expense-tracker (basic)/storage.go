package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadFile() ([]Expense, error) {
	data, _ := os.ReadFile(expenseFile)

	// if err != nil {
	// 	if errors.Is(err,os.ErrNotExist){
	// 		return []Expense{},err
	// 	}
	// 	fmt.Println("error while loading the file: ", err)
	// 	return  nil, err
	// }

	// // if it's empty
	// if len(data) == 0{
	// 	return []Expense{}, nil
	// }
	var expenses []Expense
	err := json.Unmarshal(data, &expenses)
	if err != nil {
		fmt.Println("error while Unmarshaling the file: ", err)
		return nil, err
	}
	return expenses, nil

}

// save file
func saveFile(expense *[]Expense) error {
	data, _ := json.MarshalIndent(expense, "", "  ")
	return os.WriteFile(expenseFile, data, 0644)

}

func getNextId(tasks []Expense) int {
	maxId := 0
	for _, t := range tasks {
		if t.Id > maxId {
			maxId = t.Id
		}
	}
	return maxId + 1
}
