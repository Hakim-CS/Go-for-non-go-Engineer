package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"net/http"
	//"os"
)

// ✅ GLOBAL - created once, shared across all requests
//var cache = NewCache() // create a new cache instance
// changed to redis cache

var redisClient = NewRedisClient()

func cacheSet(city string, data WheatherRespond) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// store with 12h expiration
	return redisClient.Set(ctx, city, jsonData, 12*time.Hour).Err()
}

func cacheGet(city string) (WheatherRespond, bool) {
	val, err := redisClient.Get(ctx, city).Result()
	if err == redis.Nil {
		return WheatherRespond{}, false
	}
	if err != nil {
		return WheatherRespond{}, false
	}

	var resp WheatherRespond
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return WheatherRespond{}, false
	}

	return resp, true
}

// wheather respond struct

// wheather handler func
func wheaterHandler(w http.ResponseWriter, r *http.Request) {
	// r : is the full http reques that client made
	// /wheather?city=london
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, " not city name provided ! ", http.StatusBadRequest)
		return
	}

	// check in cache first
	if resp, found := cacheGet(city); found {
		//json.NewEncoder(w).Encode(resp)
		fmt.Fprintf(w, "Data retrived from Cache redis server: %s\n", resp.City)
		fmt.Fprintf(w, "City: %s\n", resp.City)
		fmt.Fprintf(w, "Temperature: %.2f°C\n", resp.Temperature)
		fmt.Fprintf(w, "Description: %s\n", resp.Description)
		return
	}

	resp, _ := fetchRealWheather(city)
	//json.NewEncoder(w).Encode(resp)

	// store in cache
	//cache.Set(city, resp, 10*time.Minute)

	_ = cacheSet(city, resp)

	// write resp in a user friendly format
	fmt.Fprintf(w, "Data retrived from Real API server: %s \n", resp.City) // how can i write this on webpage
	// instead of console : answer
	fmt.Fprintf(w, "City: %s\n", resp.City)
	fmt.Fprintf(w, "Temperature: %.2f°C\n", resp.Temperature)
	fmt.Fprintf(w, "Description: %s\n", resp.Description)

}

// main func
func main() {
	http.HandleFunc("/wheather", wheaterHandler)
	fmt.Println("server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
