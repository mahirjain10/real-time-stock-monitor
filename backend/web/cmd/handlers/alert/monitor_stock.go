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

// Manually start stock monitoring of a stock by user (SSE client)
func ManualStartStockMonitoring(c *gin.Context, r *gin.Engine, app *types.App) {
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

	// Find the alert by user id and alert id
	stockAlert, err := models.FindAlertNameByUserIDAndID(app, userID, c.Query("alert_id"))
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	//  If alert not found
	if stockAlert.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "Alert not found", nil, nil, false)
		return
	}

	// Update the active status of the alert to true
	err = models.UpdateActiveStatusByID(app, true, stockAlert.ID)
	log.Println("err while updating active status", err)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	// Find FCM token using user id
	fcmTokenData, err := models.FindFCMTokenUsingUserID(app, userID)
	log.Println("err while finding fcm token", err)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	// if error is nil and fcm token is empty
	if err == nil && fcmTokenData.FCMToken == "" {
		helpers.SendResponse(c, http.StatusNotFound, "FCM token not found", nil, nil, false)
		return
	}

	// Cache the data into redis
	alertData := map[string]interface{}{
		"fcm_token":       fcmTokenData.FCMToken,
		"user_id":         userID,
		"ticker":          stockAlert.Ticker.TickerToMonitor,
		"alert_price":     stockAlert.AlertPrice,
		"alert_condition": stockAlert.Condition,
		"active":          true,
	}

	_, err = app.RedisClient.HSet(ctx, stockAlert.ID, alertData).Result()
	// if val == 0 {
	// 	log.Println("Data could not saved in redis")
	// }
	if err != nil {
		log.Printf("Error saving alert to Redis: %v\n", err)
		return
	}

	log.Println("stockAlert id", stockAlert.ID)		
	log.Println("stockAlert ticker", stockAlert.Ticker.TickerToMonitor)

	// Start monitoring the stock
	go sse.Client(stockAlert.ID, stockAlert.Ticker.TickerToMonitor)
	helpers.SendResponse(c, http.StatusOK, "Stock monitoring started", nil, nil, true)
}

// Manually stop stock monitoring of a stock by user (SSE client)
func ManualStopStockMonitoring(c *gin.Context, r *gin.Engine, app *types.App) {
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

	// Find the alert by user id and alert id
	stockAlert, err := models.FindAlertNameByUserIDAndID(app, userID, c.Query("alert_id"))
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	//  If alert not found
	if stockAlert.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "Alert not found", nil, nil, false)
		return
	}

	// Update the active status of the alert to true
	err = models.UpdateActiveStatusByID(app, false, stockAlert.ID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	err = app.RedisClient.Del(ctx, stockAlert.ID).Err()
	if err != nil {
		log.Printf("Error deleting from Redis: %v\n", err)
	}

	// Introduce a 5-second delay
	// time.Sleep(5 * time.Second)

	// Stop monitoring the stock
	events.ClientDisconnect(stockAlert.ID, stockAlert.Ticker.TickerToMonitor)
	helpers.SendResponse(c, http.StatusOK, "Stock monitoring stopped", nil, nil, true)
}