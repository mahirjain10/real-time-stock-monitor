package cron

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github/mahirjain_10/sse-backend/backend/internal/models"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/websocket"
)
	const redisChannel    = "monitor"

const maxRetries = 3
func StartMonitoringJob(app *types.App) error {
	ctx := context.Background()
	
	monitorStockData, err := models.GetAllActiveStocks(app)
	if err != nil {
		log.Printf("[StartMonitoringJob] Failed to fetch active stocks: %v", err)
		// Optionally, send an alert via email or monitoring service
		return err
	}
	
	fmt.Println("monitorStockData : ",monitorStockData)
	
	_ = app.RedisClient.Subscribe(ctx, redisChannel)
	log.Println("Re-subscribed to Redis channel")
	for _, stock := range monitorStockData {
		alertData := map[string]interface{}{
			"user_id":         stock.UserID,
			"ticker":          stock.TickerToMonitor,
			"alert_price":     stock.AlertPrice,
			"alert_condition": stock.Condition,
			"active":          strconv.FormatBool(stock.Active),
		}
		val, err := app.RedisClient.HSet(ctx, stock.ID, alertData).Result()
		if val == 0 {
			log.Println("Data could not saved in redis")
		}
		if err != nil {
			log.Printf("Error saving alert to Redis: %v\n", err)
		}
		fmt.Printf("Monitoring stock: %s\n", stock)
		// utils.Publish(app.RedisClient, ctx, stock.TickerToMonitor, stock.ID)
		slog.Info("Published ")

	}
	return nil
}

func StartMonitoringWithRetry(app *types.App) {
    retries := 0

    for retries < maxRetries {
        err := StartMonitoringJob(app) // ✅ Now correctly using the returned error
        if err == nil {
            log.Println("[StartMonitoringWithRetry] Monitoring started successfully.")
            return
        }

        log.Printf("[StartMonitoringWithRetry] Attempt %d failed: %v", retries+1, err)
        retries++
        time.Sleep(10 * time.Second) // Wait 10 seconds before retrying
    }

    log.Println("[StartMonitoringWithRetry] Max retries reached, monitoring not started.")
}

func StopMonitoringJob(app *types.App, hub *websocket.Hub) {
    ctx := context.Background()

    // 1. Fetch all relevant keys (Assuming active alerts are stored with a prefix "alert:")
    keys, err := app.RedisClient.Keys(ctx, "alert:*").Result()
    if err != nil {
        log.Printf("[StopMonitoringJob] Failed to fetch active alert keys: %v", err)
        return
    }

    // 2. Iterate over each alert key
    for _, key := range keys {
        result, err := app.RedisClient.HGetAll(ctx, key).Result()
        if err != nil {
            log.Printf("[StopMonitoringJob] Failed to fetch alert data for key %s: %v", key, err)
            continue
        }

        // 3. Extract "alert_id" and unregister client
        _, exists := result["alert_id"]
        if !exists {
            log.Printf("[StopMonitoringJob] Alert key %s missing 'alert_id'", key)
            continue
        }

        // hub.UnregisterClientByAlertID(alertID)
    }

    log.Println("[StopMonitoringJob] Monitoring stopped successfully.")
}

