package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 1 {
		fmt.Println("No username provided ! ")
		return
	}

	// github username
	// os.Args[0] is the program name
	userName := os.Args[1]
	fmt.Println(userName)

	event, _ := fetchEvents(userName)

	if len(event) == 0 {
		fmt.Println("no activity yet ")
		return
	}

	for _, e := range event {
		fmt.Println("- ", formatEvent(e))
	}
}

// roadMap: make URL
// https://api.github.com/users/<username>/events
// send a http request and get respond
// res is a slice of json , unmarshal it
// extract the required data from unmarshaled json
