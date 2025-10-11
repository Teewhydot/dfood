package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateID generates a unique ID for entities
func GenerateID() string {
	// Generate a simple unique ID using timestamp and random bytes
	timestamp := time.Now().UnixNano()

	// Generate 4 random bytes
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	// Combine timestamp and random bytes to create a unique ID
	return fmt.Sprintf("%d-%x", timestamp, randomBytes)
}

// GenerateUserID generates a user-specific ID
func GenerateUserID() string {
	return "user-" + GenerateID()
}

// GenerateOrderID generates an order-specific ID
func GenerateOrderID() string {
	return "order-" + GenerateID()
}

// GenerateRestaurantID generates a restaurant-specific ID
func GenerateRestaurantID() string {
	return "restaurant-" + GenerateID()
}

// GenerateFoodID generates a food-specific ID
func GenerateFoodID() string {
	return "food-" + GenerateID()
}

// GenerateRandomPassword generates a random password with specified length
func GenerateRandomPassword(length int) string {
	if length < 8 {
		length = 8 // Minimum password length
	}

	// Character sets for password generation
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	special := "!@#$%^&*"

	// Combine all character sets
	allChars := lowercase + uppercase + numbers + special

	// Generate random password
	password := make([]byte, length)
	for i := range password {
		randomIndex := make([]byte, 1)
		rand.Read(randomIndex)
		password[i] = allChars[int(randomIndex[0])%len(allChars)]
	}

	return string(password)
}
