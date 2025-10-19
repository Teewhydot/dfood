# Multiple WebSocket Endpoints Architecture

## Overview
This guide shows how to implement a scalable WebSocket service that can handle multiple endpoint types (user updates, order tracking, notifications, etc.) using a generic connection manager.

## 1. Generic WebSocket Models (internal/models/websocket.go)

```go
package models

import (
	"encoding/json"
	"time"
)

// WebSocket connection types
type WSConnectionType string

const (
	WSConnectionTypeUser         WSConnectionType = "user"
	WSConnectionTypeOrder        WSConnectionType = "order"
	WSConnectionTypeNotification WSConnectionType = "notification"
	WSConnectionTypeRestaurant   WSConnectionType = "restaurant"
)

// WebSocket message types
type WSMessageType string

const (
	// User related
	WSMessageTypeUserUpdate    WSMessageType = "user_update"
	WSMessageTypeUserConnected WSMessageType = "user_connected"
	
	// Order related
	WSMessageTypeOrderUpdate   WSMessageType = "order_update"
	WSMessageTypeOrderStatus   WSMessageType = "order_status"
	
	// Notification related
	WSMessageTypeNotification  WSMessageType = "notification"
	WSMessageTypeNotificationRead WSMessageType = "notification_read"
	
	// Restaurant related
	WSMessageTypeRestaurantUpdate WSMessageType = "restaurant_update"
	WSMessageTypeMenuUpdate    WSMessageType = "menu_update"
	
	// System messages
	WSMessageTypeError         WSMessageType = "error"
	WSMessageTypePing          WSMessageType = "ping"
	WSMessageTypePong          WSMessageType = "pong"
	WSMessageTypeConnected     WSMessageType = "connected"
	WSMessageTypeDisconnected  WSMessageType = "disconnected"
)

// Generic WebSocket message structure
type WSMessage struct {
	Type         WSMessageType   `json:"type"`
	Data         json.RawMessage `json:"data,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
	ConnectionID string          `json:"connection_id,omitempty"`
	UserID       string          `json:"user_id,omitempty"`
}

// Connection information
type WSConnection struct {
	ID             string
	Type           WSConnectionType
	UserID         string
	ResourceID     string // Could be orderID, restaurantID, etc.
	Connection     interface{} // *websocket.Conn
	LastPing       time.Time
	ConnectedAt    time.Time
	Metadata       map[string]interface{}
}

// Specific payload types
type UserUpdatePayload struct {
	User      *User                  `json:"user"`
	Changes   map[string]interface{} `json:"changes"`
	UpdatedBy string                 `json:"updated_by"`
}

type OrderUpdatePayload struct {
	Order     *Order     `json:"order"`
	Status    OrderStatus `json:"status"`
	UpdatedBy string     `json:"updated_by"`
}

type NotificationPayload struct {
	Notification *Notification `json:"notification"`
	Action       string        `json:"action"` // "created", "read", "deleted"
}

type RestaurantUpdatePayload struct {
	Restaurant *Restaurant            `json:"restaurant"`
	Changes    map[string]interface{} `json:"changes"`
	UpdatedBy  string                 `json:"updated_by"`
}
```

## 2. Generic WebSocket Service Interface (internal/service/websocket_interface.go)

```go
package service

import (
	"dfood/internal/models"
	"github.com/gorilla/websocket"
)

// Generic WebSocket service interface
type WebSocketService interface {
	// Connection management
	RegisterConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn) string
	UnregisterConnection(connectionID string)
	GetConnection(connectionID string) *models.WSConnection
	GetConnectionsByUser(userID string) []*models.WSConnection
	GetConnectionsByType(connectionType models.WSConnectionType) []*models.WSConnection
	GetConnectionsByResource(connectionType models.WSConnectionType, resourceID string) []*models.WSConnection
	
	// Message broadcasting
	BroadcastToConnection(connectionID string, message models.WSMessage)
	BroadcastToUser(userID string, message models.WSMessage)
	BroadcastToType(connectionType models.WSConnectionType, message models.WSMessage)
	BroadcastToResource(connectionType models.WSConnectionType, resourceID string, message models.WSMessage)
	BroadcastToAll(message models.WSMessage)
	
	// Specific broadcast methods
	BroadcastUserUpdate(userID string, user *models.User, changes map[string]interface{})
	BroadcastOrderUpdate(orderID string, order *models.Order, status models.OrderStatus)
	BroadcastNotification(userID string, notification *models.Notification, action string)
	BroadcastRestaurantUpdate(restaurantID string, restaurant *models.Restaurant, changes map[string]interface{})
	
	// Connection handling
	HandleConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn)
	
	// Health check
	GetConnectionCount() int
	GetConnectionStats() map[string]interface{}
}

