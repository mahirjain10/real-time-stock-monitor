package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github/mahirjain_10/sse-backend/backend/internal/events"
	"github/mahirjain_10/sse-backend/backend/internal/helpers"
	model "github/mahirjain_10/sse-backend/backend/internal/models"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/utils"
	"github/mahirjain_10/sse-backend/backend/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context, r *gin.Engine, app *types.App) {
	var user types.RegisterUser
	// Bind the incoming JSON request to 'user' struct
	if !helpers.BindAndValidateJSON(c, &user) {
		return
	}

	// DEBUG : Printing user and user.ID
	fmt.Println(user.ID)
	fmt.Println(user)

	vError := validator.ValidateRegisterUser(user)
	fmt.Println("error ", vError)

	fmt.Println(len(vError))
	if len(vError) != 0 {
		helpers.SendResponse(c, http.StatusBadRequest, "Validation error", nil, vError, false)
		return
	}
	// Check if a user with the given email already exists
	// retrievedUser, err := model.FindUserByEmail(app, user.Email)
	// if err != nil {
	// 	helpers.SendResponse(c, http.StatusInternalServerError, "Error while checking user existence", nil, nil, false)
	// 	return
	// }

	// if retrievedUser.ID != "" {
	// 	helpers.SendResponse(c, http.StatusConflict, "User already exists with given email", nil, nil, false)
	// 	return
	// }

	// Set a new UUID for the user
	user.ID = uuid.New().String()
	fmt.Printf("Generated UUID: %s\n", user.ID)

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error while hashing the password: %v\n", err)
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error during password hashing", nil, nil, false)
		return
	}
	user.Password = string(hashedPassword)

	// Save the new user to the database
	err = model.InsertUser(app, user)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062 (23000): Duplicate entry") {
			helpers.SendResponse(c, http.StatusConflict, "User already exists with given email", nil, nil, false)
			return
		}
		helpers.SendResponse(c, http.StatusInternalServerError, "Error while saving the user", nil, nil, false)
		return
	}

	// Respond with success if the user was created
	helpers.SendResponse(c, http.StatusCreated, "User account created successfully", nil, nil, true)

}

// Supports multiple FCM tokens per user, ensuring each token is linked to one user.
func LoginUser(c *gin.Context, r *gin.Engine, app *types.App) {
	var user types.LoginUser
	var fcmToken types.FCMToken

	// Unmarshal and validate JSON
	if !helpers.BindAndValidateJSON(c, &user) {
		return
	}

	// Validate user object
	vError := validator.ValidateLoginUser(user)
	if len(vError) != 0 {
		helpers.SendResponse(c, http.StatusBadRequest, "Validation error", nil, vError, false)
		return
	}

	// Find user by email
	retrievedUser, err := model.FindUserByEmail(app, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.SendResponse(c, http.StatusNotFound, "User with given email not found, please create a new account", nil, nil, false)
			return
		}
		helpers.SendResponse(c, http.StatusInternalServerError, "Error checking user existence", nil, nil, false)
		return
	}

	// Check if user exists
	if retrievedUser.ID == "" {
		helpers.SendResponse(c, http.StatusNotFound, "User with given email not found, please create a new account", nil, nil, false)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			helpers.SendResponse(c, http.StatusUnauthorized, "Incorrect password", nil, nil, false)
			return
		}
		helpers.SendResponse(c, http.StatusInternalServerError, "Error comparing passwords", nil, nil, false)
		return
	}

	// Handle FCM token
	if user.FcmToken != "" {
		// Check if FCM token exists
		fcmTokenData, err := model.FindFCMTokenUsingFCMToken(app, user.FcmToken)
		if err != nil && err != sql.ErrNoRows {
			helpers.SendResponse(c, http.StatusInternalServerError, "Error checking FCM token", nil, nil, false)
			return
		}

		// If token exists
		if fcmTokenData.ID != "" {
			// Case 1: Token already linked to this user, no action needed
			if fcmTokenData.UserID == retrievedUser.ID {
				log.Println("FCM token already linked to user", retrievedUser.ID)
			} else {
				// Case 2: Token linked to another user, reassign it
				err = model.DeleteFCMToken(app, user.FcmToken, fcmTokenData.UserID)
				if err != nil {
					helpers.SendResponse(c, http.StatusInternalServerError, "Error reassigning FCM token", nil, nil, false)
					return
				}
				fcmToken.ID = uuid.New().String()
				fcmToken.UserID = retrievedUser.ID
				fcmToken.FCMToken = user.FcmToken
				if err := model.InsertFCMToken(app, fcmToken); err != nil {
					helpers.SendResponse(c, http.StatusInternalServerError, "Error inserting FCM token", nil, nil, false)
					return
				}
				log.Println("Reassigned FCM token to user", retrievedUser.ID)
			}
		} else {
			// Case 3: Token doesn’t exist, insert new token
			fcmToken.ID = uuid.New().String()
			fcmToken.UserID = retrievedUser.ID
			fcmToken.FCMToken = user.FcmToken
			if err := model.InsertFCMToken(app, fcmToken); err != nil {
				helpers.SendResponse(c, http.StatusInternalServerError, "Error inserting FCM token", nil, nil, false)
				return
			}
			log.Println("Inserted new FCM token for user", retrievedUser.ID)
		}
	}

	// Generate and set auth token
	token, err := utils.CreateToken(retrievedUser.ID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Error generating auth token", nil, nil, false)
		return
	}
	utils.SetCookie(c, "auth_token", token, 1800) // 30-min expiry

	// Prepare response data
	data := map[string]interface{}{
		"email": retrievedUser.Email,
		"name":  retrievedUser.Name,
	}
	helpers.SendResponse(c, http.StatusOK, "User logged in successfully", data, nil, true)
}


// Logout User
// 1. Get the auth token from the cookie
// 2. Verify the token
// 3. Delete the auth token from the cookie
// 4. Get the user id from the token
// 5. Find alert which are active and are assscoiated with the user id
// 6. delete the cache 
// 7. stop the monitor stock
// 8. send the response
func LogoutUser(c *gin.Context, r *gin.Engine, app *types.App) {
	authToken, err := utils.GetCookie(c, "auth_token")
	if err != nil {
		helpers.SendResponse(c, http.StatusUnauthorized, "Authentication required", nil, nil, false)
		return
	}
	claims, err := utils.VerifyToken(authToken)
	if err != nil {
		helpers.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, nil, false)
		return
	}
	userID := claims.ID
	activeStocks, err := model.GetAllActiveStocksByUserId(app, userID)
	if err != nil {
		helpers.SendResponse(c, http.StatusInternalServerError, "Internal server error", nil, nil, false)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important: cancel the context when you're done

	for _, stockAlert := range activeStocks {
		events.ClientDisconnect(stockAlert.ID, stockAlert.Ticker.TickerToMonitor)

		deleted, err := app.RedisClient.Del(ctx, stockAlert.AlertID).Result()
		if err != nil {
			fmt.Printf("Error deleting cache for alert %s: %v\n", stockAlert.AlertID, err)
		} else if deleted > 0 {
			fmt.Printf("Successfully deleted key '%s' from Redis\n", stockAlert.AlertID)
		} else {
			fmt.Printf("Key '%s' not found in Redis\n", stockAlert.AlertID)
		}
	}

	utils.DeleteCookie(c, "auth_token")

	helpers.SendResponse(c, http.StatusOK, "Logout successful", nil, nil, true)
}
