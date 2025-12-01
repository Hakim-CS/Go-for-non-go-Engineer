package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetchRealWheather(city string) (WheatherRespond, error) {

	// do a request
	//get the resp
	// and render it in go struct format
	wheather_api := os.Getenv("WHEATHER_API_KEY")
	if wheather_api == "" {
		fmt.Print("env variable not specified")
		return WheatherRespond{}, nil
	}
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s",
		city, wheather_api)

	// make a http req
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error while creating req: ", err)
		return WheatherRespond{}, nil
	}
	req.Header.Set("user-agent", "wheather-api-practice")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("erroe while doing req: ", err)
		return WheatherRespond{}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var api_resp WheatherApiRespond
	if err := json.Unmarshal(body, &api_resp); err != nil {
		fmt.Println("error while unmarshalling ", err)
		return WheatherRespond{}, err
	}

	result := WheatherRespond{
		City:        api_resp.Address,
		Temperature: float32(api_resp.Days[0].Temp),
		Description: api_resp.Days[0].Wheather,
	}

	return result, nil

}
