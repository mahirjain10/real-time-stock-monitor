package helpers

import (
	"github/mahirjain_10/sse-backend/backend/internal/types"

	"github.com/gin-gonic/gin"
	// "github/mahirjain_10/sse-backend/backend/internal/types"
)

// SendResponse sends a standardized JSON response with the given status code, message, data, and error details
func SendResponse(c *gin.Context, statusCode int, message string, data map[string]interface{}, errors map[string]string, success bool) {
	// Create response struct
	response := types.Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
		Errors:  errors,
		Success: success,
	}

	// Send the response as JSON
	c.JSON(statusCode, response)
}
