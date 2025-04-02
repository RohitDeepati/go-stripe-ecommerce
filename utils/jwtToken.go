package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("secretkey")

func GenerateJwtToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user"] = userID
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix() // 1 hour expiration

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey) // Sign with the secret key
}


func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure we're using the correct signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid signing method")
			}
			return secretKey, nil // Return the secret key used for signing
	})
	if err != nil {
			return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
