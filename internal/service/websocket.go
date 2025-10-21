package service

import (
	"dfood/internal/models"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketService - like a service in Flutter
// Manages WebSocket connections and sends messages to users
type WebSocketService struct {
	// connections maps userID -> WebSocket connection
	// We keep it simple: one connection per user
	connections map[string]*websocket.Conn

	// mutex protects the connections map from concurrent access
	// Multiple users can connect/disconnect at the same time
	mutex sync.RWMutex
}

// NewWebSocketService creates a new WebSocket service
// Like creating a service in Flutter's dependency injection
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		connections: make(map[string]*websocket.Conn),
	}
}

// AddConnection registers a new user connection
// If user already has a connection, we replace it (single device per user)
func (s *WebSocketService) AddConnection(userID string, conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Close existing connection if user reconnects
	if existingConn, exists := s.connections[userID]; exists {
		existingConn.Close()
		log.Printf("WebSocket: Replaced existing connection for user %s", userID)
	}

	// Store the new connection
	s.connections[userID] = conn
	log.Printf("WebSocket: User %s connected", userID)

	// Send connection confirmation to the user
	s.sendToUser(userID, models.WSMessage{
		Type:      models.MessageTypeConnected,
		Data:      map[string]string{"message": "Connected successfully"},
		Timestamp: time.Now(),
	})
}

// RemoveConnection removes a user's connection
// Called when user disconnects or connection fails
func (s *WebSocketService) RemoveConnection(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if conn, exists := s.connections[userID]; exists {
		conn.Close()
		delete(s.connections, userID)
		log.Printf("WebSocket: User %s disconnected", userID)
	}
}

// SendUserUpdate sends a user profile update to the user
// This is the main method you'll call from your business logic
func (s *WebSocketService) SendUserUpdate(userID string, user *models.User, changes map[string]interface{}) {
	// Create the update data
	updateData := models.UserUpdateData{
		User:    user,
		Changes: changes,
	}

	// Create the WebSocket message
	message := models.WSMessage{
		Type:      models.MessageTypeUserUpdate,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	// Send it to the user
	s.sendToUser(userID, message)
	log.Printf("WebSocket: Sent user update to %s, changes: %v", userID, changes)
}

// GetConnectionCount returns how many users are connected
// Useful for monitoring and debugging
func (s *WebSocketService) GetConnectionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.connections)
}

// sendToUser sends a message to a specific user (private method)
// Handles errors and cleans up bad connections automatically
func (s *WebSocketService) sendToUser(userID string, message models.WSMessage) {
	s.mutex.RLock()
	conn, exists := s.connections[userID]
	s.mutex.RUnlock()

	if !exists {
		// User not connected, that's okay - they might be offline
		return
	}

	// Try to send the message
	if err := conn.WriteJSON(message); err != nil {
		log.Printf("WebSocket: Error sending to user %s: %v", userID, err)
		// Connection is bad, clean it up
		s.RemoveConnection(userID)
	}
}

// HandleConnection manages a WebSocket connection lifecycle
// This method blocks until the connection closes
func (s *WebSocketService) HandleConnection(userID string, conn *websocket.Conn) {
	// Register the connection
	s.AddConnection(userID, conn)

	// Ensure cleanup when this function exits
	defer s.RemoveConnection(userID)

	// Keep the connection alive by reading messages
	// We don't need to do anything with the messages, just read them
	// to detect when the connection closes
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Connection closed or error occurred
			log.Printf("WebSocket: Connection closed for user %s: %v", userID, err)
			break
		}
		// If we wanted to handle incoming messages from client, we'd do it here
		// For now, we just keep the connection alive
	}
}
