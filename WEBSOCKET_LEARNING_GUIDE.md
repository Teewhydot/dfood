# WebSocket Learning Guide for dfood

## 🎯 What Are WebSockets?

WebSockets provide **real-time, bidirectional communication** between client and server. Unlike traditional HTTP requests where the client always initiates communication, WebSockets allow the server to push data to clients instantly.

### Traditional HTTP vs WebSocket

```
HTTP Request/Response:
Client → Server: "Give me user data"
Server → Client: "Here's the user data"
[Connection closes]

WebSocket:
Client ↔ Server: [Persistent connection]
Server → Client: "User data changed!" (anytime)
Client → Server: "Thanks, got it!"
```

## 🏗️ How Our WebSocket System Works

### 1. Connection Flow

```
1. Client connects to WebSocket endpoint
   GET /ws/users/123/watch
   
2. Server upgrades HTTP to WebSocket
   HTTP 101 Switching Protocols
   
3. Server registers connection in memory
   connectionID: "abc-123"
   type: "user"
   userID: "123"
   
4. Client receives confirmation
   {"type": "connected", "connection_id": "abc-123"}
   
5. Connection stays open for real-time updates
```

### 2. Message Broadcasting Flow

```
1. User updates profile via REST API
   PUT /api/v1/users/123 {"first_name": "John"}
   
2. UserService.Update() saves to database
   
3. UserService.Update() calls WebSocket broadcast
   wsService.BroadcastUserUpdate(userID, user, changes)
   
4. WebSocket service finds all connections watching user 123
   
5. Sends update message to all those connections
   {"type": "user_update", "data": {...}}
   
6. Client receives message and updates UI instantly
```

## 📋 Key Components Explained

### Connection Registry

```go
// We store all active connections in memory
connections map[string]*models.WSConnection

// Each connection has metadata
type WSConnection struct {
    ID         string              // Unique identifier
    Type       WSConnectionType    // What data they want (user, order, etc.)
    UserID     string              // Who owns this connection
    ResourceID string              // Specific resource (optional)
    Connection *websocket.Conn     // The actual WebSocket
}
```

**Why this matters:** We need to know which clients want which updates so we can route messages correctly.

### Message Types

```go
// Every WebSocket message has a type so clients know how to handle it
const (
    WSMessageTypeUserUpdate    = "user_update"    // Profile changed
    WSMessageTypeOrderUpdate   = "order_update"   // Order status changed
    WSMessageTypeNotification  = "notification"   // New notification
    // ... etc
)
```

**Why this matters:** Clients can have different UI logic for different message types.

### Broadcasting Strategies

```go
// Send to one specific connection
BroadcastToConnection(connectionID, message)

// Send to all connections of one user (multiple devices)
BroadcastToUser(userID, message)

// Send to all connections watching a specific resource
BroadcastToResource(connectionType, resourceID, message)
```

**Why this matters:** Different scenarios need different targeting strategies.

## 🔄 Real-World Examples

### Example 1: User Profile Update

**Scenario:** User updates their name on their phone, tablet should update automatically.

```go
// 1. REST API call updates database
func (s *userService) Update(userID string, updates map[string]interface{}) error {
    // Save to database
    err := s.userRepo.Update(userID, updates)
    if err != nil {
        return err
    }
    
    // 2. Broadcast to WebSocket connections
    if s.wsService != nil {
        user, _ := s.userRepo.GetByID(userID)
        s.wsService.BroadcastUserUpdate(userID, user, updates)
    }
    
    return nil
}

// 3. WebSocket service routes message
func (s *webSocketService) BroadcastUserUpdate(userID string, user *User, changes map[string]interface{}) {
    // Create message payload
    payload := UserUpdatePayload{
        User:    user,
        Changes: changes,  // {"first_name": "John"}
    }
    
    // Send to all connections watching this user
    s.BroadcastToResource(WSConnectionTypeUser, userID, message)
}
```

**Result:** All devices connected to `/ws/users/123/watch` receive the update instantly.

### Example 2: Order Status Tracking

**Scenario:** Restaurant marks order as "ready", customer gets notified immediately.

```go
// 1. Restaurant updates order status
func (s *orderService) UpdateStatus(orderID string, status OrderStatus) error {
    // Update database
    err := s.orderRepo.UpdateStatus(orderID, status)
    
    // 2. Broadcast to WebSocket
    order, _ := s.orderRepo.GetByID(orderID)
    s.wsService.BroadcastOrderUpdate(orderID, order, status)
    
    return nil
}

// 3. Message goes to order watchers AND order owner
func (s *webSocketService) BroadcastOrderUpdate(orderID string, order *Order, status OrderStatus) {
    // Send to anyone watching this specific order
    s.BroadcastToResource(WSConnectionTypeOrder, orderID, msg)
    
    // Also send to the customer who owns the order
    s.BroadcastToUser(order.UserID, msg)
}
```

**Result:** Customer sees "Order Ready!" notification in real-time.

## 🛠️ Implementation Patterns

### Pattern 1: Service Layer Integration

