package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	var URL string = "https://github.com/Hakim-CS/web-scraper"
	//client request with 5 second timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// create a new request
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("the request return an error : ", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("request failed: ", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("response status code : ", resp.Status)

}