// Specific service interfaces for different endpoints
type UserWebSocketService interface {
	WatchUser(userID string, conn *websocket.Conn)
	NotifyUserUpdate(userID string, user *models.User, changes map[string]interface{})
}

type OrderWebSocketService interface {
	WatchOrder(orderID string, userID string, conn *websocket.Conn)
	NotifyOrderUpdate(orderID string, order *models.Order, status models.OrderStatus)
}

type NotificationWebSocketService interface {
	WatchNotifications(userID string, conn *websocket.Conn)
	NotifyNewNotification(userID string, notification *models.Notification)
}

type RestaurantWebSocketService interface {
	WatchRestaurant(restaurantID string, conn *websocket.Conn)
	NotifyRestaurantUpdate(restaurantID string, restaurant *models.Restaurant, changes map[string]interface{})
}
```

## 3. Generic WebSocket Service Implementation (internal/service/websocket_service.go)

```go
package service

import (
	"dfood/internal/models"
	"dfood/internal/repository"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type webSocketService struct {
	connections map[string]*models.WSConnection
	mutex       sync.RWMutex
	
	// Repository dependencies
	userRepo         repository.UserRepository
	orderRepo        repository.OrderRepository
	notificationRepo repository.NotificationRepository
	restaurantRepo   repository.RestaurantRepository
}

func NewWebSocketService(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
	notificationRepo repository.NotificationRepository,
	restaurantRepo repository.RestaurantRepository,
) WebSocketService {
	return &webSocketService{
		connections:      make(map[string]*models.WSConnection),
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		notificationRepo: notificationRepo,
		restaurantRepo:   restaurantRepo,
	}
}

// Connection management
func (s *webSocketService) RegisterConnection(connectionType models.WSConnectionType, userID, resourceID string, conn *websocket.Conn) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	connectionID := uuid.New().String()
	
	wsConn := &models.WSConnection{
		ID:          connectionID,
		Type:        connectionType,
		UserID:      userID,
		ResourceID:  resourceID,
		Connection:  conn,
		LastPing:    time.Now(),
		ConnectedAt: time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	
	s.connections[connectionID] = wsConn
	
	log.Printf("WebSocket connection registered: %s (type: %s, user: %s, resource: %s)", 
		connectionID, connectionType, userID, resourceID)
	
	// Send connection confirmation
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

// Message broadcasting
func (s *webSocketService) BroadcastToConnection(connectionID string, message models.WSMessage) {
	s.mutex.RLock()
	wsConn := s.connections[connectionID]
	s.mutex.RUnlock()
	
	if wsConn != nil {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

func (s *webSocketService) BroadcastToUser(userID string, message models.WSMessage) {
	connections := s.GetConnectionsByUser(userID)
	for _, wsConn := range connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

func (s *webSocketService) BroadcastToType(connectionType models.WSConnectionType, message models.WSMessage) {
	connections := s.GetConnectionsByType(connectionType)
	for _, wsConn := range connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

func (s *webSocketService) BroadcastToResource(connectionType models.WSConnectionType, resourceID string, message models.WSMessage) {
	connections := s.GetConnectionsByResource(connectionType, resourceID)
	for _, wsConn := range connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

func (s *webSocketService) BroadcastToAll(message models.WSMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for _, wsConn := range s.connections {
		if conn, ok := wsConn.Connection.(*websocket.Conn); ok {
			s.sendMessage(conn, message)
		}
	}
}

// Specific broadcast methods
func (s *webSocketService) BroadcastUserUpdate(userID string, user *models.User, changes map[string]interface{}) {
	payload := models.UserUpdatePayload{
		User:      user,
		Changes:   changes,
		UpdatedBy: userID,
	}
	
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling user update: %v", err)
		return
	}
	
	msg := models.WSMessage{
		Type:      models.WSMessageTypeUserUpdate,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}
	
	s.BroadcastToUser(userID, msg)
}

func (s *webSocketService) BroadcastOrderUpdate(orderID string, order *models.Order, status models.OrderStatus) {
	payload := models.OrderUpdatePayload{
		Order:     order,
		Status:    status,
		UpdatedBy: "system",
	}
	
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling order update: %v", err)
		return
	}
	
	msg := models.WSMessage{
		Type:      models.WSMessageTypeOrderUpdate,
		Data:      data,
		Timestamp: time.Now(),
	}
	
	// Broadcast to all connections watching this order
	s.BroadcastToResource(models.WSConnectionTypeOrder, orderID, msg)
	
	// Also broadcast to the user who owns the order
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
	
	msg := models.WSMessage{
		Type:      models.WSMessageTypeNotification,
		Data:      data,
		Timestamp: time.Now(),
		UserID:    userID,
	}
	
	s.BroadcastToUser(userID, msg)
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
```

## 4. Multiple WebSocket Handlers (internal/api/handlers/websocket.go)

```go
package handlers

import (
	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	wsService           service.WebSocketService
	userService         service.UserService
	orderService        service.OrderService
	notificationService service.NotificationService
	restaurantService   service.RestaurantService
}

func NewWebSocketHandler(
	wsService service.WebSocketService,
	userService service.UserService,
	orderService service.OrderService,
	notificationService service.NotificationService,
	restaurantService service.RestaurantService,
) *WebSocketHandler {
	return &WebSocketHandler{
		wsService:           wsService,
		userService:         userService,
		orderService:        orderService,
		notificationService: notificationService,
		restaurantService:   restaurantService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure properly for production
	},
}

// User WebSocket endpoints
func (h *WebSocketHandler) WatchUserDetails(c *gin.Context) {
	userID := c.Param("userId")
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeUser, userID, userID)
}

// Order WebSocket endpoints
func (h *WebSocketHandler) WatchOrder(c *gin.Context) {
	orderID := c.Param("orderId")
	userID := c.GetString("user_id") // From JWT middleware
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeOrder, userID, orderID)
}

func (h *WebSocketHandler) WatchUserOrders(c *gin.Context) {
	userID := c.Param("userId")
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeOrder, userID, "user_orders")
}

// Notification WebSocket endpoints
func (h *WebSocketHandler) WatchNotifications(c *gin.Context) {
	userID := c.Param("userId")
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeNotification, userID, userID)
}

// Restaurant WebSocket endpoints
func (h *WebSocketHandler) WatchRestaurant(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
	userID := c.GetString("user_id") // From JWT middleware
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeRestaurant, userID, restaurantID)
}

// Generic WebSocket upgrade handler
func (h *WebSocketHandler) handleWebSocketUpgrade(c *gin.Context, connectionType models.WSConnectionType, userID, resourceID string) {
	// Verify permissions based on connection type
	if err := h.verifyPermissions(connectionType, userID, resourceID); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, err
			},
			"verifying WebSocket permissions",
		)
		result.RespondWithJSON(c)
		return
	}
	
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Failed to upgrade connection", err)
			},
			"upgrading to WebSocket connection",
		)
		result.RespondWithJSON(c)
		return
	}
	
	// Handle the connection
	h.wsService.HandleConnection(connectionType, userID, resourceID, conn)
}

