package events

import (
	"fmt"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"log"
	"net/http"
)

// StopMonitoring stops monitoring for a given alert
func StopMonitoring(app *types.App, alertID string, ticker string) error {
	app.SSEServer.Mu.Lock()
	defer app.SSEServer.Mu.Unlock()

	cancelFunc, exists := app.SSEServer.ActiveCtxMap[alertID]
	if !exists {
		return fmt.Errorf("AlertID not found or already stopped")
	}

	// Cancel the monitoring context
	cancelFunc()

	// Remove client and context from maps
	delete(app.SSEServer.ClientsMap, alertID)
	delete(app.SSEServer.ActiveCtxMap, alertID)

	// Remove clientChan from ActiveTickersMap
	if clientChans, ok := app.SSEServer.ActiveTickersMap[ticker]; ok {
		newClientChans := make([]*chan string, 0, len(clientChans))
		for _, ch := range clientChans {
			if ch != app.SSEServer.ClientsMap[alertID] {
				newClientChans = append(newClientChans, ch)
			}
		}
		if len(newClientChans) == 0 {
			delete(app.SSEServer.ActiveTickersMap, ticker)
		} else {
			app.SSEServer.ActiveTickersMap[ticker] = newClientChans
		}
	}

	log.Println("Stopped monitoring for alert:", alertID, "ticker:", ticker)
	return nil
} 

func ClientDisconnect(alertID string, ticker string) {
	// var stockResponse stockResponse
	// fmt.Println("here we are on line 10")
	url := fmt.Sprintf("http://localhost:8080/disconnect?alertID=%s&ticker=%s", alertID, ticker)
	fmt.Println(`url from disconnect : `, url)
	resp, err := http.Get(url)
	fmt.Println("here we are on line 13")
	fmt.Println("response from get : ", resp)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Disconnected successfully")
	}


	// scanner := bufio.NewScanner(resp.Body)
	// for scanner.Scan() {
	// 	fmt.Println("printing here")
	// 	fmt.Println(scanner.Text()) // Print received SSE event
	// 	json.Unmarshal([]byte(scanner.Text()),&stockResponse)
	// 	// utils.ComparePriceAndThreshold(app.RedisClient,alertID,stockResponse.Price)
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Println("Error reading from server:", err)
	// }
}
