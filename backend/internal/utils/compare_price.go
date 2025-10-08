package utils

import (
	"context"
	"fmt"

	// "github/mahirjain_10/sse-backend/backend/internal/sse"

	"github/mahirjain_10/sse-backend/backend/internal/types"
	"log"
	"strconv"
)

var (
	count int32 = 0
)

func ComparePriceAndThreshold(app *types.App, ctx context.Context, alertID string, currentPrice float64) {
	count++
	fmt.Println("alert ID  : ", alertID)
	alertData, err := app.RedisClient.HGetAll(ctx, alertID).Result()
	if err != nil {
		log.Printf("Error retrieving alert data from Redis: %v", err)
		return
	}
	fmt.Println("alert data from redis : ", alertData)
	alertPrice, err := strconv.ParseFloat(alertData["alert_price"], 64)
	if err != nil {
		log.Printf("Error parsing alert price: %v", err)
		return
	}

	fmt.Println(count)
	if count == 5 || count == 15 || count == 5 && alertData["ticker"] == "LICI.NS" {
		// if count == 20 && alertData["ticker"] == "LICI.NS" {
		currentPrice = alertPrice
		fmt.Println("current price := ", currentPrice)
	}

	fmt.Printf("current price : %f , alert price %f\n", currentPrice, alertPrice)
	isConditionMet, err := CompareUsingSymbol(alertData["alert_condition"], currentPrice, alertPrice)
	fmt.Println(isConditionMet)
	if err != nil {
		log.Printf("Error evaluating alert condition: %v", err)
		return
	}
	if isConditionMet {
		responseData := types.UpdateActiveStatus{
			UserID: alertData["user_id"],
			ID:     alertID,
			Active: false,
		}
		fmt.Println("response data : ", responseData)
	
		log.Printf("Alert condition met for alert ID: %s\n", alertID)
		go NotifyUserAboutAlert(app, responseData,alertData)

		fmt.Println("calling pub sub ")
		if err != nil {
			log.Printf("Error publishing to Pub/Sub: %v", err)
		}
	} else {
		log.Printf("Alert condition not met for alert ID: %s\n", alertID)
	}
}
