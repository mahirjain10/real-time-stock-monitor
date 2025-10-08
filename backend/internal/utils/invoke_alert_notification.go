package utils

import (
	"context"
	"fmt"

	"github/mahirjain_10/sse-backend/backend/internal/events"
	"github/mahirjain_10/sse-backend/backend/internal/models"
	"github/mahirjain_10/sse-backend/backend/internal/types"

	"firebase.google.com/go/v4/messaging"
)

func NotifyUserAboutAlert(app *types.App, userData types.UpdateActiveStatus, alertDataFromRedis map[string]string) error {
	fmt.Println("notify user about alert")
	ctx := context.Background()

	// Use the alert data passed from the caller
	fmt.Println("alert data from caller: ", alertDataFromRedis)

	// 1. Update alert status in the database
	err := models.UpdateActiveStatusByID(app, false, userData.ID)
	if err != nil {
		return fmt.Errorf("failed to update alert status in DB: %v", err)
	}

	// 2. Stop monitoring before clearing Redis data
	fmt.Printf("Stopping monitoring for alertID: %s, ticker: %s\n", userData.ID, alertDataFromRedis["ticker"])
	err = events.StopMonitoring(app, userData.ID, alertDataFromRedis["ticker"])
	if err != nil {
		fmt.Printf("Warning: Failed to stop monitoring: %v\n", err)
	}

	// 3. Get FCM token from Redis
	fcmToken := alertDataFromRedis["fcm_token"]
	if fcmToken == "" {
		return fmt.Errorf("FCM token not found in Redis for alert ID: %s", userData.ID)
	}
	fmt.Println("Using FCM token from Redis:", fcmToken)

	// 4. Send FCM notification
	message := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: "Price Alert",
			Body:  fmt.Sprintf("Your price alert for %s has been triggered", alertDataFromRedis["ticker"]),
		},
		Data: map[string]string{
			"alert_id": userData.ID,
			"ticker":   alertDataFromRedis["ticker"],
			"price":    alertDataFromRedis["alert_price"],
		},
	}
	fmt.Println("Sending FCM message:", message)

	_, err = app.FCMClient.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send FCM notification: %v", err)
	}
	fmt.Println("FCM notification sent successfully")

	// 5. Remove the Redis cache for this alert
	err = app.RedisClient.Del(ctx, userData.ID).Err()
	if err != nil {
		fmt.Printf("Warning: Failed to clear Redis cache for alert %s: %v\n", userData.ID, err)
	}

	return nil
}