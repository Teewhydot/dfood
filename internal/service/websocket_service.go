package service

import (
	"dfood/internal/models"
	"dfood/internal/repository"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// webSocketService implements the WebSocketService interface
// It manages all WebSocket connections and handles message broadcasting
type webSocketService struct {
	// connections stores all active WebSocket connections in memory
	// Key: connectionID, Value: connection details
	connections map[string]*models.WSConnection

	// mutex protects the connections map from concurrent access
	// Multiple goroutines can register/unregister connections simultaneously
	mutex sync.RWMutex

	// Repository dependencies - used to fetch data when broadcasting updates
	userRepo         repository.UserRepository
	orderRepo        repository.OrderRepository
	notificationRepo repository.NotificationRepository
	restaurantRepo   repository.RestaurantRepository
	addressRepo      repository.AddressRepository
	favoritesRepo    repository.FavoritesRepository
}

// NewWebSocketService creates a new WebSocket service instance
// It initializes the connection registry and sets up repository dependencies
func NewWebSocketService(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
	notificationRepo repository.NotificationRepository,
	restaurantRepo repository.RestaurantRepository,
	addressRepo repository.AddressRepository,
	favoritesRepo repository.FavoritesRepository,
) WebSocketService {
	return &webSocketService{
		connections:      make(map[string]*models.WSConnection), // Initialize empty connection registry
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		notificationRepo: notificationRepo,
		restaurantRepo:   restaurantRepo,
		addressRepo:      addressRepo,
		favoritesRepo:    favoritesRepo,
	}
}

// === Connection Management Methods ===

// RegisterConnection adds a new WebSocket connection to our registry
// This is called when a client successfully connects to a WebSocket endpoint
func (s *webSocketService) RegisterConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn) string {
	// Lock the mutex to prevent concurrent access to the connections map
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Generate a unique ID for this connection
	connectionID := uuid.New().String()

	// Create a connection record with all the metadata we need
	wsConn := &models.WSConnection{
		ID:          connectionID,                 // Unique identifier
		Type:        connectionType,               // What kind of data this connection wants (user, order, etc.)
		UserID:      userID,                       // Which user owns this connection
		ResourceID:  resourceID,                   // Specific resource being watched (optional)
		Connection:  conn,                         // The actual WebSocket connection
		LastPing:    time.Now(),                   // Track connection health
		ConnectedAt: time.Now(),                   // When connection was established
		Metadata:    make(map[string]interface{}), // Additional data if needed
	}

	// Store the connection in our registry
	s.connections[connectionID] = wsConn

	log.Printf("WebSocket connection registered: %s (type: %s, user: %s, resource: %s)",
		connectionID, connectionType, userID, resourceID)

	// Send a confirmation message to the client so they know connection succeeded
	msg := models.WSMessage{
		Type:         models.WSMessageTypeConnected,
		Timestamp:    time.Now(),
		ConnectionID: connectionID,
		UserID:       userID,
	}
	s.sendMessage(conn, msg)

	return connectionID
}

func (s *webSocketService) UnregisterConnection(connectionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if wsConn, exists := s.connections[connectionID]; exists {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			conn.Close()
		}
		delete(s.connections, connectionID)
		log.Printf("WebSocket connection unregistered: %s", connectionID)
	}
}

func (s *webSocketService) GetConnection(connectionID string) *models.WSConnection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.connections[connectionID]
}

func (s *webSocketService) GetConnectionsByUser(userID string) []*models.WSConnection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var connections []*models.WSConnection
	for _, conn := range s.connections {
		if conn.UserID == userID {
			connections = append(connections, conn)
		}
	}
	return connections
}

func (s *webSocketService) GetConnectionsByType(connectionType models.WSConnectionType) []*models.WSConnection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var connections []*models.WSConnection
	for _, conn := range s.connections {
		if conn.Type == connectionType {
			connections = append(connections, conn)
		}
	}
	return connections
}

func (s *webSocketService) GetConnectionsByResource(connectionType models.WSConnectionType, resourceID string) []*models.WSConnection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var connections []*models.WSConnection
	for _, conn := range s.connections {
		if conn.Type == connectionType && conn.ResourceID == resourceID {
			connections = append(connections, conn)
		}
	}
	return connections
}

// === Message Broadcasting Methods ===
// These methods send messages to different groups of connected clients

