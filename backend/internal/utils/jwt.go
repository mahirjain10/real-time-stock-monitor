package utils

import (
	"os"
	"time"

	"github/mahirjain_10/sse-backend/backend/internal/types"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": id,
			"exp": time.Now().Add(time.Minute * 30).Unix(),
		})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func VerifyToken(tokenString string) (*types.CustomClaims, error) {
	claims := &types.CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func RefreshToken(tokenString string) (string, error) {
	claims, err := VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": claims.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})

	return newToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
