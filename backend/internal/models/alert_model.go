package models

import (
	"database/sql"
	"fmt"
	"time"

	"github/mahirjain_10/sse-backend/backend/internal/types"
	// "github/mahirjain_10/sse-backend/backend/internal/types"
)

// FindAlertNameByUserIDAndAlertName retrieves a stock alert by user ID and alert name.
func FindAlertNameByUserIDAndAlertName(app *types.App, userID, alertName string) (types.StockAlert, error) {
	var stockAlert types.StockAlert

	stmt := `
	SELECT 
		id,
		user_id,
		ticker,
		alert_name,
		current_fetched_price,
		current_fetched_time,
		alert_condition,
		alert_price,
		is_active,
		created_on,
		updated_on
	FROM stock_alert
	WHERE user_id = ? AND alert_name = ?;
	`

	err := app.DB.QueryRow(stmt, userID, alertName).Scan(
		&stockAlert.ID,
		&stockAlert.UserID,
		&stockAlert.Ticker.TickerToMonitor,
		&stockAlert.AlertName,
		&stockAlert.CurrentFetchedPrice,
		&stockAlert.CurrentFetchedTime,
		&stockAlert.Condition,
		&stockAlert.AlertPrice,
		&stockAlert.Active,
		&stockAlert.CreatedAt,
		&stockAlert.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.StockAlert{}, nil // Return empty if no rows found
		}
		fmt.Printf("Error while fetching data: %v\n", err)
		return types.StockAlert{}, err
	}

	return stockAlert, nil
}

// FindAlertNameByUserIDAndID retrieves a stock alert by user ID and alert ID.
func FindAlertNameByUserIDAndID(app *types.App, userID, alertID string) (types.StockAlert, error) {
	var stockAlert types.StockAlert

	stmt := `
	SELECT 
		id,
		user_id,
		ticker,
		alert_name,
		current_fetched_price,
		current_fetched_time,
		alert_condition,
		alert_price,
		is_active,
		created_on,
		updated_on
	FROM stock_alert
	WHERE user_id = ? AND id = ?;
	`

	err := app.DB.QueryRow(stmt, userID, alertID).Scan(
		&stockAlert.ID,
		&stockAlert.UserID,
		&stockAlert.Ticker.TickerToMonitor,
		&stockAlert.AlertName,
		&stockAlert.CurrentFetchedPrice,
		&stockAlert.CurrentFetchedTime,
		&stockAlert.Condition,
		&stockAlert.AlertPrice,
		&stockAlert.Active,
		&stockAlert.CreatedAt,
		&stockAlert.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.StockAlert{}, nil // Return empty if no rows found
		}
		fmt.Printf("Error while fetching data: %v\n", err)
		return types.StockAlert{}, err
	}

	return stockAlert, nil
}

