# Simplified WebSocket Implementation

## 🎯 Clean Architecture Approach (Flutter-Style)

This simplified implementation follows the same clean architecture you know from Flutter:
- **Models**: Data structures
- **Services**: Business logic 
- **Handlers**: Presentation layer (like Controllers in Flutter)
- **Routes**: Routing (like GoRouter in Flutter)

---

## Step 1: Simple WebSocket Models

**File: `internal/models/websocket_simple.go`**

```go
package models

import (
    "time"
    "github.com/gorilla/websocket"
)

// Simple message types - just what we need
const (
    MessageTypeUserUpdate = "user_update"
    MessageTypeConnected  = "connected"
    MessageTypeError      = "error"
)

// Simple WebSocket message - like a DTO in Flutter
type SimpleWSMessage struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}

// Simple connection info - like a model in Flutter
type SimpleConnection struct {
    UserID     string
    Connection *websocket.Conn
}

// User update data - like a specific DTO
type UserUpdateData struct {
    User    *User                  `json:"user"`
    Changes map[string]interface{} `json:"changes"`
}
```

**Why this is simple:**
- Only 3 message types (not 15+)
- Basic data structures
- No complex routing logic

---

## Step 2: Simple WebSocket Service

**File: `internal/service/websocket_simple.go`**

```go
package service

import (
    "dfood/internal/models"
    "encoding/json"
    "log"
    "sync"
    "time"
    "github.com/gorilla/websocket"
)

// Simple WebSocket service - like a service in Flutter
type SimpleWebSocketService struct {
    // Map of userID -> connection (one connection per user for simplicity)
    connections map[string]*websocket.Conn
    mutex       sync.RWMutex
}

// Constructor - like creating a service in Flutter
func NewSimpleWebSocketService() *SimpleWebSocketService {
    return &SimpleWebSocketService{
        connections: make(map[string]*websocket.Conn),
    }
}

// Add a user connection
func (s *SimpleWebSocketService) AddConnection(userID string, conn *websocket.Conn) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    // Close existing connection if any
    if existingConn, exists := s.connections[userID]; exists {
        existingConn.Close()
    }
    
    s.connections[userID] = conn
    log.Printf("WebSocket: User %s connected", userID)
    
    // Send connection confirmation
    s.sendToUser(userID, models.SimpleWSMessage{
        Type:      models.MessageTypeConnected,
        Data:      map[string]string{"message": "Connected successfully"},
        Timestamp: time.Now(),
    })
}

// Remove a user connection
func (s *SimpleWebSocketService) RemoveConnection(userID string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if conn, exists := s.connections[userID]; exists {
        conn.Close()
        delete(s.connections, userID)
        log.Printf("WebSocket: User %s disconnected", userID)
    }
}

// Send user update - this is the main function you'll use
func (s *SimpleWebSocketService) SendUserUpdate(userID string, user *models.User, changes map[string]interface{}) {
    updateData := models.UserUpdateData{
        User:    user,
        Changes: changes,
    }
    
    message := models.SimpleWSMessage{
        Type:      models.MessageTypeUserUpdate,
        Data:      updateData,
        Timestamp: time.Now(),
    }
    
    s.sendToUser(userID, message)
}

// Private method to send message to a specific user
func (s *SimpleWebSocketService) sendToUser(userID string, message models.SimpleWSMessage) {
    s.mutex.RLock()
    conn, exists := s.connections[userID]
    s.mutex.RUnlock()
    
    if !exists {
        return // User not connected
    }
    
    if err := conn.WriteJSON(message); err != nil {
        log.Printf("WebSocket: Error sending to user %s: %v", userID, err)
        s.RemoveConnection(userID) // Clean up bad connection
    }
}

// Handle a WebSocket connection (keeps it alive)
func (s *SimpleWebSocketService) HandleConnection(userID string, conn *websocket.Conn) {
    s.AddConnection(userID, conn)
    defer s.RemoveConnection(userID)
    
    // Keep connection alive by reading messages (even if we ignore them)
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket: Connection error for user %s: %v", userID, err)
            break
        }
    }
}
```

**Why this is simple:**
- One connection per user (not multiple types)
- Only handles user updates (not orders, notifications, etc.)
- Clear, single-purpose methods
- Like a Flutter service with clear responsibilities

---

## Step 3: Update User Service (Minimal Changes)

**File: `internal/service/user_service.go`** - Add these methods:

