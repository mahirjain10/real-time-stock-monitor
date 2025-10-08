package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github/mahirjain_10/sse-backend/backend/internal/events"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/utils"
)

func NewSSEServer() *types.SSEServer {
	return &types.SSEServer{
		// clientChan:       make(chan string),
		ClientsMap:       make(map[string]*chan string), // map of alertID to client channel
		ActiveCtxMap:     make(map[string]context.CancelFunc), // map of alertID to monitoring context
		ActiveTickersMap: make(map[string][]*chan string), // map of ticker to list of client channels
		Quit:             make(chan struct{}), // channel to quit the server
	}
}

func SSEHandler(app *types.App, w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Get alertID and ticker from query parameters
	query := r.URL.Query()
	alertID := query.Get("alertID")
	ticker := query.Get("ticker")

	if alertID == "" || ticker == "" {
		http.Error(w, "Missing alertID or ticker", http.StatusBadRequest)
		return
	}

	fmt.Println("New SSE connection for alert:", alertID, "ticker:", ticker)

	clientChan := make(chan string)
	monitorCtx, cancelMonitor := context.WithCancel(context.Background())

	// Store the client in the map
	app.SSEServer.Mu.Lock()
	app.SSEServer.ClientsMap[alertID] = &clientChan
	app.SSEServer.ActiveTickersMap[ticker] = append(app.SSEServer.ActiveTickersMap[ticker], &clientChan)
	app.SSEServer.ActiveCtxMap[alertID] = cancelMonitor
	app.SSEServer.Mu.Unlock()

	defer func() {
		// Cleanup when client disconnects
		app.SSEServer.Mu.Lock()
		delete(app.SSEServer.ClientsMap, alertID)
		delete(app.SSEServer.ActiveCtxMap, alertID)

		// Remove clientChan from ActiveTickersMap
		if clientChans, ok := app.SSEServer.ActiveTickersMap[ticker]; ok {
			newClientChans := make([]*chan string, 0, len(clientChans))
			for _, ch := range clientChans {
				if ch != &clientChan {
					newClientChans = append(newClientChans, ch)
				}
			}
			if len(newClientChans) == 0 {
				delete(app.SSEServer.ActiveTickersMap, ticker)
			} else {
				app.SSEServer.ActiveTickersMap[ticker] = newClientChans
			}
		}
		app.SSEServer.Mu.Unlock()

		cancelMonitor()
		close(clientChan)
		log.Println("Stopped monitoring for:", ticker)
	}()

	tickerChan := time.NewTicker(2 * time.Second)
	defer tickerChan.Stop()

	// Start monitoring stock prices in a separate goroutine
	go func() {
		for {
			select {
			case <-monitorCtx.Done():
				log.Printf("Price fetching goroutine for %s stopped", ticker)
				return
			case <-tickerChan.C:
				// Fetch stock price
				currentPrice, currentTime, err := utils.GetCurrentStockPriceAndTime(ticker)
				if err != nil {
					log.Printf("Error fetching stock price for %s: %v", ticker, err)
					continue
				}
				log.Printf("Fetched stock price for %s: %f at %v", ticker, currentPrice, currentTime)

				// Compare price and threshold
				go utils.ComparePriceAndThreshold(app, monitorCtx, alertID, currentPrice)

				// Send stock update only if context is still active
				message := fmt.Sprintf("data: {\"price\": %v, \"time\": \"%v\"}\n\n", currentPrice, currentTime)
				select {
				case <-monitorCtx.Done():
					log.Printf("Context canceled, skipping send for %s", ticker)
					return
				case clientChan <- message:
					// Message sent successfully
				}
			}
		}
	}()

	// Stream data to client
	for {
		select {
		case <-monitorCtx.Done():
			return
		case msg := <-clientChan:
			_, err := fmt.Fprintf(w, msg)
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func StopMonitoringHandler(app *types.App, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	alertID := query.Get("alertID")
	ticker := query.Get("ticker")

	if alertID == "" || ticker == "" {
		http.Error(w, "Missing alertID or ticker", http.StatusBadRequest)
		return
	}

	err := events.StopMonitoring(app, alertID, ticker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error(), "statusCode": "404"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Stopped monitoring for alert", "statusCode": "200"})
}