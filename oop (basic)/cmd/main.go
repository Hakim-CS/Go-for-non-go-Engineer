package main

import (
	"fmt"

	"github.com/example/oop/person" // Import the package
)

func main() {
	addr := person.NewAddress("turkey", "denizli", 212)
	person, err := person.NewPerson("hakim", 21, "hakim@gmail.com", *addr)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println(addr)
	fmt.Println("person's object ")
	fmt.Println(person)
	fmt.Println(person.IsAdult("Hakim"))
	//fmt.Println(person.ToString(person))


}
