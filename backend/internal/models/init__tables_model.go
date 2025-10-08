package models

import (
	"database/sql"
	"fmt"
)

// InitUserTable initializes the "user" table in the database.
func InitUserTable(db *sql.DB) error {
	// Query to create the user table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,  
		email VARCHAR(100) UNIQUE NOT NULL,  
		password CHAR(60) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	`

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create user table: %v", err)
	}

	fmt.Println("User table created successfully!")
	return nil
}

// InitStockAlertTable initializes the "stock_alert" table in the database.
func InitStockAlertTable(db *sql.DB) error {
	// Query to create the stock_alert table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS stock_alert (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		alert_name VARCHAR(50) NOT NULL,
		ticker VARCHAR(20) NOT NULL,
		current_fetched_price DECIMAL(6,2) NOT NULL,
		current_fetched_time DATETIME NOT NULL, 
		alert_condition VARCHAR(2) NOT NULL,
		alert_price DECIMAL(6,2) NOT NULL,
		is_active BOOLEAN DEFAULT TRUE NOT NULL,
		created_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	);
	`

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create stock_alert table: %v", err)
	}

	fmt.Println("Stock alert table created successfully!")
	return nil
}

func InitializeMonitorStockTable(db *sql.DB) error {
	// Query to create the monitor_stock table if it doesn't exist
	query := `
		CREATE TABLE IF NOT EXISTS monitor_stock(
			id VARCHAR(36) PRIMARY KEY,
			alert_id VARCHAR(36) NOT NULL,
			ticker VARCHAR(20) NOT NULL,
			is_active BOOLEAN DEFAULT TRUE NOT NULL,
			created_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (alert_id) REFERENCES stock_alert(id) ON DELETE CASCADE
		);
	`
	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create monitor_stock table: %v", err)
	}

	fmt.Println("Monitor stock table created successfully!")
	return nil
}

func InitializeFCMTokenTable(db *sql.DB) error {
	// Query to create the monitor_stock table if it doesn't exist
	query := `
		CREATE TABLE IF NOT EXISTS fcm_token(
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			fcm_token VARCHAR(255) NOT NULL,
			created_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_on DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		);
	`
	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create fcm_token table: %v", err)
	}

	fmt.Println("FCM token table created successfully!")
	return nil
}