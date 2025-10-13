package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	currencycache "s/currencyCache"
	"time"

	"github.com/joho/godotenv"
)

func apiUrl() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Errore loading .env file : %v", err)
	}
	api_key := os.Getenv("API_KEY")

	api_url := "http://api.exchangeratesapi.io/v1/latest"

	full_URL := api_url + "?access_key=" + api_key + "&format=1"

	return full_URL
}

type Currency struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
}

func PrettyJSON(data []byte) error {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		return err
	}

	if err := os.WriteFile("pretty_cache.json", pretty.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

var cache = currencycache.NewCache(24 * time.Hour)

func ApiResponse(url string) (Currency, error) {
	// check if it in cache
	var resp Currency
	fmt.Println("Using cached data:", url)
	if data, ok := cache.Get(url); ok {
		err := PrettyJSON(data)
		if err != nil {
			fmt.Println("Error prettifying JSON:", err)
		} else {
			fmt.Println("Pretty JSON written to pretty_cache.json")
		}
		err = json.Unmarshal(data, &resp)
		if err != nil {
			return resp, fmt.Errorf("%v", err)
		}
		return resp, nil
	}

	fmt.Println("Cache unmarshal failed, fetchingAPI...")

	req, err := http.Get(url)
	if err != nil {
		log.Fatalf("fatal : faild to get req %v", err)
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("fatal : faild to get response %v", err)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		panic(err)
	}
	cache.Add(url, body)
	return resp, nil
}

func GetRate(resp Currency, code string) (float64, bool) {
	rate, ok := resp.Rates[code]
	return rate, ok
}
