package main

import (
	"time"
)

// file to be saved
const taskFile = "tasks.json"

// represents each status
const (
	statusToDo       = "todo"
	statusInProgress = "in-progress"
	statusDone       = "done"
)

type Task struct {
	ID          int       `json:"ID"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	//variable which starts with uppercase letter will be exported in json
	// and small letter will not be exported
}
