package main

import (
	"encoding/json"
	"errors"
	"os"
)

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(taskFile)
	// data variable data type is []byte , what that means is that it holds raw bytes read from the file
	//it is like a slice of bytes for example : []byte{0x7b, 0x22, 0x49, 0x44, 0x22, ...} which represents the json content
	//[]byte{0x7b, 0x22, 0x49, 0x44, 0x22} what is that ?
	// that is the json representation of {"ID" but how ? answer
	//
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Task{}, nil
	}
	var tasks []Task
	//json.Unmarshal(data, &tasks) is used to parse the JSON-encoded data and store the result in the variable tasks

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err

	}
	return os.WriteFile(taskFile, data, 0644)
}

// get the next Id
func getNextId(tasks []Task) int {
	maxId := 0
	for _, t := range tasks {
		if t.ID > maxId {
			maxId = t.ID
		}
	}
	return maxId + 1
}

func getbyID(tasks []Task, id int) (*Task, int) {

	for i := range tasks {
		if id == tasks[i].ID {
			return &tasks[i], i
		}
	}
	return nil, -1
}
