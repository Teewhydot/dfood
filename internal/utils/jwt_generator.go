package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

// RotateSecretKey changes the JWT secret, invalidating all existing tokens
func RotateSecretKey(newSecret string) {
	jwtKey = []byte(newSecret)
	// Clear blacklist since all tokens are now invalid anyway
	tokenBlacklist = make(map[string]bool)
}

func GenerateJwtToken(email string, isRefresh bool) (string, error) {
	var expirationTime time.Time
	if isRefresh {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	} else {
		expirationTime = time.Now().Add(15 * time.Minute)
	}
	claims := &jwt.MapClaims{
		"sub": email,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	// Check if token is blacklisted first
	if IsTokenBlacklisted(tokenStr) {
		return nil, jwt.ErrTokenInvalidClaims
	}

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
