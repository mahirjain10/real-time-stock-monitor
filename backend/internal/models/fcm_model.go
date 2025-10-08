package models

import (
	"database/sql"
	"fmt"

	"github/mahirjain_10/sse-backend/backend/internal/types"
)

func FindFCMTokenUsingFCMToken(app *types.App, fcmToken string) (types.FCMToken,error) {
	stmt := `SELECT id ,user_id ,fcm_token FROM fcm_token WHERE fcm_token = ?`
	var fcmTokenData types.FCMToken
	err := app.DB.QueryRow(stmt, fcmToken).Scan(&fcmTokenData.ID, &fcmTokenData.UserID, &fcmTokenData.FCMToken)
	if err == sql.ErrNoRows{
		return types.FCMToken{}, nil
	}
	if err != nil {
		return types.FCMToken{}, fmt.Errorf("error finding fcm token: %w", err)
	}
	return fcmTokenData, nil
}

func FindFCMTokenUsingUserID(app *types.App, userID string) (types.FCMToken,error) {
	stmt := `SELECT id, user_id ,fcm_token FROM fcm_token WHERE user_id = ?`
	var fcmTokenData types.FCMToken
	err := app.DB.QueryRow(stmt, userID).Scan(&fcmTokenData.ID, &fcmTokenData.UserID, &fcmTokenData.FCMToken)
	if err == sql.ErrNoRows{
		return types.FCMToken{}, nil
	}
	if err != nil {
		return types.FCMToken{}, fmt.Errorf("error finding fcm token: %w", err)
	}
	return fcmTokenData, nil
}
func InsertFCMToken(app *types.App, fcmToken types.FCMToken) error {
	query := `INSERT INTO fcm_token (id, user_id, fcm_token) VALUES (?, ?, ?)`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(fcmToken.ID, fcmToken.UserID, fcmToken.FCMToken)
	if err != nil {
		fmt.Printf("Error executing insert: %v\n", err)
		return fmt.Errorf("error inserting fcm token: %w", err)
	}
	return nil
}

// DeleteFCMToken deletes an FCM token for a specific user from the fcm_token table.
func DeleteFCMToken(app *types.App, fcmToken, userID string) error {

	query := `DELETE FROM fcm_token WHERE fcm_token = ? AND user_id = ?`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare DELETE statement for fcm_token: %w", err)
	}
	defer stmt.Close()

	// Execute query and check affected rows
	result, err := stmt.Exec(fcmToken, userID)
	if err != nil {
		return fmt.Errorf("failed to execute DELETE for fcm_token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no fcm_token found for user_id %s and token %s", userID, fcmToken)
	}

	return nil
}

// func UpdateFCMToken(app *types.App, fcmToken types.FCMToken) error {
// 	query := `UPDATE fcm_token SET user_id = ? WHERE fcm_token = ?`
// 	stmt, err := app.DB.Prepare(query)
// 	if err != nil {
// 		return fmt.Errorf("error preparing statement: %w", err)
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(fcmToken.UserID, fcmToken.FCMToken)
// 	if err != nil {
// 		fmt.Printf("Error executing update: %v\n", err)
// 		return fmt.Errorf("error updating fcm token: %w", err)
// 	}
// 	return nil
// }
