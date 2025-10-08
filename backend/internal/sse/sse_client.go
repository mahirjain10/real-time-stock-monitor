package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type stockResponse struct {
	Price float64 `json:"price"`
	Time  string  `json:"time"`
}

func Client(alertID string, ticker string) {
	log.Println("alertID from client sse", alertID)
	log.Println("ticker from client sse", ticker)
	var stockResponse stockResponse
	fmt.Println("here we are on line 10")
	url := fmt.Sprintf("http://localhost:8080/events?alertID=%s&ticker=%s", alertID, ticker)
	log.Println("url from client sse", url)
	resp, err := http.Get(url)
	fmt.Println("here we are on line 13")
	fmt.Println("response from get : ", resp)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println("printing here")
		fmt.Println(scanner.Text()) // Print received SSE event
		json.Unmarshal([]byte(scanner.Text()), &stockResponse)
		// utils.ComparePriceAndThreshold(app.RedisClient,alertID,stockResponse.Price)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from server:", err)
	}
}

