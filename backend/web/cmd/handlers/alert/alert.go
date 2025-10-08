package alert

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github/mahirjain_10/sse-backend/backend/internal/helpers"
	"github/mahirjain_10/sse-backend/backend/internal/models"
	"github/mahirjain_10/sse-backend/backend/internal/sse"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/utils"
	"github/mahirjain_10/sse-backend/backend/internal/websocket"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCurrentStockPriceAndTime(c *gin.Context, r *gin.Engine, app *types.App) {
	// var stock types.GetCurrentPrice
	var TTM types.Ticker
	var stockData types.StockData

	if !helpers.BindAndValidateJSON(c, &TTM) {
		return
	}
	latestPrice, currentTime, err := utils.GetCurrentStockPriceAndTime(TTM.TickerToMonitor)
	if err != nil {
		if err.Error() == "failed to fetch stock price, try again" {
			helpers.SendResponse(c, http.StatusInternalServerError, "Failed to fetch stock price,try again", nil, nil, false)
			return
		}
		if err.Error() == "failed to decode response json" {
			c.JSON(http.StatusBadRequest, gin.H{"statusCode": http.StatusBadRequest, "message": "failed to decode response json", "error": err.Error()})
			return

		}
		if err.Error() == "no data found" {
			c.JSON(http.StatusNotFound, gin.H{"statusCode": http.StatusNotFound, "message": "no data found", "error": nil})
			return

		}
	}
	// Prepare response
	response := map[string]any{
		"statusCode": http.StatusOK,
		"message":    "Latest price fetched successfully",
		"data": types.GetCurrentPrice{
			CurrentFetchedPrice: latestPrice,
			CurrentFetchedTime:  currentTime,
		},
		"error": nil,
	}

	// Return the response
	fmt.Println(stockData)
	// c.JSON(http.StatusOK, response)
	helpers.SendResponse(c, http.StatusOK, "Current price fetched successfully", response, nil, true)
}

func CreateStockAlert(c *gin.Context, r *gin.Engine, app *types.App) {
	ctx := context.Background()
	var alertInput types.StockAlert
	// var monitorStockPrice types.MonitorStockPrice

	// Bind and validate JSON input
	if !helpers.BindAndValidateJSON(c, &alertInput) {
		return
	}

	// Check if the user exists
	user, err := models.FindUserByID(app, alertInput.UserID)
	if err != nil {
		log.Printf("Error finding user by ID: %v", err)
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	if user.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "User not found", nil, nil, false)
		return
	}

	// Check if alert name already exists for the user
	existingAlert, err := models.FindAlertNameByUserIDAndAlertName(app, alertInput.UserID, alertInput.AlertName)
	if err != nil {
		log.Printf("Error finding alert name: %v", err)
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	if existingAlert.ID != "" {
		helpers.SendResponse(c, http.StatusConflict, "Alert name already exists. Use a different name.", nil, nil, false)
		return
	}

	if alertInput.CurrentFetchedPrice == alertInput.AlertPrice {
		helpers.SendResponse(c, http.StatusBadRequest, "Alert price cannot be same as current price", nil, nil, false)
		return
	}
	// Generate a unique ID for the alert
	alertInput.ID = uuid.New().String()
	alertInput.Active = true
	// Insert stock alert data into the database
	if err := models.InsertStockAlertData(app, alertInput); err != nil {
		log.Printf("Error inserting stock alert data: %v", err)
		helpers.SendResponse(c, http.StatusInternalServerError, "Error saving stock alert", nil, nil, false)
		return
	}
	log.Printf("acitve status : %t", alertInput.Active)

	fcmToken, err := models.FindFCMTokenUsingUserID(app,alertInput.UserID)
	if err != nil {
		log.Printf("Error finding FCM token: %v", err)
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}
	// Save alert data in Redis
	alertData := map[string]interface{}{
		"fcm_token": 	fcmToken.FCMToken,
		"user_id":         user.ID,
		"ticker":          alertInput.TickerToMonitor,
		"alert_price":     alertInput.AlertPrice,
		"alert_condition": alertInput.Condition,

		"active":          strconv.FormatBool(alertInput.Active),
	}
	val, err := app.RedisClient.HSet(ctx, alertInput.ID, alertData).Result()
	if val == 0 {
		log.Println("Data could not saved in redis")
	}
	if err != nil {
		log.Printf("Error saving alert to Redis: %v\n", err)
	}

	// Insert stock monitoring data into database
	// monitorStockPrice.ID=uuid.NewString()
	// monitorStockPrice.AlertID=alertInput.ID
	// monitorStockPrice.TickerToMonitor=alertInput.TickerToMonitor
	// monitorStockPrice.IsActive=true

	// err = models.InsertMonitorStockData(app,monitorStockPrice)
	// if err != nil {
	// 	log.Printf("Error inserting stock monitoring data: %v", err)
	// 	helpers.SendResponse(c, http.StatusInternalServerError, "Error saving stock monitoring data", nil, nil, false)
	// 	return
	// }
	// monitorStockHashKey := "monitor_stock : " + monitorStockPrice.ID
	// monitorStockRedis := make(map[string]string)
	// monitorStockRedis["id"]=monitorStockPrice.ID
	// monitorStockRedis["alert_id"]=monitorStockPrice.AlertID
	// monitorStockRedis["ticker"]=monitorStockPrice.TickerToMonitor
	// monitorStockRedis["is_active"]=strconv.FormatBool(monitorStockPrice.IsActive)

	// fmt.Printf("Printing hash key : %s\n ",monitorStockHashKey)

	// val, err = app.RedisClient.HSet(ctx, monitorStockHashKey, monitorStockRedis).Result()
	// if val == 0 {
	// 	log.Println("Data could not saved in redis")
	// }
	// if err != nil {
	// 	log.Printf("Error saving stock monitoring data to Redis: %v\n", err)
	// 	return;
	// }

	// Publish alert to Redis channel
	utils.Publish(app.RedisClient, ctx, alertInput.TickerToMonitor, alertInput.ID)
    data:=make(map[string]interface{})
	fmt.Println(alertInput.ID)
	data["alert_id"]=alertInput.ID
	// Send success response
	go sse.Client(alertInput.ID, alertInput.TickerToMonitor)
	helpers.SendResponse(c, http.StatusCreated, "Stock alert created successfully", data, nil, true)
}