```go
// Add to UserService interface
type UserService interface {
    // ... existing methods ...
    SetWebSocketService(ws *SimpleWebSocketService)
}

// Add to userService struct
type userService struct {
    // ... existing fields ...
    wsService *SimpleWebSocketService
}

// Add these methods
func (s *userService) SetWebSocketService(ws *SimpleWebSocketService) {
    s.wsService = ws
}

// Update your existing Update method
func (s *userService) Update(userID string, updates map[string]interface{}) error {
    // 1. Update database (existing code)
    err := s.userRepo.Update(userID, updates)
    if err != nil {
        return err
    }
    
    // 2. Send WebSocket update (NEW - just 3 lines!)
    if s.wsService != nil {
        user, _ := s.userRepo.GetByID(userID)
        s.wsService.SendUserUpdate(userID, user, updates)
    }
    
    return nil
}

// Update your existing UpdateField method
func (s *userService) UpdateField(userID, field string, value interface{}) error {
    // 1. Update database (existing code)
    err := s.userRepo.UpdateField(userID, field, value)
    if err != nil {
        return err
    }
    
    // 2. Send WebSocket update (NEW - just 3 lines!)
    if s.wsService != nil {
        user, _ := s.userRepo.GetByID(userID)
        changes := map[string]interface{}{field: value}
        s.wsService.SendUserUpdate(userID, user, changes)
    }
    
    return nil
}
```

**Why this is simple:**
- Only 3 lines added to existing methods
- No complex interfaces or dependencies
- Clear separation of concerns

---

## Step 4: Simple WebSocket Handler

**File: `internal/api/handlers/websocket_simple.go`**

```go
package handlers

import (
    "dfood/internal/service"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type SimpleWebSocketHandler struct {
    wsService   *service.SimpleWebSocketService
    userService service.UserService
}

func NewSimpleWebSocketHandler(wsService *service.SimpleWebSocketService, userService service.UserService) *SimpleWebSocketHandler {
    return &SimpleWebSocketHandler{
        wsService:   wsService,
        userService: userService,
    }
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for development
    },
}

// Simple user profile WebSocket endpoint
func (h *SimpleWebSocketHandler) WatchUserProfile(c *gin.Context) {
    userID := c.Param("userId")
    
    // Verify user exists
    _, err := h.userService.GetByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    // Upgrade to WebSocket
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
        return
    }
    
    // Handle the connection (this blocks until connection closes)
    h.wsService.HandleConnection(userID, conn)
}
```

**Why this is simple:**
- One handler method
- Basic error handling
- No complex routing logic

---

## Step 5: Simple Routes

**File: `internal/api/routes/routes.go`** - Add this:

```go
func SetupRoutes(deps *Dependencies) *gin.Engine {
    // ... existing code ...
    
    // Initialize simple WebSocket service
    wsService := service.NewSimpleWebSocketService()
    
    // Set WebSocket service in user service
    deps.UserService.SetWebSocketService(wsService)
    
    // Initialize simple WebSocket handler
    wsHandler := handlers.NewSimpleWebSocketHandler(wsService, deps.UserService)
    
    // ... existing routes ...
    
    // Add simple WebSocket route
    v1.GET("/users/:userId/watch", wsHandler.WatchUserProfile)
    
    return router
}
```

**Why this is simple:**
- One route
- Clear initialization
- No complex dependencies

---

## Step 6: Update Main (Minimal)

**File: `cmd/main.go`** - No changes needed! The routes handle everything.

---

## 🎯 How to Use (Like Flutter)

### 1. **Connect from Flutter:**

```dart
// Connect to WebSocket
final ws = WebSocketChannel.connect(
    Uri.parse('ws://localhost:8080/api/v1/users/user123/watch')
);

// Listen for messages
ws.stream.listen((message) {
    final data = jsonDecode(message);
    
    switch (data['type']) {
        case 'connected':
            print('Connected to WebSocket');
            break;
        case 'user_update':
            // Update your user model
            final userData = data['data'];
            updateUserInUI(userData['user'], userData['changes']);
            break;
    }
});
```

### 2. **Update User Profile:**

```bash
# This will automatically send WebSocket update
curl -X PUT http://localhost:8080/api/v1/users/user123 \
  -H "Content-Type: application/json" \
  -d '{"first_name": "John Updated"}'
```

### 3. **WebSocket Message Received:**

```json
{
  "type": "user_update",
  "data": {
    "user": {
      "id": "user123",
      "first_name": "John Updated",
      "last_name": "Doe"
    },
    "changes": {
      "first_name": "John Updated"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 🎓 Key Simplifications

| Complex Version | Simple Version |
|----------------|----------------|
| 6+ connection types | 1 connection type (user) |
| 15+ message types | 3 message types |
| Complex routing | Direct user-to-connection mapping |
| Multiple broadcast methods | 1 send method |
| Resource-based connections | User-based connections |
| Complex payload types | Simple data structures |

## 🚀 Benefits

1. **Easy to Understand**: Like Flutter services you're used to
2. **Easy to Extend**: Add more message types as needed
3. **Easy to Debug**: Simple connection mapping
4. **Production Ready**: Still handles errors and cleanup
5. **Clean Architecture**: Clear separation of concerns

## 📱 Flutter Analogy

```dart
// This WebSocket service is like a Flutter service:
class SimpleWebSocketService {
  Map<String, WebSocketConnection> connections = {};
  
  void addConnection(String userId, WebSocketConnection conn) { }
  void sendUserUpdate(String userId, User user, Map changes) { }
}

// The handler is like a Flutter controller:
class SimpleWebSocketHandler {
  void watchUserProfile(String userId) { }
}
```

This simplified version gives you 80% of the functionality with 20% of the complexity!