package websocket

import (
	"context"
	"sync"
)

// Hub manages WebSocket clients and their mappings.
type Hub struct {
	mu               sync.RWMutex
	clients          map[*Client]bool     // Set of all active clients
	clientsMap       map[string]*Client   // Maps alertId/userId to client
	activeTickersMap map[string][]*Client // Maps ticker to a list of clients
	activeCtxMap     map[string]context.CancelFunc
	register         chan *Client  // Channel for registering new clients
	unregister       chan *Client  // Channel for unregistering clients
	quit             chan struct{} // Channel to stop the hub gracefully
}

// NewHub creates and returns a new Hub instance with initialized channels and maps.
func NewHub() *Hub {
	return &Hub{
		clients:          make(map[*Client]bool),
		clientsMap:       make(map[string]*Client),
		activeTickersMap: make(map[string][]*Client),
		activeCtxMap: make(map[string]context.CancelFunc),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		quit:             make(chan struct{}),
	}
}

// // Run starts the Hub and handles client registration, unregistration, and graceful shutdown.
// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case client := <-h.register:
// 			h.mu.Lock()
// 			h.clients[client] = true
// 			log.Printf("Client registered: %v", client)
// 			log.Printf("Current clients: %v", h.clients)
// 			h.mu.Unlock()

// 		case client := <-h.unregister:
// 			h.mu.Lock()
// 			if _, ok := h.clients[client]; ok {
// 				delete(h.clients, client)
// 				log.Printf("Client unregistered: %v", client)
// 				log.Printf("Current clients: %v", h.clients)

// 				// Clean up corresponding mappings in clientsMap
// 				for alertID, c := range h.clientsMap {
// 					if c == client {
// 						delete(h.clientsMap, alertID)
// 						log.Printf("Deleted client from clientsMap: %s", alertID)
// 					}
// 				}

// 				// Clean up corresponding mappings in activeTickersMap
// 				for ticker, clients := range h.activeTickersMap {
// 					for i, c := range clients {
// 						if c == client {
// 							h.activeTickersMap[ticker] = append(clients[:i], clients[i+1:]...)
// 							log.Printf("Removed client from activeTickersMap: %s", ticker)
// 							break
// 						}
// 					}
// 					if len(h.activeTickersMap[ticker]) == 0 {
// 						delete(h.activeTickersMap, ticker)
// 						log.Printf("Deleted ticker from activeTickersMap: %s", ticker)
// 					}
// 				}

// 				close(client.send)
// 			}
// 			h.mu.Unlock()

// 		case <-h.quit:
// 			h.mu.Lock()
// 			for client := range h.clients {
// 				close(client.send)
// 			}
// 			h.clients = nil
// 			h.clientsMap = nil
// 			h.activeTickersMap = nil
// 			log.Println("Hub stopped, all clients and maps cleared")
// 			h.mu.Unlock()
// 			return
// 		}
// 	}
// }

// // Stop stops the Hub gracefully by

// // UnregisterClientByAlertID removes a client based on its alert ID and closes its connection
// func (h *Hub) UnregisterClientByAlertID(alertID string) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	// Find the client associated with this alert ID
// 	if client, exists := h.clientsMap[alertID]; exists {
// 		// Send a WebSocket close frame before closing the connection
// 		err := client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing connection"))
// 		if err != nil {
// 			fmt.Println("Failed to send close message:", err)
// 		}

// 		// Trigger unregistration through the channel
// 		h.unregister <- client
// 		fmt.Println("unregistered from func 104")

// 		// Close the WebSocket connection properly
// 		client.conn.Close()
// 		fmt.Println("closed from func 107")

// 		// Signal the monitoring goroutine to stop
// 		select {
// 		case <-client.done:
// 			// Already closed, do nothing
// 		default:
// 			close(client.done)
// 		}
// 		fmt.Println("unregistered from func done 110")
// 	}

// 	// Stop the monitoring goroutine if it exists
// 	if cancelMonitor, exists := h.activeCtxMap[alertID]; exists {
// 		cancelMonitor()
// 		log.Printf("Before cancelling ,activeCtxMap : %v\n",h.activeCtxMap)
// 		log.Println("Cancelling context from line 133 hub")
// 		delete(h.activeCtxMap, alertID)
// 		log.Printf("After cancelling ,activeCtxMap : %v\n",h.activeCtxMap)
// 	}
// }

// new
