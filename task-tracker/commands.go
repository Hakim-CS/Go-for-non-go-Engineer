package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// why we don't use scanln here?
// because scanln stops reading at the first space
// so if the description has spaces, it will not be read completely

// what does args represent here?
// args is a slice of strings which contains the command-line arguments passed to the "add" command
// For example, if the user runs the command "task-tracker add Buy groceries",
// then args will be []string{"Buy", "groceries"}
// So we need to join these arguments to form the complete description of the task
// we can use strings.Join() function to join the arguments with space as separator
// strings.Join(args, " ") will give us "Buy groceries"

// cmdAdd handles the "add" command to add a new task
// args: command-line arguments passed to the "add" command
// It processes the arguments, creates a new task, and saves it to the task list
func cmdAdd(args []string) {
	//the value of args[0] is "add"
	//args is a slice of strings which contains the command-line arguments passed to the "add" command
	//args be like : ["add","description","status"]
	if len(args) < 1 {
		fmt.Println("add description ")
	}
	description := strings.Join(args, " ") // Join all args to form the description
	//strings.join() joins the elements of a slice into a single string with a specified separator

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("en error occured : ", err)
	}

	now := time.Now
	newTask := Task{
		ID:          getNextId(tasks),
		Description: description,
		Status:      statusToDo, // default status is "todo",declared in task.go
		CreatedAt:   now(),
		UpdatedAt:   now(),
	}
	tasks = append(tasks, newTask)

	err = saveTasks(tasks)
	if err != nil {
		fmt.Println("En error occured while saving the task: ", err)
	}

	fmt.Printf("Task added successfully with ID %d\n", newTask.ID)
}

func cmdDeleteByID(args []string) {
	//so the args[0] is "3", which is the id to be deleted

	if len(args) < 1 {
		fmt.Println("Please provide the task ID to delete.")
	}
	for _, arg := range args {
		// Convert the argument to an integer ID
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Invalid task ID: %s\n", arg)
			continue
		}
		tasks, err := loadTasks()
		if err != nil {
			fmt.Println("Error loading tasks: ", err)
			return
		}
		_, index := getbyID(tasks, id)
		if index == -1 {
			fmt.Printf("Task with ID %d not found.\n", id)
			continue
		}
		// Remove the task from the slice
		//what really this line do? answ
		tasks = append(tasks[:index], tasks[index+1:]...)
		if err := saveTasks(tasks); err != nil {
			fmt.Println("Error saving tasks: ", err)
			return
		}
		fmt.Printf("Task with ID %d deleted successfully.\n", id)

	}
}

func cmdList() {

	data, err := loadTasks()
	if err != nil {
		fmt.Println("En error occured while loading the task: ", err)
	}

	for i := 0; i < len(data); i++ {
		fmt.Printf("task %d detail : ", i+1)
		fmt.Println()
		fmt.Println("ID:  ", data[i].ID)
		fmt.Println("Description:  ", data[i].Description)
		fmt.Println("Status:  ", data[i].Status)
		fmt.Println("CreatedAt:  ", data[i].CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("UpdatedAt:  ", data[i].UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
}

// update status
func cmdUpdate(args []string) {
	if len(args) < 2 {
		fmt.Println("provide the status ")
	}
	// for instance : args = [1,"done"]
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error while loading the file:  ", err)
	}
	status := args[1]

	if status != statusToDo && status != statusDone && status != statusInProgress {
		fmt.Println("invalied status")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Error while loading the file:  ", err)
		return
	}
	task, index := getbyID(tasks, id)
	if index == -1 {
		fmt.Printf("task with ID %d not found !", index)
		return
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	// save to file
	if err := saveTasks(tasks); err != nil {
		fmt.Println("Error saving the tasks")
		return
	}
	fmt.Println("task updated successfully ")

}
