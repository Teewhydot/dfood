package models

import (
	"time"

	"github.com/gorilla/websocket"
)

// Message types - just what we need
const (
	MessageTypeUserUpdate = "user_update"
	MessageTypeConnected  = "connected"
	MessageTypeError      = "error"
)

// WSMessage - like a DTO in Flutter
// Every WebSocket message follows this simple structure
type WSMessage struct {
	Type      string      `json:"type"`      // What kind of message (user_update, connected, error)
	Data      interface{} `json:"data"`      // The actual content
	Timestamp time.Time   `json:"timestamp"` // When it was sent
}

// WSConnection - like a model in Flutter
// Represents one user's WebSocket connection
type WSConnection struct {
	UserID     string          // Which user this connection belongs to
	Connection *websocket.Conn // The actual WebSocket connection
}

// UserUpdateData - specific data for user updates
// This goes in the Data field of WSMessage
type UserUpdateData struct {
	User    *User                  `json:"user"`    // Complete updated user object
	Changes map[string]interface{} `json:"changes"` // Only the fields that changed
}
