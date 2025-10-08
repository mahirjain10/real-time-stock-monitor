package websocket

import (
	"time"

	"github.com/gorilla/websocket"
	// "github/mahirjain_10/sse-backend/backend/internal/types"
)

const (
	writeWait      = 10 * time.Second    // Time to wait for writing a message
	pongWait       = 30 * time.Second    // Time to wait before considering the connection dead if no Pong is received
	pingPeriod     = (pongWait * 9) / 10 // Ping interval (90% of pongWait)
	maxMessageSize = 512                 // Max allowed message size for WebSocket
)

var (
	newline = []byte{'\n'} // Byte array for newline
	space   = []byte{' '}  // Byte array for space
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub  *Hub            // Reference to the hub managing all clients
	conn *websocket.Conn // WebSocket connection for the client
	send chan []byte     // Channel for sending messages to the client
	done chan struct{}   // Channel for signaling that the client is done (for graceful shutdown)
}

// func (c *Client) ReadPump(ctx *gin.Context) {
// 	defer func() {
// 		c.hub.unregister <- c
// 		log.Printf("Client unregistered in ReadPump: %v", c)

// 		c.hub.mu.Lock()
// 		if c.conn != nil {
// 			c.conn.Close()
// 			c.conn = nil
// 		}
// 		c.hub.mu.Unlock()

// 		select {
// 		case <-c.done:
// 		default:
// 			close(c.done)
// 		}
// 	}()

// 	c.conn.SetReadLimit(maxMessageSize)
// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.conn.SetPongHandler(func(string) error {
// 		c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 		return nil
// 	})

// 	log.Printf("Client map at start of ReadPump: %v", c.hub.clientsMap)

// 	monitorCtx, cancelMonitor := context.WithCancel(context.Background())
// 	defer cancelMonitor()

// 	for {
// 		select {
// 		case <-c.done:
// 			log.Println("Client disconnected, exiting ReadPump")
// 			return
// 		default:
// 			log.Println("IN HERE")
// 			_, message, err := c.conn.ReadMessage()
// 			log.Println("AFTER READ PUMP ")
// 			if err != nil {
// 				log.Printf("WebSocket read error: %v", err)
// 				return
// 			}
// 			var monitoringData types.MonitorStockPrice
// 			fmt.Printf("consoling incoming message from websocket : %s",message)
// 			if err := json.Unmarshal(message, &monitoringData); err != nil {
// 				log.Printf("Invalid message: %v", err)
// 				continue
// 			}

// 			fmt.Printf("consoling monitor stock data after parsing : %v",monitoringData)

// 			c.hub.mu.Lock()
// 			c.hub.clientsMap[monitoringData.AlertID] = c
// 			log.Printf("Updated clientsMap: %v", c.hub.clientsMap)

// 			if _, exists := c.hub.activeTickersMap[monitoringData.TickerToMonitor]; !exists {
// 				c.hub.activeTickersMap[monitoringData.TickerToMonitor] = []*Client{c}
// 				c.hub.activeCtxMap[monitoringData.AlertID] = cancelMonitor
// 				log.Printf("Added to activeTickersMap: %s", monitoringData.TickerToMonitor)
// 				log.Printf("Updated activeCtxMap: %v", c.hub.activeCtxMap)

// 				go func() {
// 					c.monitorStockPrice(monitoringData, monitoringData.AlertID, monitorCtx)
// 				}()
// 			} else {
// 				c.hub.activeTickersMap[monitoringData.TickerToMonitor] = append(c.hub.activeTickersMap[monitoringData.TickerToMonitor], c)
// 				c.hub.activeCtxMap[monitoringData.AlertID] = cancelMonitor
// 				log.Printf("Appended to activeTickersMap: %s", monitoringData.TickerToMonitor)
// 			}
// 			c.hub.mu.Unlock()
// 		}
// 	}
// }

// func (c *Client) monitorStockPrice(monitoringData types.MonitorStockPrice, alertID string, monitorCtx context.Context) {
// 	tickerChan := time.NewTicker(2 * time.Second)
// 	defer tickerChan.Stop()

// 	for {
// 		select {
// 		case <-monitorCtx.Done():
// 			log.Println("Stopping monitoring for:", monitoringData.TickerToMonitor)

