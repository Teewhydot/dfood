package models

import (
	"time"

	"github.com/gorilla/websocket"
)

// Message types - just what we need
const (
	// User messages
	MessageTypeUserUpdate = "user_update"
	MessageTypeConnected  = "connected"
	MessageTypeError      = "error"

	// Address messages
	MessageTypeAddressAdd     = "address_add"
	MessageTypeAddressUpdate  = "address_update"
	MessageTypeAddressDelete  = "address_delete"
	MessageTypeAddressDefault = "address_default"

	// Favorites messages
	MessageTypeFavoriteAdd    = "favorite_add"
	MessageTypeFavoriteRemove = "favorite_remove"
	MessageTypeFavoritesClear = "favorites_clear"

	// Notification messages
	MessageTypeNotificationNew    = "notification_new"
	MessageTypeNotificationRead   = "notification_read"
	MessageTypeNotificationDelete = "notification_delete"
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

// AddressUpdateData - for address changes
type AddressUpdateData struct {
	Address *Address               `json:"address"`           // Address object
	Action  string                 `json:"action"`            // "add", "update", "delete", "set_default"
	Changes map[string]interface{} `json:"changes,omitempty"` // What changed (for updates)
}

// FavoriteUpdateData - for favorites changes
type FavoriteUpdateData struct {
	ResourceType string      `json:"resource_type"` // "food" or "restaurant"
	ResourceID   string      `json:"resource_id"`   // ID of the item
	Action       string      `json:"action"`        // "add", "remove", "clear"
	Resource     interface{} `json:"resource"`      // The actual Food or Restaurant object
}

// NotificationUpdateData - for notification changes
type NotificationUpdateData struct {
	Notification *Notification `json:"notification"` // Notification object
	Action       string        `json:"action"`       // "new", "read", "delete"
}