// InsertStockAlertData inserts a new stock alert record into the database.
func InsertStockAlertData(app *types.App, stockAlertData types.StockAlert) error {
	// Parse the current fetched time (assumed to be in '06-03-2025 20:26:16' format)
	currentFetchedTime, err := time.Parse("02-01-2006 15:04:05", stockAlertData.CurrentFetchedTime)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return fmt.Errorf("failed to parse current fetched time: %w", err)
	}

	// Format the parsed time into the database-friendly format (YYYY-MM-DD HH:MM:SS)
	formattedTime := currentFetchedTime.Format("2006-01-02 15:04:05")
	
	query := `
	INSERT INTO stock_alert (
		user_id, id, alert_name, ticker, 
		current_fetched_price, current_fetched_time, 
		alert_condition, alert_price
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		stockAlertData.UserID,
		stockAlertData.ID,
		stockAlertData.AlertName,
		stockAlertData.Ticker.TickerToMonitor,
		stockAlertData.CurrentFetchedPrice,
		formattedTime,
		stockAlertData.Condition,
		stockAlertData.AlertPrice,
	)

	if err != nil {
		fmt.Printf("Error executing insert: %v\n", err)
		return fmt.Errorf("failed to insert stock alert data: %w", err)
	}

	return nil
}

// UpdateStockAlertData updates an existing stock alert record in the database.
func UpdateStockAlertData(app *types.App, updateData types.UpdateStockAlert) error {
	query := `
	UPDATE stock_alert
	SET alert_name = ?, alert_condition = ?, alert_price = ?
	WHERE user_id = ? AND id = ?
	`

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		updateData.AlertName,
		updateData.Condition,
		updateData.AlertPrice,
		updateData.UserID,
		updateData.ID,
	)

	if err != nil {
		fmt.Printf("Error executing update: %v\n", err)
		return fmt.Errorf("failed to update stock alert data: %w", err)
	}

	return nil
}

// DeleteStockAlertByID deletes a stock alert by its ID.
func DeleteStockAlertByID(app *types.App, alertID string) (int64, error) {
	tx, err := app.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `DELETE FROM stock_alert WHERE id = ?`
	result, err := tx.Exec(query, alertID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete stock alert data: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return rowsAffected, nil
}

// UpdateActiveStatusByID updates the 'is_active' status of a stock alert by ID.
func UpdateActiveStatusByID(app *types.App, status bool, alertID string) error {
	query := `
	UPDATE stock_alert
	SET is_active = ?
	WHERE id = ?
	`

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, alertID)
	if err != nil {
		fmt.Printf("Error executing update: %v\n", err)
		return fmt.Errorf("failed to update stock alert status: %w", err)
	}

	return nil
}

func InsertMonitorStockData(app *types.App, MSP types.MonitorStockPrice) error {
	query := `
		INSERT INTO monitor_stock(
			id, alert_id, ticker, is_active
		)
		VALUES (?, ?, ?, ?) 
	`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		MSP.ID,
		MSP.AlertID,
		MSP.TickerToMonitor,
		MSP.IsActive,
	)
	if err != nil {
		fmt.Printf("Error executing insert: %v\n", err)
		return fmt.Errorf("failed to insert monitor stock data: %w", err)
	}
	return nil
}

func ChangeStockMonitoringStatus(app *types.App, isActive bool, id string) error {
	query := `
		UPDATE monitor_stock 
		SET = is_active=?
		WHERE id = ?
	`
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		isActive,
		id,
	)
	if err != nil {
		fmt.Printf("Error executing update: %v\n", err)
		return fmt.Errorf("failed to update monitor stock data: %w", err)
	}
	return nil
}

func GetAllActiveStocks(app *types.App) ([]types.StockAlert, error) {
	var monitorStocks []types.StockAlert
	query := `
    SELECT id,
           user_id,
           ticker,
           alert_condition,
           alert_price,
           is_active
    FROM stock_alert 
    WHERE alert_name LIKE "c%";`

	rows, err := app.DB.Query(query)
	if err != nil {
		fmt.Printf("Error while fetching data: %v\n", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		var startMonitoring types.StockAlert
		err = rows.Scan(
			&startMonitoring.ID,
			&startMonitoring.UserID,
			&startMonitoring.TickerToMonitor,
			&startMonitoring.Condition,
			&startMonitoring.AlertPrice,
			&startMonitoring.Active,
		)
		if err != nil {
			fmt.Printf("Error while scanning rows: %v\n", err)
			return nil, err
		}
		monitorStocks = append(monitorStocks, startMonitoring)
	}
	return monitorStocks, nil
}


// GetAllActiveStocksByUserId retrieves all active stock alerts for a given user ID.
func GetAllActiveStocksByUserId(app *types.App, userID string) ([]types.StockAlert, error) {
	var activeStocks []types.StockAlert

	// query to get all the active stocks's details using user id
	query := `
	SELECT
		id,
		user_id,
		ticker,
		alert_condition,
		alert_price,
		is_active,
		created_on,
		updated_on
	FROM stock_alert
	WHERE user_id = ? AND is_active = true;
`


	// prepare the statement
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// execute the query
	rows, err := stmt.Query(userID)
	if err != nil {
		fmt.Printf("Error while querying data: %v\n", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stockAlert types.StockAlert
		// var tickerToMonitor string // Temporary variable to scan ticker

		err := rows.Scan(
			&stockAlert.ID,
			&stockAlert.UserID,
			&stockAlert.TickerToMonitor,
			&stockAlert.Condition,
			&stockAlert.AlertPrice,
			&stockAlert.Active,
			&stockAlert.CreatedAt,
			&stockAlert.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		} // Assign the scanned ticker
		activeStocks = append(activeStocks, stockAlert)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error during rows iteration: %v\n", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return activeStocks, nil
}