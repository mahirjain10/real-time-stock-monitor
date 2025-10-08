package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"github/mahirjain_10/sse-backend/backend/internal/helpers"
	"github/mahirjain_10/sse-backend/backend/internal/models"
	"github/mahirjain_10/sse-backend/backend/internal/types"
	"github/mahirjain_10/sse-backend/backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the user using a cookie that contains the user ID
func AuthMiddleware(app *types.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cookie by name
		token, err := utils.GetCookie(c, "auth_token")
		if err != nil {
			c.Abort()
			return
		}

		// fmt.Printf("User ID type : %T ", userID)
		// fmt.Println("User ID type : ", userID)

		jwtClaims, err := utils.VerifyToken(token)
        if err != nil {
            helpers.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil, nil, false)
            c.Abort()
            return
        }

		
		// Find user by ID
		user, err := models.FindUserByID(app, jwtClaims.ID)
		fmt.Println(user)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helpers.SendResponse(c, http.StatusUnauthorized, "User not found", nil, nil, false)
				c.Abort()
				return
			}
			helpers.SendResponse(c, http.StatusInternalServerError, "Error retrieving user", nil, nil, false)
			c.Abort()
			return
		}

		// If user exists, store user info in context for later use
		c.Set("user", user.ID)

		c.Next()
	}
}