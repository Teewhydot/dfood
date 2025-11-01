package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

// InitJWT initializes the JWT secret key from configuration
func InitJWT(secret string) error {
	if secret == "" {
		return errors.New("JWT secret cannot be empty")
	}
	if len(secret) < 16 {
		return errors.New("JWT secret must be at least 16 characters long for security")
	}
	if secret == "your_secret_key" || secret == "your_jwt_secret_key_here" {
		return errors.New("JWT secret must be changed from default value")
	}
	jwtKey = []byte(secret)
	return nil
}

// RotateSecretKey changes the JWT secret, invalidating all existing tokens
func RotateSecretKey(newSecret string) {
	jwtKey = []byte(newSecret)
	// Clear blacklist since all tokens are now invalid anyway
	tokenBlacklist = make(map[string]bool)
}

func GenerateJwtToken(userID, email string, isRefresh bool) (string, error) {
	var expirationTime time.Time
	if isRefresh {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	} else {
		expirationTime = time.Now().Add(15 * time.Minute)
	}
	claims := &jwt.MapClaims{
		"sub":        email,
		"user_id":    userID,
		"exp":        expirationTime.Unix(),
		"iat":        time.Now().Unix(),
		"is_refresh": isRefresh,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	// Check if token is blacklisted first
	if IsTokenBlacklisted(tokenStr) {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenSignatureInvalid
}

// In-memory blacklist (in production, use Redis or database)
var tokenBlacklist = make(map[string]bool)

func InvalidateToken(tokenStr string) error {
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return err
	}

	// Add token to blacklist
	tokenBlacklist[tokenStr] = true

	// Optional: Store with expiration time for cleanup
	// In production, you'd store this in Redis with TTL
	_ = claims
	return nil
}

func IsTokenBlacklisted(tokenStr string) bool {
	return tokenBlacklist[tokenStr]
}

// InvalidateAllUserTokens invalidates all tokens for a specific user
func InvalidateAllUserTokens(email string) {
	// In a real implementation, you'd query your token storage
	// and invalidate all tokens for this user
	for token := range tokenBlacklist {
		claims, err := ValidateTokenWithoutBlacklistCheck(token)
		if err == nil {
			if sub, ok := (*claims)["sub"].(string); ok && sub == email {
				tokenBlacklist[token] = true
			}
		}
	}
}

// ValidateTokenWithoutBlacklistCheck validates token without checking blacklist
func ValidateTokenWithoutBlacklistCheck(tokenStr string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenSignatureInvalid
}

// CleanupExpiredTokensFromBlacklist removes expired tokens from blacklist
func CleanupExpiredTokensFromBlacklist() {
	for token := range tokenBlacklist {
		claims, err := ValidateTokenWithoutBlacklistCheck(token)
		if err != nil {
			// Token is invalid/expired, remove from blacklist
			delete(tokenBlacklist, token)
		} else if claims != nil {
			// Check if token is expired
			if exp, ok := (*claims)["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					delete(tokenBlacklist, token)
				}
			}
		}
	}
}
