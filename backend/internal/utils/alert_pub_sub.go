// package utils

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github/mahirjain_10/sse-backend/backend/internal/types"

// 	"github.com/redis/go-redis/v9"
// )

// // PublishToPubSub publishes a message to a Redis Pub/Sub topic
// func PublishToPubSub(redisClient *redis.Client, ctx context.Context, topicName string, message types.UpdateActiveStatus) error {
// 	// Marshal the message to JSON

// 	messageJSON, err := json.Marshal(message)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal message: %w", err)
// 	}

// 	// Publish the message to the Redis topic
// 	if err := redisClient.Publish(ctx, topicName, messageJSON).Err(); err != nil {
// 		return fmt.Errorf("failed to publish message to Redis: %w", err)
// 	}

// 	log.Printf("Message published to topic '%s': %s", topicName, string(messageJSON))
// 	return nil
// }

// // SubscribeToPubSub subscribes to a Redis Pub/Sub topic and listens for messages
// func SubscribeToPubSub(redisClient *redis.Client, ctx context.Context, topicName string) {
// 	// Subscribe to the given topic
// 	subscription := redisClient.Subscribe(ctx, topicName)
// 	defer subscription.Close()
// 	channel := subscription.Channel()
// 	fmt.Println("alert pub sub recieved")
// 	// Loop to listen for messages on the channel
// 	for msg := range channel {
// 		// Unmarshal the message from JSON into the UpdateActiveStatus struct
// 		var message types.UpdateActiveStatus
// 		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
// 			log.Printf("Error unmarshalling message from Pub/Sub: %v", err)
// 			continue
// 		}
// 		// Call the API to handle the notification for the alert
// 		if err := InvokeAlertNotificationAPI(message); err != nil {
// 			log.Printf("Error invoking alert notification API: %v", err)
// 		}
// 	}
// }

package utils
