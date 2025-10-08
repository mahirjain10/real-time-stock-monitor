package models

import (
	"database/sql"
	"fmt"

	"github/mahirjain_10/sse-backend/backend/internal/types"
)

// Helper function to handle query execution and error handling
// Using vardic parameters to accept any type of parameter
func executeQueryRow(app *types.App, stmt string, args ...interface{}) (*sql.Row, error) {
	row := app.DB.QueryRow(stmt, args...)

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("query row execution error: %w", err)
	}
	return row, nil
}

// Helper function to scan a user into the struct
func scanUser(row *sql.Row) (types.RegisterUser, error) {
	var user types.RegisterUser
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return types.RegisterUser{}, fmt.Errorf("scanning user data failed: %w", err)
	}
	return user, nil
}

// FindUserByID retrieves a user by their unique ID
func FindUserByID(app *types.App, userID string) (types.RegisterUser, error) {
	stmt := `SELECT * FROM user WHERE id = ?`

	row, err := executeQueryRow(app, stmt, userID)
	if err != nil {
		return types.RegisterUser{}, err
	}

	user, err := scanUser(row)
	if err != nil {
		return types.RegisterUser{}, err
	}

	return user, nil
}

func FindUserByEmail(app *types.App, email string) (types.RegisterUser, error) {
    stmt := `SELECT * FROM user WHERE email = ?`
    
    row, err := executeQueryRow(app, stmt, email)
    if err != nil {
        return types.RegisterUser{}, fmt.Errorf("error executing query: %w", err)
    }

    var user types.RegisterUser
    err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt) // Modify fields as needed

	fmt.Println(err)
    if err == sql.ErrNoRows {
        return types.RegisterUser{}, sql.ErrNoRows
    } else if err != nil {
        return types.RegisterUser{}, fmt.Errorf("error scanning user: %w", err)
    }

    return user, nil
}


// InsertUser adds a new user to the database
func InsertUser(app *types.App, user types.RegisterUser) error {
	query := `
		INSERT INTO user(id, name, email, password) 
		VALUES(?, ?, ?, ?)
	`

	// Prepare the SQL statement
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement with user data
	if _, err := stmt.Exec(user.ID, user.Name, user.Email, user.Password); err != nil {
		return fmt.Errorf("error executing insert statement: %w", err)
	}

	return nil
}