func UpdateStockAlert(c *gin.Context, r *gin.Engine, app *types.App) {
	ctx := context.Background()
	var updateAlertInput types.UpdateStockAlert
	if !helpers.BindAndValidateJSON(c, &updateAlertInput) {
		return
	}

	// Check if user exists or not
	user, err := models.FindUserByID(app, updateAlertInput.UserID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	if user.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "User not found", nil, nil, false)
		return
	}

	// Checking for alert data with given ID exists
	retrieveStockAlertData, err := models.FindAlertNameByUserIDAndAlertName(app, updateAlertInput.UserID, updateAlertInput.AlertName)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	fmt.Println(retrieveStockAlertData)

	//If alert name is already present in your account other than current alertID then send error
	if strings.TrimSpace(retrieveStockAlertData.ID) != strings.TrimSpace(updateAlertInput.ID) &&
		strings.TrimSpace(retrieveStockAlertData.UserID) == strings.TrimSpace(updateAlertInput.UserID) {
		fmt.Println("in if func")
		helpers.SendResponse(c, http.StatusNotFound, "Alert name already exists in your account,Use different alert name", nil, nil, false)
		return
	}

	err = models.UpdateStockAlertData(app, updateAlertInput)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Unable to update alert data ,Try again later", nil, nil, false)
		return
	}

	if retrieveStockAlertData.AlertPrice != updateAlertInput.AlertPrice {
		fmt.Println("in if for alert price update")
		// Update the data to redis
		val, err := app.RedisClient.HSet(ctx, updateAlertInput.ID, "alert_price", updateAlertInput.AlertPrice).Result()
		if val == 0 {
			log.Println("Data could not saved in redis")
		}
		if err != nil {
			// Log the error and return it or handle it as per your application's error handling policy
			log.Printf("Error updating alert in Redis for ID %s: %v", updateAlertInput.ID, err)
		}
	}
	if retrieveStockAlertData.Condition != updateAlertInput.Condition {
		fmt.Println("in if for alert condition update")
		// Update the data to redis
		val, err := app.RedisClient.HSet(ctx, updateAlertInput.ID, "alert_condition", updateAlertInput.Condition).Result()
		if val == 0 {
			log.Println("Data could not saved in redis")
		}
		if err != nil {
			// Log the error and return it or handle it as per your application's error handling policy
			log.Printf("Error updating alert in Redis for ID %s: %v", updateAlertInput.ID, err)
		}

	}
	helpers.SendResponse(c, http.StatusOK, "Stock alert updated successfully", nil, nil, true)
}

