package models

import (
	"encoding/json"
	"time"
)

// WSConnectionType categorizes different types of WebSocket connections
// This helps us organize and route messages to the right subscribers
type WSConnectionType string

const (
	// Each connection type represents a different data stream
	WSConnectionTypeUser         WSConnectionType = "user"         // User profile updates
	WSConnectionTypeOrder        WSConnectionType = "order"        // Order status changes
	WSConnectionTypeNotification WSConnectionType = "notification" // New notifications
	WSConnectionTypeRestaurant   WSConnectionType = "restaurant"   // Restaurant info updates
	WSConnectionTypeAddress      WSConnectionType = "address"      // Address changes
	WSConnectionTypeFavorites    WSConnectionType = "favorites"    // Favorites add/remove
)

// WSMessageType defines what kind of message is being sent over WebSocket
// Clients use this to determine how to handle the incoming message
type WSMessageType string

const (
	// User related messages - sent when user profile changes
	WSMessageTypeUserUpdate    WSMessageType = "user_update"    // Profile field updated
	WSMessageTypeUserConnected WSMessageType = "user_connected" // Successfully connected

	// Order related messages - sent during order lifecycle
	WSMessageTypeOrderUpdate WSMessageType = "order_update" // Order details changed
	WSMessageTypeOrderStatus WSMessageType = "order_status" // Order status changed (pending → confirmed → delivered)

	// Notification related messages - sent for user notifications
	WSMessageTypeNotification     WSMessageType = "notification"      // General notification event
	WSMessageTypeNotificationRead WSMessageType = "notification_read" // User marked notification as read
	WSMessageTypeNotificationNew  WSMessageType = "notification_new"  // New notification created

	// Restaurant related messages - sent when restaurant data changes
	WSMessageTypeRestaurantUpdate WSMessageType = "restaurant_update" // Restaurant info updated
	WSMessageTypeMenuUpdate       WSMessageType = "menu_update"       // Menu items changed

	// Address related messages - sent when user addresses change
	WSMessageTypeAddressUpdate WSMessageType = "address_update" // Existing address modified
	WSMessageTypeAddressNew    WSMessageType = "address_new"    // New address added
	WSMessageTypeAddressDelete WSMessageType = "address_delete" // Address removed

	// Favorites related messages - sent when user favorites change
	WSMessageTypeFavoriteAdd    WSMessageType = "favorite_add"    // Item added to favorites
	WSMessageTypeFavoriteRemove WSMessageType = "favorite_remove" // Item removed from favorites
	WSMessageTypeFavoritesClear WSMessageType = "favorites_clear" // All favorites cleared

	// System messages - for connection management and health checks
	WSMessageTypeError        WSMessageType = "error"        // Error occurred
	WSMessageTypePing         WSMessageType = "ping"         // Server → Client heartbeat
	WSMessageTypePong         WSMessageType = "pong"         // Client → Server heartbeat response
	WSMessageTypeConnected    WSMessageType = "connected"    // Connection established successfully
	WSMessageTypeDisconnected WSMessageType = "disconnected" // Connection closed
)

// WSMessage is the standard format for all WebSocket messages
// Every message sent over WebSocket follows this structure
type WSMessage struct {
	Type         WSMessageType   `json:"type"`                    // What kind of message this is
	Data         json.RawMessage `json:"data,omitempty"`          // The actual payload (varies by type)
	Timestamp    time.Time       `json:"timestamp"`               // When this message was created
	ConnectionID string          `json:"connection_id,omitempty"` // Which connection this relates to
	UserID       string          `json:"user_id,omitempty"`       // Which user this message is for
}

// WSConnection represents an active WebSocket connection
// We store metadata about each connection to manage them properly
type WSConnection struct {
	ID          string                 // Unique identifier for this connection
	Type        WSConnectionType       // What kind of data this connection subscribes to
	UserID      string                 // Which user owns this connection
	ResourceID  string                 // Specific resource being watched (orderID, restaurantID, etc.)
	Connection  any                    // The actual WebSocket connection (*websocket.Conn)
	LastPing    time.Time              // Last time we received a ping (for health checking)
	ConnectedAt time.Time              // When this connection was established
	Metadata    map[string]interface{} // Additional connection-specific data
}

// Payload types - these go inside WSMessage.Data for specific message types

// UserUpdatePayload is sent when a user's profile changes
// Contains the updated user object and what specifically changed
type UserUpdatePayload struct {
	User      *User                  `json:"user"`       // Complete updated user object
	Changes   map[string]interface{} `json:"changes"`    // Only the fields that changed (e.g., {"first_name": "John"})
	UpdatedBy string                 `json:"updated_by"` // Who made the change (usually the user themselves)
}

// OrderUpdatePayload is sent when an order status or details change
// Useful for real-time order tracking
type OrderUpdatePayload struct {
	Order     *Order      `json:"order"`      // Complete updated order object
	Status    OrderStatus `json:"status"`     // New order status
	UpdatedBy string      `json:"updated_by"` // Who updated it (user, restaurant, system)
}

// NotificationPayload is sent when notifications are created, read, or deleted
// Allows real-time notification updates without polling
type NotificationPayload struct {
	Notification *Notification `json:"notification"` // The notification object
	Action       string        `json:"action"`       // What happened: "created", "read", "deleted"
}

// RestaurantUpdatePayload is sent when restaurant information changes
// Useful for live menu updates or restaurant status changes
type RestaurantUpdatePayload struct {
	Restaurant *Restaurant            `json:"restaurant"` // Complete updated restaurant object
	Changes    map[string]interface{} `json:"changes"`    // What changed (e.g., {"is_open": true})
	UpdatedBy  string                 `json:"updated_by"` // Who made the change
}

// AddressUpdatePayload is sent when user addresses are added, updated, or deleted
// Keeps address lists in sync across devices
type AddressUpdatePayload struct {
	Address   *Address               `json:"address"`           // The address object
	Action    string                 `json:"action"`            // "created", "updated", "deleted", "set_default"
	Changes   map[string]interface{} `json:"changes,omitempty"` // What changed (for updates)
	UpdatedBy string                 `json:"updated_by"`        // Who made the change
}

// FavoriteUpdatePayload is sent when favorites are added or removed
// Keeps favorite lists synchronized in real-time
type FavoriteUpdatePayload struct {
	UserID       string `json:"user_id"`            // Which user's favorites changed
	ResourceType string `json:"resource_type"`      // "food" or "restaurant"
	ResourceID   string `json:"resource_id"`        // ID of the favorited item
	Action       string `json:"action"`             // "add", "remove", "clear"
	Resource     any    `json:"resource,omitempty"` // The actual Food or Restaurant object (for convenience)
}
