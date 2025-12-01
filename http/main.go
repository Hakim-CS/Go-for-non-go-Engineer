package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	resp, err := http.Get("https://www.google.com")

	if err != nil {
		fmt.Printf("repont return an erro , %v", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("response status code : ", resp.StatusCode)
	fmt.Println("header : ", len(resp.Header))

	//resp.Header is a map (Go’s key-value dictionary).
	for key, value := range resp.Header {
		fmt.Printf("key: %s,  value: %s \n", key, value)
	}
	fmt.Println()
	fmt.Println("reposnse body : ", resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("body return an error : ", err)
		return
	}
	defer fmt.Println("body : ", string(body))

}