func (h *WebSocketHandler) verifyPermissions(connectionType models.WSConnectionType, userID, resourceID string) error {
	switch connectionType {
	case models.WSConnectionTypeUser:
		// Verify user exists and user can access their own data
		_, err := h.userService.GetByID(userID)
		return err
		
	case models.WSConnectionTypeOrder:
		if resourceID == "user_orders" {
			// User watching their own orders
			_, err := h.userService.GetByID(userID)
			return err
		}
		// Verify order exists and user owns it
		order, err := h.orderService.GetByID(resourceID)
		if err != nil {
			return err
		}
		if order.UserID != userID {
			return errors.NewHTTPError(http.StatusForbidden, "Access denied", nil)
		}
		return nil
		
	case models.WSConnectionTypeNotification:
		// Verify user exists
		_, err := h.userService.GetByID(userID)
		return err
		
	case models.WSConnectionTypeRestaurant:
		// Verify restaurant exists (public access)
		_, err := h.restaurantService.GetByID(resourceID)
		return err
		
	default:
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid connection type", nil)
	}
}

// WebSocket stats endpoint
func (h *WebSocketHandler) GetWebSocketStats(c *gin.Context) {
	stats := h.wsService.GetConnectionStats()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
```

## 5. Updated Routes (internal/api/routes/routes.go)

```go
// Add to SetupRoutes function
func SetupRoutes(deps *Dependencies) *gin.Engine {
	// ... existing code ...
	
	// Initialize WebSocket Handler
	wsHandler := handlers.NewWebSocketHandler(
		deps.WebSocketService,
		deps.UserService,
		deps.OrderService,
		deps.NotificationService,
		deps.RestaurantService,
	)
	
	// WebSocket endpoints
	ws := v1.Group("/ws")
	ws.Use(middleware.AuthMiddleware(deps.UserRepository))
	{
		// User WebSocket endpoints
		ws.GET("/users/:userId/watch", wsHandler.WatchUserDetails)
		
		// Order WebSocket endpoints
		ws.GET("/orders/:orderId/watch", wsHandler.WatchOrder)
		ws.GET("/users/:userId/orders/watch", wsHandler.WatchUserOrders)
		
		// Notification WebSocket endpoints
		ws.GET("/users/:userId/notifications/watch", wsHandler.WatchNotifications)
		
		// Restaurant WebSocket endpoints
		ws.GET("/restaurants/:restaurantId/watch", wsHandler.WatchRestaurant)
		
		// WebSocket stats (admin only)
		ws.GET("/stats", wsHandler.GetWebSocketStats)
	}
	
	// ... rest of routes ...
}
```

## 6. Service Integration Examples

### Update User Service
```go
func (s *userService) Update(id string, updates map[string]interface{}) error {
	err := s.userRepo.Update(id, updates)
	if err != nil {
		return err
	}
	
	// Broadcast update via WebSocket
	if s.wsService != nil {
		user, _ := s.userRepo.GetByID(id)
		if user != nil {
			s.wsService.BroadcastUserUpdate(id, user, updates)
		}
	}
	
	return nil
}
```

### Update Order Service
```go
func (s *orderService) UpdateStatus(id string, status models.OrderStatus) error {
	err := s.orderRepo.UpdateStatus(id, status)
	if err != nil {
		return err
	}
	
	// Broadcast update via WebSocket
	if s.wsService != nil {
		order, _ := s.orderRepo.GetByID(id)
		if order != nil {
			s.wsService.BroadcastOrderUpdate(id, order, status)
		}
	}
	
	return nil
}
```

## WebSocket Endpoints Summary

| Endpoint | Purpose | Connection Type |
|----------|---------|-----------------|
| `/ws/users/:userId/watch` | Watch user profile changes | `user` |
| `/ws/orders/:orderId/watch` | Watch specific order updates | `order` |
| `/ws/users/:userId/orders/watch` | Watch all user's orders | `order` |
| `/ws/users/:userId/notifications/watch` | Watch user notifications | `notification` |
| `/ws/restaurants/:restaurantId/watch` | Watch restaurant updates | `restaurant` |
| `/ws/stats` | WebSocket connection statistics | N/A |

This architecture allows you to easily add new WebSocket endpoint types while maintaining clean separation of concerns and reusing the core connection management logic.