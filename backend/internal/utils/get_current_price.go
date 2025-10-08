package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github/mahirjain_10/sse-backend/backend/internal/types"
)

// func GetCurrentStockPriceAndTime(TTM string) (float64,string,error) {
// 	// Fetch stock data from external API
// 	var stockData types.StockData;
// 	apiURL := fmt.Sprintf("%s%s?range=1d&interval=1m", os.Getenv("STOCK_API_URL"), TTM)
// // fmt.Println("API URL:", apiURL)

// 	fmt.Println("API URL:", apiURL)

// 	res, err := http.Get(apiURL)
// 	fmt.Println("error ",err)
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		// return
// 	}

// 	// Print raw response
// 	fmt.Println("Response:", string(body))
// 	if err != nil {
// 		// helpers.SendResponse(c, http.StatusInternalServerError, "Failed to fetch stock price,try again", nil, nil, false)
// 		return 0,"",fmt.Errorf("failed to fetch stock price, try again")
// 	}
// 	fmt.Println("befor res.body closes ")
// 	defer res.Body.Close() // Ensure the response body is closed after reading
// 	fmt.Println("after res.body closes ")

// 	fmt.Println("error from helper func : ",err)
// 	// Decode the JSON response into stockData struct
// 	if err := json.NewDecoder(res.Body).Decode(&stockData); err != nil {
// 		fmt.Println(err)
// 		// c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "failed to decode response json", "error": err.Error()})
// 		return 0,"",fmt.Errorf("failed to decode response json")
// 	}
// 	fmt.Println("stock data : ",stockData)
// 	// Check if we have valid data
// 	if len(stockData.Chart.Result) == 0 || len(stockData.Chart.Result[0].Indicators.Quote) == 0 || len(stockData.Chart.Result[0].Indicators.Quote[0].Close) == 0 {
// 		// c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "no data found", "error": nil})
// 		return 0,"",fmt.Errorf("no data found")
// 	}

// 	// Get the latest price and current time
// 	latestPrice := stockData.Chart.Result[0].Indicators.Quote[0].Close[len(stockData.Chart.Result[0].Indicators.Quote[0].Close)-1]
// 	currentTime := time.Now().Format("02-01-2006 15:04:05")
// 	fmt.Printf("printing price %f and time %s",latestPrice,currentTime)
// 	return latestPrice,currentTime ,nil
// }


func GetCurrentStockPriceAndTime(TTM string) (float64, string, error) {
	var stockData types.StockData

	// Construct API URL
	apiURL := fmt.Sprintf("%s%s?range=1d&interval=1m", os.Getenv("STOCK_API_URL"), TTM)
	fmt.Println("API URL:", apiURL)

	// Create HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout to avoid hanging requests
	}

	// Create new request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set Headers (Mimic a browser)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return 0, "", fmt.Errorf("failed to fetch stock price: %w", err)
	}
	defer res.Body.Close() // Ensure response body is closed

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0, "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Debugging: Print raw response
	// fmt.Println("Response Body:", string(body))

	// Decode JSON response
	if err := json.Unmarshal(body, &stockData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0, "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	// fmt.Println("stock data : ",stockData)
	// Validate stock data
	if len(stockData.Chart.Result) == 0 ||
		len(stockData.Chart.Result[0].Indicators.Quote) == 0 ||
		len(stockData.Chart.Result[0].Indicators.Quote[0].Close) == 0 {
		fmt.Println("No data found in API response")
		return 0, "", fmt.Errorf("no stock data found or stock ticker invalid")
	}

	// Get the latest stock price
	latestPrice := stockData.Chart.Result[0].Indicators.Quote[0].Close[len(stockData.Chart.Result[0].Indicators.Quote[0].Close)-1]
	currentTime := time.Now().Format("02-01-2006 15:04:05")

	fmt.Printf("Stock Price: %f, Time: %s\n", latestPrice, currentTime)
	return math.Floor(latestPrice*100)/100, currentTime, nil
}
