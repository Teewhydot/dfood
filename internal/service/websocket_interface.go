package service

import (
	"dfood/internal/models"

	"github.com/gorilla/websocket"
)

// WebSocketService is the main interface for managing WebSocket connections and broadcasting messages
// This service acts as a central hub for all real-time communication in the application
type WebSocketService interface {
	// Connection Management - handles the lifecycle of WebSocket connections

	// RegisterConnection adds a new WebSocket connection to our registry
	// Returns a unique connectionID that can be used to reference this connection
	RegisterConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn) string

	// UnregisterConnection removes a connection and cleans up resources
	// Called when a client disconnects or connection fails
	UnregisterConnection(connectionID string)

	// GetConnection retrieves a specific connection by its ID
	GetConnection(connectionID string) *models.WSConnection

	// GetConnectionsByUser finds all connections belonging to a specific user
	// Useful when a user has multiple devices/tabs connected
	GetConnectionsByUser(userID string) []*models.WSConnection

	// GetConnectionsByType finds all connections of a specific type (e.g., all "user" connections)
	GetConnectionsByType(connectionType models.WSConnectionType) []*models.WSConnection

	// GetConnectionsByResource finds connections watching a specific resource (e.g., specific order)
	GetConnectionsByResource(connectionType models.WSConnectionType, resourceID string) []*models.WSConnection

	// Message Broadcasting - sends messages to different groups of connections

	// BroadcastToConnection sends a message to one specific connection
	BroadcastToConnection(connectionID string, message models.WSMessage)

	// BroadcastToUser sends a message to all of a user's connections
	// Example: User updates profile on phone, tablet gets notified
	BroadcastToUser(userID string, message models.WSMessage)

	// BroadcastToType sends a message to all connections of a specific type
	// Example: Broadcast to all users watching notifications
	BroadcastToType(connectionType models.WSConnectionType, message models.WSMessage)

	// BroadcastToResource sends a message to all connections watching a specific resource
	// Example: Order status changes, notify all watchers of that order
	BroadcastToResource(connectionType models.WSConnectionType, resourceID string, message models.WSMessage)

	// BroadcastToAll sends a message to every connected client
	// Use sparingly - only for system-wide announcements
	BroadcastToAll(message models.WSMessage)

	// Specific Broadcast Methods - convenience methods for common use cases
	// These handle the message formatting and routing automatically

	// BroadcastUserUpdate notifies when a user's profile changes
	BroadcastUserUpdate(userID string, user *models.User, changes map[string]interface{})

	// BroadcastOrderUpdate notifies when an order status or details change
	BroadcastOrderUpdate(orderID string, order *models.Order, status models.OrderStatus)

	// BroadcastNotification notifies when notifications are created/read/deleted
	BroadcastNotification(userID string, notification *models.Notification, action string)

	// BroadcastRestaurantUpdate notifies when restaurant info changes
	BroadcastRestaurantUpdate(restaurantID string, restaurant *models.Restaurant, changes map[string]interface{})

	// BroadcastAddressUpdate notifies when user addresses change
	BroadcastAddressUpdate(userID string, address *models.Address, action string, changes map[string]interface{})

	// BroadcastFavoriteUpdate notifies when favorites are added/removed
	BroadcastFavoriteUpdate(userID, resourceType, resourceID, action string, resource interface{})

	// Connection Handling - manages the WebSocket connection lifecycle

	// HandleConnection manages a WebSocket connection from start to finish
	// This includes ping/pong heartbeats, message reading, and cleanup
	HandleConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn)

	// Health Monitoring - provides insights into connection status

	// GetConnectionCount returns the total number of active connections
	GetConnectionCount() int

	// GetConnectionStats returns detailed statistics about connections
	// Useful for monitoring and debugging
	GetConnectionStats() map[string]interface{}
}

// Specialized WebSocket Service Interfaces
// These are optional interfaces that can be implemented for specific use cases
// The main WebSocketService interface above handles most scenarios

// UserWebSocketService provides user-specific WebSocket operations
type UserWebSocketService interface {
	WatchUser(userID string, conn *websocket.Conn)                                     // Start watching a user's profile changes
	NotifyUserUpdate(userID string, user *models.User, changes map[string]interface{}) // Send user update notification
}

// OrderWebSocketService provides order-specific WebSocket operations
type OrderWebSocketService interface {
	WatchOrder(orderID string, userID string, conn *websocket.Conn)                   // Start watching a specific order
	NotifyOrderUpdate(orderID string, order *models.Order, status models.OrderStatus) // Send order update notification
}

// NotificationWebSocketService provides notification-specific WebSocket operations
type NotificationWebSocketService interface {
	WatchNotifications(userID string, conn *websocket.Conn)                 // Start watching user's notifications
	NotifyNewNotification(userID string, notification *models.Notification) // Send new notification alert
}

// RestaurantWebSocketService provides restaurant-specific WebSocket operations
type RestaurantWebSocketService interface {
	WatchRestaurant(restaurantID string, conn *websocket.Conn)                                                 // Start watching restaurant changes
	NotifyRestaurantUpdate(restaurantID string, restaurant *models.Restaurant, changes map[string]interface{}) // Send restaurant update
}

// AddressWebSocketService provides address-specific WebSocket operations
type AddressWebSocketService interface {
	WatchAddresses(userID string, conn *websocket.Conn)                                                        // Start watching user's addresses
	NotifyAddressUpdate(userID string, address *models.Address, action string, changes map[string]interface{}) // Send address update
}

// FavoritesWebSocketService provides favorites-specific WebSocket operations
type FavoritesWebSocketService interface {
	WatchFavorites(userID string, resourceType string, conn *websocket.Conn)                    // Start watching user's favorites
	NotifyFavoriteUpdate(userID, resourceType, resourceID, action string, resource interface{}) // Send favorite update
}