```go
// Every service that needs WebSocket updates follows this pattern:

type UserService interface {
    Update(userID string, updates map[string]interface{}) error
    SetWebSocketService(wsService WebSocketService)  // Inject WebSocket service
}

type userService struct {
    userRepo  UserRepository
    wsService WebSocketService  // Store reference
}

func (s *userService) Update(userID string, updates map[string]interface{}) error {
    // 1. Update database
    err := s.userRepo.Update(userID, updates)
    
    // 2. Broadcast via WebSocket (if available)
    if s.wsService != nil {
        // Fetch updated data and broadcast
        user, _ := s.userRepo.GetByID(userID)
        s.wsService.BroadcastUserUpdate(userID, user, updates)
    }
    
    return err
}
```

### Pattern 2: Connection Type Routing

```go
// Different connection types get different message routing:

switch connectionType {
case WSConnectionTypeUser:
    // User profile updates go to that user's connections
    s.BroadcastToResource(WSConnectionTypeUser, userID, message)
    
case WSConnectionTypeOrder:
    // Order updates go to order watchers AND order owner
    s.BroadcastToResource(WSConnectionTypeOrder, orderID, message)
    s.BroadcastToUser(order.UserID, message)
    
case WSConnectionTypeNotification:
    // Notifications go to that user only
    s.BroadcastToUser(userID, message)
}
```

### Pattern 3: Message Payload Structure

```go
// Every WebSocket message follows the same structure:
type WSMessage struct {
    Type      WSMessageType   `json:"type"`       // What kind of message
    Data      json.RawMessage `json:"data"`       // The actual payload
    Timestamp time.Time       `json:"timestamp"`  // When it happened
    UserID    string          `json:"user_id"`    // Who it's for
}

// Specific payloads go in the Data field:
type UserUpdatePayload struct {
    User    *User                  `json:"user"`     // Complete updated object
    Changes map[string]interface{} `json:"changes"`  // What changed
}
```

## 🔧 Connection Management

### Health Monitoring

```go
// We use ping/pong to detect dead connections
func (s *webSocketService) HandleConnection(connectionType, userID, resourceID, conn) {
    // Set up ping/pong handlers
    conn.SetPongHandler(func(string) error {
        // Update last ping time when we receive pong
        wsConn.LastPing = time.Now()
        return nil
    })
    
    // Send ping every 30 seconds
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            conn.WriteMessage(websocket.PingMessage, nil)
        }
    }()
}
```

### Cleanup

```go
// Always clean up when connections close
func (s *webSocketService) HandleConnection(...) {
    connectionID := s.RegisterConnection(...)
    defer s.UnregisterConnection(connectionID)  // Cleanup on exit
    
    // Handle messages until connection closes
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break  // Connection closed, cleanup will happen
        }
    }
}
```

## 📱 Client-Side (Flutter) Integration

### Connection Management

```dart
class WebSocketService {
    WebSocketChannel? _channel;
    
    Future<void> connect(String endpoint) async {
        // Connect to WebSocket endpoint
        _channel = WebSocketChannel.connect(Uri.parse(endpoint));
        
        // Listen for messages
        _channel!.stream.listen((message) {
            final wsMessage = WSMessage.fromJson(jsonDecode(message));
            _handleMessage(wsMessage);
        });
    }
    
    void _handleMessage(WSMessage message) {
        switch (message.type) {
            case 'user_update':
                // Update user UI
                final payload = UserUpdatePayload.fromJson(message.data);
                _updateUserUI(payload.user, payload.changes);
                break;
                
            case 'order_update':
                // Update order status UI
                final payload = OrderUpdatePayload.fromJson(message.data);
                _updateOrderUI(payload.order, payload.status);
                break;
        }
    }
}
```

### UI Updates

```dart
// Use Provider or similar state management
class UserProvider extends ChangeNotifier {
    User? _user;
    
    void connectWebSocket() {
        wsService.userUpdateStream.listen((updatedUser) {
            _user = updatedUser;
            notifyListeners();  // Triggers UI rebuild
        });
    }
}

// UI automatically updates when data changes
Consumer<UserProvider>(
    builder: (context, userProvider, child) {
        return Text(userProvider.user?.firstName ?? '');
    },
)
```

## 🚀 Benefits of This Architecture

### 1. **Real-Time Updates**
- No need to refresh or poll for changes
- Instant UI updates across all devices
- Better user experience

### 2. **Scalable Design**
- Easy to add new WebSocket endpoint types
- Centralized connection management
- Efficient message routing

### 3. **Clean Integration**
- WebSocket logic separated from business logic
- Existing API endpoints automatically broadcast updates
- No code duplication

### 4. **Production Ready**
- Connection health monitoring
- Automatic cleanup
- Error handling and logging
- Performance monitoring

## 🎓 Key Learning Points

1. **WebSockets enable real-time communication** - server can push data to clients anytime
2. **Connection registry is crucial** - we need to track who wants what updates
3. **Message types provide structure** - clients know how to handle different updates
4. **Broadcasting strategies matter** - different scenarios need different targeting
5. **Integration with existing services** - WebSocket updates happen automatically when data changes
6. **Client-side state management** - UI updates reactively to WebSocket messages

This architecture provides a solid foundation for real-time features while maintaining clean, maintainable code!