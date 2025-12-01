package main

import (
	"fmt"
	"os"
)

// let's  implement this with only cmdAdd function

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No command provided.")
		return
	}

	//os.Args[0] is the program name like "task-tracker"
	//os.Args[1] is the command like "add"
	//os.Args[2:] are the arguments to the command like ["Buy","groceries"]
	command := os.Args[1]
	arg := os.Args[2:]

	switch command {
	case "add":
		cmdAdd(arg)
	case "delete":
		cmdDeleteByID(arg)
	case "list":
		cmdList()
	case "update":
		cmdUpdate(arg)

	default:
		fmt.Printf("Unknown command: %s\n", command)

	}

}
