package person

import (
	"strconv"
	"strings"
)

type Person struct {
	name    string
	age     int
	email   string
	address Address
}

type Address struct {
	country string
	city    string
	zip     int
}

func NewAddress(country string, city string, zip int) *Address {
	result := Address{
		country: country,
		city:    city,
		zip:     zip,
	}
	return &result
}

func NewPerson(name string, age int, email string, address Address) (*Person, error) {
	result := Person{
		name:    name,
		age:     age,
		email:   email,
		address: address,
	}
	return &result, nil
}

func (person *Person) IsAdult(name string) bool {
	return person.age >= 18
}

func ToString(person Person) string {
	var myString []string = make([]string, 0)
	myString = append(myString, person.name)
	myString = append(myString, strconv.Itoa(person.age))
	myString = append(myString, person.email)
	return strings.Join(myString, ", ")
}
