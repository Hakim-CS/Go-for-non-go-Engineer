package main

import "fmt"

var isConnected bool = false

func Connect() {
	if isConnected == true {
		fmt.Println("Database connected ")
	}
}

func main() {

}