func DeleteStockAlert(c *gin.Context, r *gin.Engine, app *types.App) {
	ctx := context.Background()

	var deleteStockAlert types.DeleteStockAlert
	if !helpers.BindAndValidateJSON(c, &deleteStockAlert) {
		return
	}

	user, err := models.FindUserByID(app, deleteStockAlert.UserID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	if user.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "User not found", nil, nil, false)
		return
	}
	retrieveStockAlertData, err := models.FindAlertNameByUserIDAndID(app, deleteStockAlert.UserID, deleteStockAlert.ID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	fmt.Println(retrieveStockAlertData)
	if retrieveStockAlertData.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "Alert with given ID not found", nil, nil, false)
		return
	}

	rowsAffected, err := models.DeleteStockAlertByID(app, retrieveStockAlertData.UserID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}
	if rowsAffected == 0 {
		helpers.SendResponse(c, http.StatusNotFound, "Stock Alert to be deleted not found", nil, nil, false)
		return
	} else {
		_, err := app.RedisClient.Del(ctx, retrieveStockAlertData.ID).Result()
		if err != nil {
			log.Printf("Error deleting alert in Redis for ID %s: %v", retrieveStockAlertData.ID, err)
		}
		helpers.SendResponse(c, http.StatusOK, "Stock Alert deleted successfully", nil, nil, true)
		return
	}
}

func UpdateActiveStatus(c *gin.Context, r *gin.Engine, app *types.App) {
	ctx := context.Background()

	var updateActiveStatus types.UpdateActiveStatus
	if !helpers.BindAndValidateJSON(c, &updateActiveStatus) {
		return
	}

	isSuccess := utils.UpdateActiveStatusUtil(c, ctx, updateActiveStatus.UserID, updateActiveStatus.ID, updateActiveStatus.Active, app)

	if isSuccess {
		helpers.SendResponse(c, http.StatusOK, "Stock alert status updated successfully", nil, nil, true)
	}
}


// using websocket 
func StartStockAlertMonitoring(c *gin.Context, r *gin.Engine, app *types.App){
	ctx := context.Background()
    
	var startMonitoring types.StartMonitoring
	if !helpers.BindAndValidateJSON(c,&startMonitoring){
		return
	}

	success := utils.UpdateActiveStatusUtil(c,ctx,startMonitoring.UserID,startMonitoring.AlertID,true,app)
    // if err != nil{
	// 	helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
	// 	return
	// }
	if !success{
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}
	// utils.Publish(app.RedisClient, ctx, startMonitoring.TickerToMonitor, startMonitoring.AlertID)
	helpers.SendResponse(c, http.StatusOK, "Stock monitoring started successgfully", nil, nil, true)
    
}

func StopStockAlertMonitoring(c *gin.Context, r *gin.Engine, app *types.App,hub *websocket.Hub){
	ctx := context.Background()
    
	var startMonitoring types.StartMonitoring
	if !helpers.BindAndValidateJSON(c,&startMonitoring){
		return
	}

	success := utils.UpdateActiveStatusUtil(c,ctx,startMonitoring.UserID,startMonitoring.AlertID,false,app)
    // if err != nil{
	// 	helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
	// 	return
	// }
	if !success{
		fmt.Println("internal")
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}
	// utils.Publish(app.RedisClient, ctx, startMonitoring.TickerToMonitor, startMonitoring.AlertID)
	// hub.UnregisterClientByAlertID(startMonitoring.AlertID)
	helpers.SendResponse(c, http.StatusOK, "Stock monitoring started successgfully", nil, nil, true)
    
}


// func LoadMonitorActiveStocks(c *gin.Context, r *gin.Engine, app *types.App){
// 	if c.Query("user_id") == "" {
// 		helpers.SendResponse(c, http.StatusBadRequest, "User ID is required", nil, nil, false)
// 		return
// 	}
// 	// ctx := context.Background()

// 	activeStocks, err := models.GetAllActiveStocksByUserId(app, c.Query("user_id"))
// 	if err != nil {
// 		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
// 		return
// 	}
// 	for _ ,stockAlerts := range activeStocks{
			
// 	}
// 	helpers.SendResponse(c, http.StatusOK, "Active stocks loaded successfully", activeStocks, nil, true)
// }
func getCurrentTime() string {
	// You can implement a function to get the current time in your desired format
	return "2025-04-22 15:05:00 IST"
}
func SendFCMNotification(c *gin.Context, r *gin.Engine, app *types.App){
	ctx := context.Background()

	// var sendFCMNotification types.SendFCMNotification
	// if !helpers.BindAndValidateJSON(c,&sendFCMNotification){
		// 	return
		// }
	registrationToken:="cYvebaHvTZ2iB9q3QtJm5U:APA91bHVZEyY7ut-qZ_eNe25PmtOQKwoGTR4mxgBVSUl1MVL0zc-2img4UoobZFARkGYht5pUsz7-pABhA4GN-H63iGtlJx9f3NlT8oGsJX681ujbra9d8w"
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Go Backend Notification",
			Body:  "Notification sent from the Go backend at " + getCurrentTime(),
		},
		Token: registrationToken,
	}
	response, err := app.FCMClient.Send(ctx, message)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}
	log.Printf("Successfully sent message: %q\n", response)
	helpers.SendResponse(c, http.StatusOK, "FCM notification sent successfully", nil, nil, true)
}
