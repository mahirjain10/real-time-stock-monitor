package alert

import (
	"context"
	"github/mahirjain_10/sse-backend/backend/internal/events"
	"github/mahirjain_10/sse-backend/backend/internal/helpers"
	"github/mahirjain_10/sse-backend/backend/internal/models"
	"github/mahirjain_10/sse-backend/backend/internal/sse"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoadUserAlerts loads active alerts for a user into cache and starts monitoring.
func LoadUserAlerts(c *gin.Context, app *types.App) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Increased timeout for initial load
	defer cancel()

	// Get raw user id from context (assuming your auth middleware sets this)
	userIDRaw, exists := c.Get("user")
	if !exists {
		helpers.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, nil, false)
		return
	}

	// Convert user id into string
	userID, ok := userIDRaw.(string)
	if !ok {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	// Find all active alerts for the user
	activeAlerts, err := models.GetAllActiveStocksByUserId(app, userID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	if len(activeAlerts) == 0 {
		helpers.SendResponse(c, http.StatusOK, "No active alerts found for this user.", nil, nil, true)
		return
	}
	log.Println("activeAlerts", activeAlerts)
	// Find FCM token for the user
	fcmTokenData, err := models.FindFCMTokenUsingUserID(app, userID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}
	if fcmTokenData.FCMToken == "" {
		helpers.SendResponse(c, http.StatusNotFound, "FCM token not found for this user.", nil, nil, false)
		return
	}

	// Iterate through active alerts, cache them, and start monitoring
	for _, alert := range activeAlerts {
		alertData := map[string]interface{}{
			"fcm_token":       fcmTokenData.FCMToken,
			"user_id":         userID,
			"ticker":          alert.Ticker.TickerToMonitor,
			"alert_price":     alert.AlertPrice,
			"alert_condition": alert.Condition,
			"active":          true,
		}

		_, err := app.RedisClient.HSet(ctx, alert.ID, alertData).Result()
		if err != nil {
			log.Printf("Error saving alert %s to Redis: %v\n", alert.ID, err)
			// Consider whether to continue or return an error here based on your application's needs
			continue // For now, we'll log and try to continue with other alerts
		}

		log.Printf("Starting monitoring for alert ID: %s, Ticker: %s\n", alert.ID, alert.Ticker.TickerToMonitor)
		go sse.Client(alert.ID, alert.Ticker.TickerToMonitor) // Start monitoring in a goroutine
	}

	helpers.SendResponse(c, http.StatusOK, "Active alerts loaded and monitoring started.", nil, nil, true)
}



// UnloadUserAlerts removes user's alerts from cache and stops monitoring.
func UnloadUserAlerts(c *gin.Context, app *types.App) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get raw user id from context
	userIDRaw, exists := c.Get("user")
	if !exists {
		helpers.SendResponse(c, http.StatusUnauthorized, "Unauthorized", nil, nil, false)
		return
	}

	// Convert user id into string
	userID, ok := userIDRaw.(string)
	if !ok {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	// Find all active alerts for the user (to get their IDs for cache removal and stopping SSE)
	activeAlerts, err := models.GetAllActiveStocksByUserId(app, userID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error while fetching active alerts", nil, nil, false)
		return
	}

	// Iterate through active alerts and perform cleanup
	for _, alert := range activeAlerts {
		// Remove alert data from Redis cache
		deleted, err := app.RedisClient.Del(ctx, alert.ID).Result()
		if err != nil {
			log.Printf("Error deleting alert %s from Redis: %v\n", alert.ID, err)
			// Consider logging or handling this error appropriately
		}
		if deleted > 0 {
			log.Printf("Removed alert %s from Redis cache.\n", alert.ID)
		}

		// Stop SSE monitoring for the alert
		events.ClientDisconnect(alert.ID, alert.Ticker.TickerToMonitor)
		log.Printf("Stopped monitoring for alert ID: %s, Ticker: %s\n", alert.ID, alert.Ticker.TickerToMonitor)
	}

	helpers.SendResponse(c, http.StatusOK, "User alerts unloaded and monitoring stopped.", nil, nil, true)
}