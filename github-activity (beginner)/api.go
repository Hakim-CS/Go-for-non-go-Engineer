package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//"go.mongodb.org/mongo-driver/event"
)

func fetchEvents(userName string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", userName)

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "github-activuty-cli")
	// do the request and get the return value in resp
	resp, _ := client.Do(req)

	body, _ := io.ReadAll(resp.Body)

	// convert Json -> Go struct

	var event []Event
	if err := json.Unmarshal(body, &event); err != nil {
		fmt.Println("Err while unmarshaling resp: ", err)
		return nil, err
	}
	return event, nil
}