// BroadcastToConnection sends a message to one specific connection
// Used when you need to send a message to a particular client
func (s *webSocketService) BroadcastToConnection(connectionID string, message models.WSMessage) {
	// Get the connection (read lock is sufficient since we're not modifying)
	s.mutex.RLock()
	wsConn := s.connections[connectionID]
	s.mutex.RUnlock()

	// Send message if connection exists and is valid
	if wsConn != nil {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

// BroadcastToUser sends a message to all connections belonging to a specific user
// Example: User updates profile on phone, notify their tablet too
func (s *webSocketService) BroadcastToUser(userID string, message models.WSMessage) {
	connections := s.GetConnectionsByUser(userID)
	for _, wsConn := range connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

// BroadcastToType sends a message to all connections of a specific type
// Example: Send notification to all users watching notifications
func (s *webSocketService) BroadcastToType(connectionType models.WSConnectionType, message models.WSMessage) {
	connections := s.GetConnectionsByType(connectionType)
	for _, wsConn := range connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

// BroadcastToResource sends a message to all connections watching a specific resource
// Example: Order status changes, notify everyone watching that order
func (s *webSocketService) BroadcastToResource(connectionType models.WSConnectionType, resourceID string, message models.WSMessage) {
	connections := s.GetConnectionsByResource(connectionType, resourceID)
	for _, wsConn := range connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

// BroadcastToAll sends a message to every connected client
// Use sparingly - only for system-wide announcements
func (s *webSocketService) BroadcastToAll(message models.WSMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, wsConn := range s.connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

// === Specific Broadcast Methods ===
// These are convenience methods that handle common broadcasting scenarios
// They automatically format the message and route it to the right connections

// BroadcastUserUpdate sends a user profile update to all connections watching that user
// Called automatically when user data changes via the API
func (s *webSocketService) BroadcastUserUpdate(userID string, user *models.User, changes map[string]interface{}) {
	// Create the payload with user data and what changed
	payload := models.UserUpdatePayload{
		User:      user,    // Complete updated user object
		Changes:   changes, // Only the fields that changed
		UpdatedBy: userID,  // Who made the change
	}

	// Convert payload to JSON for transmission
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling user update: %v", err)
		return
	}

	// Create the WebSocket message
	msg := models.WSMessage{
		Type:      models.WSMessageTypeUserUpdate,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	// Send to all connections watching this user's profile
	s.BroadcastToResource(models.WSConnectionTypeUser, userID, msg)
}

// BroadcastOrderUpdate sends order status changes to relevant connections
// This enables real-time order tracking for customers
func (s *webSocketService) BroadcastOrderUpdate(orderID string, order *models.Order, status models.OrderStatus) {
	// Create payload with order details and new status
	payload := models.OrderUpdatePayload{
		Order:     order,    // Complete order object
		Status:    status,   // New status (pending, confirmed, delivered, etc.)
		UpdatedBy: "system", // Usually updated by system or restaurant
	}

	// Convert to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling order update: %v", err)
		return
	}

	// Create WebSocket message
	msg := models.WSMessage{
		Type:      models.WSMessageTypeOrderUpdate,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Send to anyone specifically watching this order
	s.BroadcastToResource(models.WSConnectionTypeOrder, orderID, msg)

	// Also send to the user who owns the order (they might be watching their orders in general)
	if order.UserID != "" {
		s.BroadcastToUser(order.UserID, msg)
	}
}

func (s *webSocketService) BroadcastNotification(userID string, notification *models.Notification, action string) {
	payload := models.NotificationPayload{
		Notification: notification,
		Action:       action,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	var msgType models.WSMessageType
	switch action {
	case "created":
		msgType = models.WSMessageTypeNotificationNew
	case "read":
		msgType = models.WSMessageTypeNotificationRead
	default:
		msgType = models.WSMessageTypeNotification
	}

	msg := models.WSMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	// Broadcast to notification connections for this user
	s.BroadcastToResource(models.WSConnectionTypeNotification, userID, msg)
}

func (s *webSocketService) BroadcastRestaurantUpdate(restaurantID string, restaurant *models.Restaurant, changes map[string]interface{}) {
	payload := models.RestaurantUpdatePayload{
		Restaurant: restaurant,
		Changes:    changes,
		UpdatedBy:  "system",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling restaurant update: %v", err)
		return
	}

	msg := models.WSMessage{
		Type:      models.WSMessageTypeRestaurantUpdate,
		Data:      data,
		Timestamp: time.Now(),
	}

	s.BroadcastToResource(models.WSConnectionTypeRestaurant, restaurantID, msg)
}

func (s *webSocketService) BroadcastAddressUpdate(userID string, address *models.Address, action string, changes map[string]interface{}) {
	payload := models.AddressUpdatePayload{
		Address:   address,
		Action:    action,
		Changes:   changes,
		UpdatedBy: userID,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling address update: %v", err)
		return
	}

	var msgType models.WSMessageType
	switch action {
	case "created":
		msgType = models.WSMessageTypeAddressNew
	case "deleted":
		msgType = models.WSMessageTypeAddressDelete
	default:
		msgType = models.WSMessageTypeAddressUpdate
	}

	msg := models.WSMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	// Broadcast to address connections for this user
	s.BroadcastToResource(models.WSConnectionTypeAddress, userID, msg)
}

func (s *webSocketService) BroadcastFavoriteUpdate(userID, resourceType, resourceID, action string, resource interface{}) {
	payload := models.FavoriteUpdatePayload{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Action:       action,
		Resource:     resource,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling favorite update: %v", err)
		return
	}

	var msgType models.WSMessageType
	switch action {
	case "add":
		msgType = models.WSMessageTypeFavoriteAdd
	case "remove":
		msgType = models.WSMessageTypeFavoriteRemove
	case "clear":
		msgType = models.WSMessageTypeFavoritesClear
	default:
		msgType = models.WSMessageTypeFavoriteAdd
	}

	msg := models.WSMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	// Broadcast to favorites connections for this user and resource type
	resourceKey := userID + "_" + resourceType
	s.BroadcastToResource(models.WSConnectionTypeFavorites, resourceKey, msg)
}

// Connection handling
func (s *webSocketService) HandleConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn) {
	connectionID := s.RegisterConnection(connectionType, userID, resourceID, conn)
	defer s.UnregisterConnection(connectionID)

	// Set up ping/pong handlers
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if wsConn := s.GetConnection(connectionID); wsConn != nil {
			wsConn.LastPing = time.Now()
		}
		return nil
	})

	// Start ping ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	// Read messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for connection %s: %v", connectionID, err)
			}
			break
		}
	}
}

// Health check
func (s *webSocketService) GetConnectionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.connections)
}

func (s *webSocketService) GetConnectionStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := make(map[string]interface{})
	typeCount := make(map[models.WSConnectionType]int)

	for _, conn := range s.connections {
		typeCount[conn.Type]++
	}

	stats["total_connections"] = len(s.connections)
	stats["connections_by_type"] = typeCount
	stats["timestamp"] = time.Now()

	return stats
}

func (s *webSocketService) sendMessage(conn *websocket.Conn, msg models.WSMessage) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending WebSocket message: %v", err)
	}
}
