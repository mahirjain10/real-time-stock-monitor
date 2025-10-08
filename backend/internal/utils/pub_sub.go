package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
// 	socketServerURL = "ws://localhost:8080/ws/get-stock-price-socket"
	redisChannel    = "monitor"
)

// var count int = 0

// Publish sends a message to the Redis Pub/Sub channel.
func Publish(redisClient *redis.Client, ctx context.Context, ticker, alertID string) {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Error loading IST timezone:", err)
		return
	}
	currentTime := time.Now().In(ist).Format("Jan 02, 2006 at 3:04pm (MST)")

	message := map[string]string{
		"ticker_to_monitor": ticker,
		"alert_id":          alertID,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("[%s] - Error marshaling message: %s\n", currentTime, err)
		return
	}

	if err := redisClient.Publish(ctx, redisChannel, string(messageJSON)).Err(); err != nil {
		log.Printf("[%s] - Error publishing message: %s\n", currentTime, err)
		return
	}

	log.Printf("[%s] - Message published: %s\n", currentTime, string(messageJSON))
}

func Subscribe(redisClient *redis.Client, ctx context.Context) {
	_, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Error loading IST timezone:", err)
		return
	}

	pubsub := redisClient.Subscribe(ctx, redisChannel)
	defer pubsub.Close()
	log.Println("Listening for messages...")

	// var conn *websocket.Conn
	// var mu sync.Mutex // Protects conn from race conditions

	// connectWebSocket := func() (*websocket.Conn, error) {
	// 	mu.Lock()
	// 	defer mu.Unlock()

	// 	newConn, _, err := websocket.DefaultDialer.Dial(socketServerURL, nil)
	// 	if err != nil {
	// 		log.Println("WebSocket dial error:", err)
	// 		return nil, err
	// 	}

	// 	log.Println("WebSocket connected successfully")

	// 	// 🛑 Ensure the reader starts BEFORE returning the connection
	// 	go startWebSocketReader(newConn, redisClient, ctx, ist)

	// 	return newConn, nil
	// }

	// // Initial WebSocket connection
	// for {
	// 	conn, err = connectWebSocket()
	// 	if err == nil {
	// 		break
	// 	}
	// 	log.Println("Retrying WebSocket connection in 5 seconds...")
	// 	time.Sleep(5 * time.Second)
	// }
	// defer conn.Close()

	// processRedisMessages(pubsub, &conn, connectWebSocket, ctx, ist, redisClient)
}

// // READ MESSAGE AND MONITORS THE PRICES USING FUNC
// func startWebSocketReader(conn *websocket.Conn, redisClient *redis.Client, ctx context.Context, ist *time.Location) {
// 	go func() {
// 		defer conn.Close()
// 		for {
// 			_, response, err := conn.ReadMessage()
// 			if err != nil {
// 				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
// 					log.Println("WebSocket closed normally")
// 				} else {
// 					log.Printf("WebSocket read error (likely disconnection): %v", err)
// 				}
// 				return
// 			}

// 			currentTime := time.Now().In(ist).Format("Jan 02, 2006 at 3:04pm (MST)")
// 			log.Printf("[%s] - WebSocket response: %s\n", currentTime, response)

// 			var parsedResponse types.Response
// 			if err := json.Unmarshal(response, &parsedResponse); err != nil {
// 				log.Printf("Error unmarshaling response: %v", err)
// 				continue
// 			}
// 			if parsedResponse.Data == nil {
// 				return
// 			}

// 			dataJSON, err := json.Marshal(parsedResponse.Data)
// 			if err != nil {
// 				log.Printf("Error marshaling response data: %v", err)
// 				return
// 			}

// 			var stockData types.GetCurrentPrice
// 			if err := json.Unmarshal(dataJSON, &stockData); err != nil {
// 				log.Printf("Error unmarshaling stock data: %v", err)
// 				return
// 			}

// 			go ComparePriceAndThreshold(redisClient, ctx, stockData.AlertID, float64(stockData.CurrentFetchedPrice))
// 		}
// 	}()
// }

