package main

type WheatherApiRespond struct {
	Address string `json:"address"`
	Days    []struct {
		Temp     float64 `json:"temp"`
		Wheather string  `json:"conditions"`
	} `json:"days"`
}

type WheatherRespond struct {
	City        string  `json:"city"`
	Temperature float32 `json:"temperature"`
	Description string  `json:"description"`
}
