package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string) (int, error) {
	key := os.Getenv("SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(key), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("user_id not found in token")
		}
		return int(uid), nil
	}

	return 0, errors.New("invalid token")
}