// func ComparePriceAndThreshold(redisClient *redis.Client, ctx context.Context, alertID string, currentPrice float64) {
// 	count++
// 	fmt.Println("alert data from redis : ",alertID)
// 	alertData, err := redisClient.HGetAll(ctx, alertID).Result()
// 	if err != nil {
// 		log.Printf("Error retrieving alert data from Redis: %v", err)
// 		return
// 	}
// 	fmt.Println("alert data from redis : ",alertData)
// 	alertPrice, err := strconv.ParseFloat(alertData["alert_price"], 64)
// 	if err != nil {
// 		log.Printf("Error parsing alert price: %v", err)
// 		return
// 	}

// 	fmt.Println(count)
// 	// if count == 10 || count == 15 || count == 5 && alertData["ticker"] == "LICI.NS" {
// 	if count == 20 && alertData["ticker"] == "LICI.NS" {

// 		currentPrice = alertPrice
// 		fmt.Println("current price := ", currentPrice)
// 	}

// 	fmt.Printf("current price : %f , alert price %f\n", currentPrice, alertPrice)
// 	isConditionMet, err := CompareUsingSymbol(alertData["alert_condition"], currentPrice, alertPrice)
// 	fmt.Println(isConditionMet)
// 	if err != nil {
// 		log.Printf("Error evaluating alert condition: %v", err)
// 		return
// 	}
// 	if isConditionMet {
// 		responseData := types.UpdateActiveStatus{
// 			UserID: alertData["user_id"],
// 			ID:     alertID,
// 			Active: false,
// 		}
// 		fmt.Println("response data : ", responseData)
// 		log.Printf("Alert condition met for alert ID: %s\n", alertID)
// 		err := PublishToPubSub(redisClient, ctx, "alert-topic", responseData)
// 		fmt.Println("calling pub sub ")
// 		if err != nil {
// 			log.Printf("Error publishing to Pub/Sub: %v", err)
// 		}
// 	} else {
// 		log.Printf("Alert condition not met for alert ID: %s\n", alertID)
// 	}
// }

// // PROCESSES THE INCOMING MESSAGE FROM PUBLISH AND SENDS THE MESSAGE TO WEBSOCKET
// func processRedisMessages(pubsub *redis.PubSub, conn **websocket.Conn, reconnectFunc func() (*websocket.Conn, error), ctx context.Context, ist *time.Location, redisClient *redis.Client) {
// 	for {
// 		msg, err := pubsub.ReceiveMessage(ctx)
// 		if err != nil {
// 			log.Printf("Error receiving Redis message: %s\n", err)
// 			continue
// 		}

// 		currentTime := time.Now().In(ist).Format("Jan 02, 2006 at 3:04pm (MST)")
// 		log.Printf("[%s] - Received message: %s\n", currentTime, msg.Payload)

// 		log.Printf("connection from line 188 : %v", *conn)
// 		// Attempt to send message over WebSocket
// 		if (*conn) == nil {
// 			log.Println("Cannot send message: WebSocket connection is nil.")
// 			continue
// 		}

// 		err = (*conn).WriteMessage(websocket.TextMessage, []byte(msg.Payload))
// 		if err != nil {
// 			log.Println("WebSocket write error. Reconnecting...")
// 			log.Printf("error of write message : %v", err)

// 			(*conn).Close()
// 			*conn = nil // ✅ Prevent writing to a closed connection

// 			continue
// 		}
// 		// Attempt to reconnect
// 		for {
// 			newConn, err := reconnectFunc()
// 			if err == nil {
// 				*conn = newConn
// 				log.Println("Reconnected successfully, waiting before sending messages...")

// 				time.Sleep(2 * time.Second) // 🛑 Give some time for WebSocket readiness

// 				// ⚡ Re-subscribe to Redis after reconnecting
// 				pubsub = redisClient.Subscribe(ctx, redisChannel)
// 				log.Println("Re-subscribed to Redis channel")
// 				break
// 			}
// 			log.Println("Retrying WebSocket connection in 5 seconds...")
// 			time.Sleep(5 * time.Second)
// 		}
// 	}
// }
