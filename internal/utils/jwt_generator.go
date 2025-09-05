package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("your_secret_key")

func GenerateJwtToken(email string, isRefresh bool) (string, error) {
	var expirationTime time.Time
	if isRefresh {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	} else {
		expirationTime = time.Now().Add(15 * time.Minute)
	}
	claims := &jwt.MapClaims{
		"sub":      email,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenSignatureInvalid
}

func InvalidateToken(tokenStr string) error {
	_, err := ValidateToken(tokenStr)
	if err != nil {
		return err
	}
	// Invalidate the token (e.g., by adding it to a blacklist)
	return nil
}