// 			c.hub.mu.Lock()
// 			if clients, exists := c.hub.activeTickersMap[monitoringData.TickerToMonitor]; exists {
// 				for i, client := range clients {
// 					if client == c {
// 						log.Println("Removing client from active ticker map")
// 						c.hub.activeTickersMap[monitoringData.TickerToMonitor] = append(clients[:i], clients[i+1:]...)
// 						// log.Println("Removing client from active ticker map")
// 						log.Println("actievtciker Map from IF : ", c.hub.activeTickersMap)
// 						break
// 					}
// 				}
// 				if len(c.hub.activeTickersMap[monitoringData.TickerToMonitor]) == 0 {
// 					delete(c.hub.activeTickersMap, monitoringData.TickerToMonitor)
// 					log.Println("actievtciker Map : ", c.hub.activeTickersMap)
// 				}
// 			}
// 			c.hub.mu.Unlock()
// 			return

// 		case <-tickerChan.C:
// 			currentPrice, currentTime, err := utils.GetCurrentStockPriceAndTime(monitoringData.Ticker, types.StockData{})
// 			if err != nil {
// 				log.Printf("Error fetching stock price: %v", err)
// 				continue
// 			}

// 			responseJSON, _ := json.Marshal(map[string]interface{}{
// 				"statusCode": http.StatusOK,
// 				"message":    "Latest price fetched successfully",
// 				"data": types.GetCurrentPrice{
// 					CurrentFetchedPrice: currentPrice,
// 					CurrentFetchedTime:  currentTime,
// 					AlertID:             alertID,
// 				},
// 			})
// 			fmt.Println("IN MONITOR STOCK : ",responseJSON)
// 			c.hub.mu.RLock()
// 			for _, client := range c.hub.activeTickersMap[monitoringData.TickerToMonitor] {
// 				select {
// 				case client.send <- responseJSON:
// 				default:
// 					log.Println("Client send buffer full, skipping message")
// 				}
// 			}
// 			c.hub.mu.RUnlock()
// 		}
// 	}
// }

// func (c *Client) WritePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		log.Println("WritePump stopping...")
// 		ticker.Stop()
// 		c.hub.mu.Lock()
// 		if c.conn != nil {
// 			c.conn.Close()
// 			c.conn = nil
// 			log.Println("WebSocket connection closed in WritePump")
// 		}
// 		c.hub.mu.Unlock()
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.send:
// 			fmt.Println("message received in write pump")
// 			c.hub.mu.Lock()
// 			if c.conn == nil {
// 				c.hub.mu.Unlock()
// 				return // Stop if connection is nil
// 			}
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			c.hub.mu.Unlock()

// 			if !ok {
// 				return
// 			}
// 			fmt.Println("writing mess")
// 			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
// 				log.Printf("WebSocket write error: %v", err)
// 				return
// 			}

// 		case <-ticker.C:
// 			c.hub.mu.Lock()
// 			if c.conn == nil {
// 				c.hub.mu.Unlock()
// 				return
// 			}
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			c.hub.mu.Unlock()

// 			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				log.Printf("WebSocket ping error: %v", err)
// 				return
// 			}
// 		}
// 	}
// }

// // UnregisterClientByAlertID will be used to remove a client based on its alert ID
// // func (h *Hub) UnregisterClientByAlertID(alertID string) {
// // Close the WebSocket connection
// // }

// func ServeWs(c *gin.Context, hub *Hub, w http.ResponseWriter, r *http.Request) {
// 	// Add check for existing connection
// 	// if oldClient, exists := hub.clientsMap[r.URL.Query().Get("alertId")]; exists {
// 	//     hub.UnregisterClientByAlertID(r.URL.Query().Get("alertId"))
// 	//     oldClient.conn.Close()
// 	// }

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}

// 	client := &Client{
// 		hub:  hub,
// 		conn: conn,
// 		send: make(chan []byte, 256),
// 		done: make(chan struct{}),
// 	}

// 	client.hub.register <- client

// 	fmt.Printf("clientMap : %v\n", client.hub.clientsMap)
// 	// // Add alert ID to client mapping
// 	// alertID := r.URL.Query().Get("alertId")
// 	// if alertID != "" {
// 	//     hub.mu.Lock()
// 	//     hub.clientsMap[alertID] = client
// 	//     hub.mu.Unlock()
// 	// }

// 	go client.WritePump()
// 	go client.ReadPump(c)
// }
