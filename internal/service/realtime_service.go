package service

import (
	"dfood/internal/models"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RealtimeService manages all WebSocket connections and real-time updates
type RealtimeService struct {
	// connections maps userID -> WebSocket connection
	connections map[string]*websocket.Conn
	mutex       sync.RWMutex
}

// NewRealtimeService creates a new realtime service
func NewRealtimeService() *RealtimeService {
	return &RealtimeService{
		connections: make(map[string]*websocket.Conn),
	}
}

// Connection Management

// AddConnection registers a new user connection
func (s *RealtimeService) AddConnection(userID string, conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Close existing connection if user reconnects
	if existingConn, exists := s.connections[userID]; exists {
		existingConn.Close()
		log.Printf("Realtime: Replaced existing connection for user %s", userID)
	}

	// Store the new connection
	s.connections[userID] = conn
	log.Printf("Realtime: User %s connected", userID)

	// Send connection confirmation
	s.sendToUser(userID, models.WSMessage{
		Type:      models.MessageTypeConnected,
		Data:      map[string]string{"message": "Connected successfully"},
		Timestamp: time.Now(),
	})
}

// RemoveConnection removes a user's connection
func (s *RealtimeService) RemoveConnection(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if conn, exists := s.connections[userID]; exists {
		conn.Close()
		delete(s.connections, userID)
		log.Printf("Realtime: User %s disconnected", userID)
	}
}

// GetConnectionCount returns number of connected users
func (s *RealtimeService) GetConnectionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.connections)
}

// HandleConnection manages a WebSocket connection lifecycle
func (s *RealtimeService) HandleConnection(userID string, conn *websocket.Conn) {
	s.AddConnection(userID, conn)
	defer s.RemoveConnection(userID)

	// Keep connection alive by reading messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Realtime: Connection closed for user %s: %v", userID, err)
			break
		}
	}
}

// User Profile Updates

// SendUserUpdate broadcasts user profile changes
func (s *RealtimeService) SendUserUpdate(userID string, user *models.User, changes map[string]interface{}) {
	updateData := models.UserUpdateData{
		User:    user,
		Changes: changes,
	}

	message := models.WSMessage{
		Type:      models.MessageTypeUserUpdate,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent user update to %s, changes: %v", userID, changes)
}

// Address Updates

// SendAddressAdd broadcasts when a new address is added
func (s *RealtimeService) SendAddressAdd(userID string, address *models.Address) {
	updateData := models.AddressUpdateData{
		Address: address,
		Action:  "add",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeAddressAdd,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent address add to %s, address: %s", userID, address.ID)
}

// SendAddressUpdate broadcasts when an address is updated
func (s *RealtimeService) SendAddressUpdate(userID string, address *models.Address, changes map[string]interface{}) {
	updateData := models.AddressUpdateData{
		Address: address,
		Action:  "update",
		Changes: changes,
	}

	message := models.WSMessage{
		Type:      models.MessageTypeAddressUpdate,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent address update to %s, changes: %v", userID, changes)
}

// SendAddressDelete broadcasts when an address is deleted
func (s *RealtimeService) SendAddressDelete(userID string, address *models.Address) {
	updateData := models.AddressUpdateData{
		Address: address,
		Action:  "delete",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeAddressDelete,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent address delete to %s, address: %s", userID, address.ID)
}

// SendAddressDefault broadcasts when default address changes
func (s *RealtimeService) SendAddressDefault(userID string, address *models.Address) {
	updateData := models.AddressUpdateData{
		Address: address,
		Action:  "set_default",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeAddressDefault,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent address default to %s, address: %s", userID, address.ID)
}

// Favorites Updates

// SendFavoriteAdd broadcasts when item is added to favorites
func (s *RealtimeService) SendFavoriteAdd(userID, resourceType, resourceID string, resource interface{}) {
	updateData := models.FavoriteUpdateData{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Action:       "add",
		Resource:     resource,
	}

	message := models.WSMessage{
		Type:      models.MessageTypeFavoriteAdd,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent favorite add to %s, %s: %s", userID, resourceType, resourceID)
}

// SendFavoriteRemove broadcasts when item is removed from favorites
func (s *RealtimeService) SendFavoriteRemove(userID, resourceType, resourceID string, resource interface{}) {
	updateData := models.FavoriteUpdateData{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Action:       "remove",
		Resource:     resource,
	}

	message := models.WSMessage{
		Type:      models.MessageTypeFavoriteRemove,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent favorite remove to %s, %s: %s", userID, resourceType, resourceID)
}

// SendFavoritesClear broadcasts when all favorites are cleared
func (s *RealtimeService) SendFavoritesClear(userID string) {
	updateData := models.FavoriteUpdateData{
		Action: "clear",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeFavoritesClear,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent favorites clear to %s", userID)
}

// Notification Updates

// SendNotificationNew broadcasts when a new notification is created
func (s *RealtimeService) SendNotificationNew(userID string, notification *models.Notification) {
	updateData := models.NotificationUpdateData{
		Notification: notification,
		Action:       "new",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeNotificationNew,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent new notification to %s, notification: %s", userID, notification.ID)
}

// SendNotificationRead broadcasts when a notification is marked as read
func (s *RealtimeService) SendNotificationRead(userID string, notification *models.Notification) {
	updateData := models.NotificationUpdateData{
		Notification: notification,
		Action:       "read",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeNotificationRead,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent notification read to %s, notification: %s", userID, notification.ID)
}

// SendNotificationDelete broadcasts when a notification is deleted
func (s *RealtimeService) SendNotificationDelete(userID string, notification *models.Notification) {
	updateData := models.NotificationUpdateData{
		Notification: notification,
		Action:       "delete",
	}

	message := models.WSMessage{
		Type:      models.MessageTypeNotificationDelete,
		Data:      updateData,
		Timestamp: time.Now(),
	}

	s.sendToUser(userID, message)
	log.Printf("Realtime: Sent notification delete to %s, notification: %s", userID, notification.ID)
}

// Private helper method to send message to a specific user
func (s *RealtimeService) sendToUser(userID string, message models.WSMessage) {
	s.mutex.RLock()
	conn, exists := s.connections[userID]
	s.mutex.RUnlock()

	if !exists {
		// User not connected, that's okay
		return
	}

	// Try to send the message
	if err := conn.WriteJSON(message); err != nil {
		log.Printf("Realtime: Error sending to user %s: %v", userID, err)
		// Connection is bad, clean it up
		s.RemoveConnection(userID)
	}
}
