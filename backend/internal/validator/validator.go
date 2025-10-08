package validator

import (
	"fmt"

	"github/mahirjain_10/sse-backend/backend/internal/types"

	"github.com/rezakhademix/govalidator"
)


func ValidateRegisterUser(user types.RegisterUser) map[string]string {
    v := govalidator.New() 
    
    // Email validation
    v.RequiredString(user.Email, "email", "Email is required").Email(user.Email, "email", "Invalid email")

    // Name validation with min/max character limits
    v.RequiredString(user.Name, "name", "Name is required").MaxString(user.Name, 20, "name", "Name must be at most 20 characters").MinString(user.Name, 3, "name", "Name must be at least 3 characters")

    // Password validation with min/max character limits
    v.RequiredString(user.Password, "password", "Password is required").MaxString(user.Password, 15, "password", "Password must be at most 15 characters").MinString(user.Password, 8, "password", "Password must be at least 8 characters")
    
    // Capture errors immediately after validation
    if v.IsFailed() {
        errors := v.Errors()
        fmt.Printf("Validation errors: %#v\n", errors)
        return errors
    }
    
    return map[string]string{}
}



func ValidateLoginUser(user types.LoginUser) map[string]string{
    v := govalidator.New() 

    // Email validation
    v.RequiredString(user.Email, "email", "Email is required").Email(user.Email, "email", "Invalid email")
    
    // Password validation with min/max character limits
    v.RequiredString(user.Password, "password", "Password is required").MaxString(user.Password, 15, "password", "Password must be at most 15 characters").MinString(user.Password, 8, "password", "Password must be at least 8 characters")
    
    // FCM Token validation
    v.RequiredString(user.FcmToken, "fcm_token", "FCM Token is required")
    // Capture errors immediately after validation
    if v.IsFailed() {
        errors := v.Errors()
        fmt.Printf("Validation errors: %#v\n", errors)
        return errors
    }
        
    return map[string]string{}
}
