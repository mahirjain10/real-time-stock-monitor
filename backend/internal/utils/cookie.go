package utils

import (
	"errors"
	"github/mahirjain_10/sse-backend/backend/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name, value string, age int) {
//    / This does NOT allow setting SameSite
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        MaxAge:   age,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode, // Explicitly setting Lax
    })
}


func GetCookie(c *gin.Context, name string) (string, error) {
    value, err := c.Cookie(name)
    if err != nil {
        if errors.Is(err, http.ErrNoCookie) {
            helpers.SendResponse(c, http.StatusUnauthorized, "Authentication required", nil, nil, false)
            return "", err
        }
        helpers.SendResponse(c, http.StatusInternalServerError, "Error retrieving authentication cookie", nil, nil, false)
        return "", err
    }
    return value, nil
}

func DeleteCookie(c *gin.Context, name string) {
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     name,
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode, // Explicitly setting Lax
    })
}
