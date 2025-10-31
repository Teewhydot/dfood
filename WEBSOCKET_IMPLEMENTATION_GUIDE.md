# Complete WebSocket Implementation Guide for Go Backend

## 📚 Table of Contents
1. [Understanding WebSockets](#understanding-websockets)
2. [Project Architecture Overview](#project-architecture-overview)
3. [The 4 Essential Pieces](#the-4-essential-pieces)
4. [Step-by-Step Implementation](#step-by-step-implementation)
5. [Practical Exercise: Implement Order Tracking WebSocket](#practical-exercise)
6. [Testing Your Implementation](#testing-your-implementation)
7. [Common Pitfalls and Solutions](#common-pitfalls-and-solutions)
8. [Troubleshooting Guide](#troubleshooting-guide)

---

## Understanding WebSockets

### What is a WebSocket?

Think of a WebSocket like a **phone call** between your app and the server:
- **Traditional HTTP**: Like sending letters - you send a request, wait for response, connection closes
- **WebSocket**: Like a phone call - connection stays open, both sides can talk anytime

### Why Use WebSockets?

**Without WebSocket (Polling):**
```
Client: "Any new messages?" → Server: "No"
(wait 5 seconds)
Client: "Any new messages?" → Server: "No"
(wait 5 seconds)
Client: "Any new messages?" → Server: "Yes, here's one!"
```
*Problems: Waste of resources, delays, battery drain*

**With WebSocket:**
```
Client: [Keeps connection open]
Server: "New message!" → Client receives instantly
Server: "Another message!" → Client receives instantly
```
*Benefits: Instant updates, efficient, real-time*

### Real-World Examples in This Project

1. **User Profile Updates**: When a user changes their name, all their devices see it instantly
2. **Address Updates**: Add an address on phone, see it on tablet immediately
3. **Favorites**: Like a food item, it appears in favorites list instantly
4. **Notifications**: Server pushes notifications to user in real-time

---

## Project Architecture Overview

### The Big Picture

```
┌─────────────────────────────────────────────────────────────┐
│                     Flutter App (Client)                     │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  Profile   │  │ Favorites  │  │  Orders    │            │
│  │  Screen    │  │  Screen    │  │  Screen    │            │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘            │
│        │                │               │                    │
│        └────────────────┼───────────────┘                    │
│                         │                                    │
│                  WebSocket Connection                        │
└─────────────────────────┼────────────────────────────────────┘
                          │
                          │ ws://localhost:8080/api/v1/realtime/...
                          │
┌─────────────────────────┼────────────────────────────────────┐
│                    Go Backend                                 │
│                         │                                     │
│  ┌──────────────────────▼───────────────────────┐            │
│  │         RealtimeHandler (Entry Point)        │            │
│  │  • Validates user                            │            │
│  │  • Upgrades HTTP → WebSocket                 │            │
│  │  • Calls RealtimeService                     │            │
│  └──────────────────┬───────────────────────────┘            │
│                     │                                         │
│  ┌──────────────────▼───────────────────────────┐            │
│  │         RealtimeService (Connection Manager)  │            │
│  │  • Stores all WebSocket connections          │            │
│  │  • Sends messages to connected users         │            │
│  │  • Manages connection lifecycle              │            │
│  └──────────────────┬───────────────────────────┘            │
│                     │                                         │
│  ┌──────────────────▼───────────────────────────┐            │
│  │      Business Logic Services                  │            │
│  │  • UserService                                │            │
│  │  • AddressService  ───► Calls RealtimeService │            │
│  │  • FavoritesService    to broadcast updates  │            │
│  │  • NotificationService                        │            │
│  └───────────────────────────────────────────────┘            │
└───────────────────────────────────────────────────────────────┘
```

### How Data Flows

**Example: User updates their profile**

```
1. User taps "Save" on Flutter app
   ↓
2. App sends HTTP PATCH to /api/v1/users/123
   ↓
3. UserService.Update() processes the update
   ↓
4. UserService calls RealtimeService.SendUserUpdate()
   ↓
5. RealtimeService finds all WebSocket connections for user123
   ↓
6. RealtimeService sends update message through WebSocket
   ↓
7. All user's devices receive update instantly via their WebSocket connections
```

---

## The 4 Essential Pieces

Every WebSocket feature needs these 4 pieces. Think of them like building blocks:

### Piece 1: Message Types and Data Models
**Location:** `internal/models/websocket.go`

**What it does:** Defines what kind of messages you can send and what data they contain.

**Analogy:** Like creating different types of envelopes for different types of mail.

```go
// Message types - like labels on envelopes
const (
    MessageTypeOrderUpdate = "order_update"  // Label for order messages
)

// Message structure - the envelope format
type WSMessage struct {
    Type      string      `json:"type"`      // What kind of message
    Data      interface{} `json:"data"`      // The actual content
    Timestamp time.Time   `json:"timestamp"` // When it was sent
}

// Data structure - what goes inside the envelope
type OrderUpdateData struct {
    Order  *Order       `json:"order"`   // The order details
    Status OrderStatus  `json:"status"`  // Current status
}
```

### Piece 2: RealtimeService (Connection Manager)
**Location:** `internal/service/realtime_service.go`

**What it does:** Manages all WebSocket connections and sends messages to users.

**Analogy:** Like a post office that knows everyone's address and delivers mail.

```go
type RealtimeService struct {
    connections map[string]*websocket.Conn  // Map of userID → their connection
    mutex       sync.RWMutex                // Lock for thread-safety
}

// Key methods:
// - AddConnection(): Register a new user connection
// - RemoveConnection(): Remove user when they disconnect
// - SendOrderUpdate(): Send order updates to specific user
// - sendToUser(): Private method that actually sends the message
```

### Piece 3: RealtimeHandler (Entry Point)
**Location:** `internal/api/handlers/realtime_handler.go`

**What it does:** Handles incoming WebSocket connection requests from clients.

**Analogy:** Like a receptionist who greets visitors and connects them to the right person.

```go
func (h *RealtimeHandler) WatchOrders(c *gin.Context) {
    userID := c.Param("userId")

    // 1. Validate user exists (check ID)
    _, err := h.userService.GetByID(userID)

    // 2. Upgrade HTTP connection to WebSocket (open phone line)
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

    // 3. Hand over to RealtimeService (transfer call)
    h.realtimeService.HandleConnection(userID, conn)
}
```

### Piece 4: Service Integration
**Location:** `internal/service/order_service.go` (or any business service)

**What it does:** Calls RealtimeService when data changes to notify connected users.

**Analogy:** Like a store employee calling you when your order is ready.

```go
func (s *orderService) UpdateStatus(orderID string, status OrderStatus) error {
    // 1. Update the order in database
    err := s.orderRepo.UpdateStatus(orderID, status)
    if err != nil {
        return err
    }

    // 2. Send real-time update to user
    if s.realtimeService != nil {
        order, _ := s.orderRepo.GetByID(orderID)
        s.realtimeService.SendOrderUpdate(order.UserID, order, status)
    }

    return nil
}
```

---

## Step-by-Step Implementation

Let's implement a complete WebSocket feature for **Order Tracking**. Follow these steps exactly:

### Step 1: Add Message Types

**File:** `internal/models/websocket.go`

**Task:** Add new message type constants for orders.

```go
// Find the const block (around line 10-31) and add:
const (
    // ... existing types ...

    // Order messages (ADD THESE)
    MessageTypeOrderUpdate  = "order_update"
    MessageTypeOrderStatus  = "order_status"
    MessageTypeOrderReady   = "order_ready"
)
```

**Why:** These constants define what types of order messages we can send. Using constants prevents typos.

### Step 2: Add Data Structures

**File:** `internal/models/websocket.go`

**Task:** Add data structure for order updates.

```go
// Add at the bottom of the file (after NotificationUpdateData):

// OrderUpdateData - for order changes
type OrderUpdateData struct {
    Order     *Order      `json:"order"`      // Complete order object
    Status    OrderStatus `json:"status"`     // Current order status
    UpdatedAt time.Time   `json:"updated_at"` // When status changed
    Message   string      `json:"message"`    // User-friendly message
}
```

**Why:** This structure defines exactly what data we send when an order updates. JSON tags tell Go how to convert this to JSON for the client.

**Key Points:**
- `*Order` = pointer to Order model (the asterisk means pointer)
- `json:"order"` = this field will be called "order" in JSON
- `OrderStatus` = type from your order model (like "pending", "preparing", etc.)

### Step 3: Add Service Methods

**File:** `internal/service/realtime_service.go`

**Task:** Add methods to send order updates.

```go
// Add these methods at the bottom of the file (before the sendToUser method):

// Order Updates

// SendOrderUpdate broadcasts when order status changes
func (s *RealtimeService) SendOrderUpdate(userID string, order *models.Order, status models.OrderStatus, message string) {
    // Create the data package
    updateData := models.OrderUpdateData{
        Order:     order,
        Status:    status,
        UpdatedAt: time.Now(),
        Message:   message,
    }

    // Create the WebSocket message
    wsMessage := models.WSMessage{
        Type:      models.MessageTypeOrderUpdate,
        Data:      updateData,
        Timestamp: time.Now(),
    }

    // Send to user
    s.sendToUser(userID, wsMessage)
    log.Printf("Realtime: Sent order update to %s, order: %s, status: %s",
        userID, order.ID, status)
}

// SendOrderReady broadcasts when order is ready for pickup
func (s *RealtimeService) SendOrderReady(userID string, order *models.Order) {
    // Create the data package
    updateData := models.OrderUpdateData{
        Order:     order,
        Status:    order.Status,
        UpdatedAt: time.Now(),
        Message:   "Your order is ready for pickup!",
    }

    // Create the WebSocket message
    wsMessage := models.WSMessage{
        Type:      models.MessageTypeOrderReady,
        Data:      updateData,
        Timestamp: time.Now(),
    }

    // Send to user
    s.sendToUser(userID, wsMessage)
    log.Printf("Realtime: Sent order ready notification to %s, order: %s",
        userID, order.ID)
}
```

**Why:** These methods make it easy to send order updates. You just call `SendOrderUpdate()` and it handles all the message creation and sending.

**Understanding the Code:**
1. `updateData` = The actual order information we want to send
2. `wsMessage` = Wraps the data in a standard message format
3. `s.sendToUser()` = Actually sends the message through the WebSocket
4. `log.Printf()` = Logs to console for debugging

### Step 4: Create Handler Methods

**File:** `internal/api/handlers/realtime_handler.go`

**Task:** Add handler for order WebSocket endpoint.

```go
// Add this method in the RealtimeHandler (after WatchNotifications):

// Order WebSocket Endpoints

// WatchOrder handles real-time order tracking WebSocket
// This allows a user to watch their specific order for updates
func (h *RealtimeHandler) WatchOrder(c *gin.Context) {
    // Extract IDs from URL parameters
    userID := c.Param("userId")
    orderID := c.Param("orderId")

    // Step 1: Validate user exists
    _, err := h.userService.GetByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "User not found",
            "user_id": userID,
        })
        return
    }

    // Step 2: Validate order exists and belongs to user
    // Note: You'll need to add this method to your handler dependencies
    // For now, we'll skip this check, but in production you should verify
    // that the order belongs to the user

    // Step 3: Upgrade HTTP connection to WebSocket
    conn, err := realtimeUpgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Failed to upgrade to WebSocket",
            "details": err.Error(),
        })
        return
    }

    // Step 4: Handle the connection (this blocks until connection closes)
    h.realtimeService.HandleConnection(userID, conn)

    // Note: Code after HandleConnection only runs when user disconnects
}

// WatchUserOrders handles real-time tracking of all user's orders
func (h *RealtimeHandler) WatchUserOrders(c *gin.Context) {
    userID := c.Param("userId")

    // Validate user exists
    _, err := h.userService.GetByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "User not found",
            "user_id": userID,
        })
        return
    }

    // Upgrade to WebSocket
    conn, err := realtimeUpgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Failed to upgrade to WebSocket",
            "details": err.Error(),
        })
        return
    }

    // Handle the connection
    h.realtimeService.HandleConnection(userID, conn)
}
```

**Why:** These handlers are the entry points - when a client connects, these methods run first. They validate the user and upgrade the connection.

**Understanding the Flow:**
1. Extract `userID` from URL (e.g., `/realtime/users/123/orders/watch` → userID = 123)
2. Check if user exists (security check)
3. Upgrade HTTP → WebSocket (like opening a phone line)
4. Pass connection to RealtimeService (transfer the call to the manager)

**Common Questions:**
- **Why check if user exists?** Security - don't want unauthorized connections
- **What does Upgrade() do?** Converts regular HTTP connection to WebSocket
- **Why does HandleConnection block?** It keeps running until user disconnects

### Step 5: Register Routes

**File:** `internal/api/routes/routes.go`

**Task:** Register the new WebSocket endpoints.

```go
// Find the realtime group (around line 226-246) and add:

realtime := v1.Group("/realtime")
realtime.Use(middleware.AuthMiddleware(deps.UserRepository))
{
    // ... existing endpoints ...

    // Order tracking endpoints (ADD THESE)
    realtime.GET("/users/:userId/orders/watch", realtimeHandler.WatchUserOrders)
    realtime.GET("/users/:userId/orders/:orderId/watch", realtimeHandler.WatchOrder)

    // Connection statistics
    realtime.GET("/stats", realtimeHandler.GetConnectionStats)
}
```

**Why:** Routes tell the server which function to call when a URL is accessed. Like a phone directory.

**Understanding the Routes:**
- `/users/:userId/orders/watch` = Watch ALL orders for a user
- `/users/:userId/orders/:orderId/watch` = Watch ONE specific order
- `:userId` and `:orderId` = Placeholders that capture values from URL

**Example URLs:**
- `ws://localhost:8080/api/v1/realtime/users/user123/orders/watch`
- `ws://localhost:8080/api/v1/realtime/users/user123/orders/order456/watch`

### Step 6: Integrate with Order Service

**File:** `internal/service/order_service.go`

**Task:** Add RealtimeService and call it when orders update.

#### 6.1: Add RealtimeService Field

```go
// Find the OrderService interface (around line 14-21) and add:

type OrderService interface {
    CreateOrder(order *models.Order) error
    GetByID(orderID string) (*models.Order, error)
    GetUserOrders(userID string) ([]models.Order, error)
    UpdateOrderStatus(orderID string, status models.OrderStatus) error
    CancelOrder(orderID string) error

    // Add this method
    SetRealtimeService(realtimeService *RealtimeService)
}

// Find the orderService struct (around line 23-27) and add:

type orderService struct {
    orderRepo    repository.OrderRepository
    userRepo     repository.UserRepository

    // Add this field
    realtimeService *RealtimeService
}
```

#### 6.2: Add Setter Method

```go
// Add this method after NewOrderService (around line 35):

func (s *orderService) SetRealtimeService(realtimeService *RealtimeService) {
    s.realtimeService = realtimeService
}
```

#### 6.3: Integrate in Update Method

```go
// Find UpdateOrderStatus method and modify it:

func (s *orderService) UpdateOrderStatus(orderID string, status models.OrderStatus) error {
    if strings.TrimSpace(orderID) == "" {
        return errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
    }

    // Validate order exists and get it
    order, err := s.orderRepo.GetByID(orderID)
    if err != nil {
        return err
    }

    // Update the status in database
    err = s.orderRepo.UpdateStatus(orderID, status)
    if err != nil {
        return err
    }

    // Send real-time update to user
    if s.realtimeService != nil {
        // Get the updated order
        updatedOrder, err := s.orderRepo.GetByID(orderID)
        if err == nil {
            // Create user-friendly message based on status
            var message string
            switch status {
            case models.OrderStatusPending:
                message = "Your order has been received"
            case models.OrderStatusConfirmed:
                message = "Your order has been confirmed"
            case models.OrderStatusPreparing:
                message = "Your order is being prepared"
            case models.OrderStatusReady:
                message = "Your order is ready for pickup!"
            case models.OrderStatusDelivering:
                message = "Your order is on the way"
            case models.OrderStatusDelivered:
                message = "Your order has been delivered"
            case models.OrderStatusCancelled:
                message = "Your order has been cancelled"
            default:
                message = "Your order status has been updated"
            }

            // Send the update
            s.realtimeService.SendOrderUpdate(
                updatedOrder.UserID,
                updatedOrder,
                status,
                message,
            )

            // Send special notification if order is ready
            if status == models.OrderStatusReady {
                s.realtimeService.SendOrderReady(updatedOrder.UserID, updatedOrder)
            }
        }
    }

    return nil
}
```

**Why:** This ensures that whenever an order status changes, connected users receive instant updates.

**Understanding the Integration:**
1. Update database first (source of truth)
2. Check if realtimeService exists (might be nil in tests)
3. Get the updated order to send fresh data
4. Create a user-friendly message for each status
5. Send the update through WebSocket
6. Send extra notification if order is ready

### Step 7: Wire Everything Together

**File:** `internal/api/routes/routes.go`

**Task:** Inject RealtimeService into OrderService.

```go
// Find where services are set up (around line 47-50) and add:

// Set RealtimeService in services that need it
deps.AddressService.SetRealtimeService(realtimeService)
deps.FavoritesService.SetRealtimeService(realtimeService)
deps.NotificationService.SetRealtimeService(realtimeService)
deps.OrderService.SetRealtimeService(realtimeService)  // ADD THIS LINE
```

**Why:** This connects the RealtimeService to the OrderService so they can work together.

**Understanding Dependency Injection:**
- Dependencies = things a service needs to work
- Injection = giving those things to the service
- Here we're giving OrderService access to RealtimeService

### Step 8: Test Your Implementation

**File:** Create `test-order-websocket.html`

```html
<!DOCTYPE html>
<html>
<head>
    <title>Order WebSocket Test</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
        }
        .status {
            padding: 10px;
            margin: 10px 0;
            border-radius: 5px;
        }
        .connected { background-color: #d4edda; color: #155724; }
        .disconnected { background-color: #f8d7da; color: #721c24; }
        .message {
            background-color: #f0f0f0;
            padding: 10px;
            margin: 5px 0;
            border-left: 4px solid #007bff;
        }
        button {
            padding: 10px 20px;
            margin: 5px;
            cursor: pointer;
        }
        #messages {
            max-height: 400px;
            overflow-y: auto;
            border: 1px solid #ddd;
            padding: 10px;
        }
    </style>
</head>
<body>
    <h1>🔌 Order WebSocket Tester</h1>

    <div id="status" class="status disconnected">
        Status: Disconnected
    </div>

    <div>
        <input type="text" id="userId" placeholder="User ID" value="user123">
        <button onclick="connect()">Connect</button>
        <button onclick="disconnect()">Disconnect</button>
    </div>

    <h3>Received Messages:</h3>
    <div id="messages"></div>

    <script>
        let ws = null;

        function connect() {
            const userId = document.getElementById('userId').value;
            const url = `ws://localhost:8080/api/v1/realtime/users/${userId}/orders/watch`;

            ws = new WebSocket(url);

            ws.onopen = function() {
                document.getElementById('status').className = 'status connected';
                document.getElementById('status').textContent = 'Status: Connected ✓';
                addMessage('Connected to server', 'info');
            };

            ws.onmessage = function(event) {
                const data = JSON.parse(event.data);
                addMessage('Received: ' + JSON.stringify(data, null, 2), 'success');

                // Show notification if order is ready
                if (data.type === 'order_ready') {
                    if (Notification.permission === 'granted') {
                        new Notification('Order Ready!', {
                            body: data.data.message
                        });
                    }
                }
            };

            ws.onerror = function(error) {
                addMessage('Error: ' + error, 'error');
            };

            ws.onclose = function() {
                document.getElementById('status').className = 'status disconnected';
                document.getElementById('status').textContent = 'Status: Disconnected ✗';
                addMessage('Disconnected from server', 'info');
            };
        }

        function disconnect() {
            if (ws) {
                ws.close();
            }
        }

        function addMessage(message, type) {
            const messagesDiv = document.getElementById('messages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message';
            messageDiv.textContent = new Date().toLocaleTimeString() + ' - ' + message;
            messagesDiv.insertBefore(messageDiv, messagesDiv.firstChild);
        }

        // Request notification permission
        if (Notification.permission === 'default') {
            Notification.requestPermission();
        }
    </script>
</body>
</html>
```

**How to Test:**

1. **Start your Go server:**
```bash
go run cmd/main.go
```

2. **Open the HTML file in a browser**

3. **Connect to WebSocket:**
   - Enter a user ID (e.g., "user123")
   - Click "Connect"
   - You should see "Connected to server" message

4. **Trigger an order update:**
   - Use your API to create an order for user123
   - Update the order status
   - Watch the WebSocket test page - you should see updates appear!

5. **Using curl to test:**
```bash
# Create an order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "items": [{"food_id": "food1", "quantity": 2}]
  }'

# Update order status (this should send WebSocket message)
curl -X PUT http://localhost:8080/api/v1/orders/order456/status \
  -H "Content-Type: application/json" \
  -d '{"status": "preparing"}'
```

---

## Practical Exercise

Now implement **Cart Updates** WebSocket yourself! Follow these steps:

### Exercise: Implement Real-time Shopping Cart Updates

**Goal:** Users should see their shopping cart update in real-time when items are added/removed.

**Step-by-Step Checklist:**

#### □ Step 1: Add Message Types
- [ ] Open `internal/models/websocket.go`
- [ ] Add constants: `MessageTypeCartAdd`, `MessageTypeCartRemove`, `MessageTypeCartClear`
- [ ] Create `CartUpdateData` struct with fields: `Item`, `Action`, `TotalItems`, `TotalPrice`

#### □ Step 2: Add Service Methods
- [ ] Open `internal/service/realtime_service.go`
- [ ] Create `SendCartAdd(userID, item, totalItems, totalPrice)` method
- [ ] Create `SendCartRemove(userID, item, totalItems, totalPrice)` method
- [ ] Create `SendCartClear(userID)` method

#### □ Step 3: Create Handler
- [ ] Open `internal/api/handlers/realtime_handler.go`
- [ ] Create `WatchCart(c *gin.Context)` method
- [ ] Validate user exists
- [ ] Upgrade to WebSocket
- [ ] Call `HandleConnection()`

#### □ Step 4: Register Route
- [ ] Open `internal/api/routes/routes.go`
- [ ] Add route: `realtime.GET("/users/:userId/cart/watch", realtimeHandler.WatchCart)`

#### □ Step 5: Integrate with Cart Service
- [ ] Open your cart service file
- [ ] Add `realtimeService *RealtimeService` field
- [ ] Add `SetRealtimeService()` method
- [ ] In `AddToCart()` method, call `SendCartAdd()` after successful add
- [ ] In `RemoveFromCart()` method, call `SendCartRemove()` after successful removal
- [ ] In `ClearCart()` method, call `SendCartClear()` after successful clear

#### □ Step 6: Wire Dependencies
- [ ] Open `internal/api/routes/routes.go`
- [ ] Add `deps.CartService.SetRealtimeService(realtimeService)`

#### □ Step 7: Test
- [ ] Build: `go build ./...`
- [ ] Run server: `go run cmd/main.go`
- [ ] Create test HTML file (similar to order test)
- [ ] Connect to WebSocket
- [ ] Add item to cart via API
- [ ] See update appear in real-time!

### Hints for the Exercise:

<details>
<summary>Click to see CartUpdateData structure hint</summary>

```go
type CartUpdateData struct {
    Item       *CartItem `json:"item"`        // The cart item
    Action     string    `json:"action"`      // "add", "remove", "clear"
    TotalItems int       `json:"total_items"` // Total number of items in cart
    TotalPrice float64   `json:"total_price"` // Total cart price
}
```
</details>

<details>
<summary>Click to see SendCartAdd method hint</summary>

```go
func (s *RealtimeService) SendCartAdd(userID string, item *models.CartItem, totalItems int, totalPrice float64) {
    updateData := models.CartUpdateData{
        Item:       item,
        Action:     "add",
        TotalItems: totalItems,
        TotalPrice: totalPrice,
    }

    message := models.WSMessage{
        Type:      models.MessageTypeCartAdd,
        Data:      updateData,
        Timestamp: time.Now(),
    }

    s.sendToUser(userID, message)
    log.Printf("Realtime: Sent cart add to %s, item: %s", userID, item.ID)
}
```
</details>

---

## Testing Your Implementation

### Manual Testing Checklist

#### □ Build Test
```bash
go build ./...
# Should complete with no errors
```

#### □ Start Server
```bash
go run cmd/main.go
# Server should start on port 8080
```

#### □ Test WebSocket Connection
1. Open test HTML file in browser
2. Enter user ID
3. Click "Connect"
4. Status should show "Connected ✓"

#### □ Test Real-time Updates
1. Keep WebSocket connection open
2. Use curl or Postman to update data:
```bash
# Update order status
curl -X PUT http://localhost:8080/api/v1/orders/ORDER_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "preparing"}'
```
3. Check test page - should show update message

#### □ Test Multiple Connections
1. Open test page in multiple browser tabs
2. Use same user ID in all tabs
3. Trigger an update
4. All tabs should receive the update

#### □ Test Disconnection
1. Click "Disconnect" button
2. Status should show "Disconnected ✗"
3. Trigger an update
4. Should not receive message (connection closed)

### Common Test Scenarios

#### Scenario 1: User Updates Profile
```
Setup:
1. Connect WebSocket: ws://localhost:8080/api/v1/realtime/users/user123/watch
2. Keep connection open

Trigger:
curl -X PATCH http://localhost:8080/api/v1/users/user123/first_name \
  -H "Content-Type: application/json" \
  -d '{"value": "NewName"}'

Expected Result:
{
  "type": "user_update",
  "data": {
    "user": {...},
    "changes": {"first_name": "NewName"}
  },
  "timestamp": "2025-10-31T..."
}
```

#### Scenario 2: Order Status Changes
```
Setup:
1. Connect: ws://localhost:8080/api/v1/realtime/users/user123/orders/watch
2. Create order: POST /api/v1/orders

Trigger:
curl -X PUT http://localhost:8080/api/v1/orders/order456/status \
  -d '{"status": "preparing"}'

Expected Result:
{
  "type": "order_update",
  "data": {
    "order": {...},
    "status": "preparing",
    "message": "Your order is being prepared"
  },
  "timestamp": "..."
}
```

---

## Common Pitfalls and Solutions

### Pitfall 1: "Connection Refused" Error

**Symptom:**
```
Error: WebSocket connection failed
net::ERR_CONNECTION_REFUSED
```

**Causes & Solutions:**
1. **Server not running**
   - Solution: Run `go run cmd/main.go`
   - Check: Server should print "Server running on :8080"

2. **Wrong port**
   - Solution: Check your config file for correct port
   - Common: Port 8080, 8000, or 3000

3. **Firewall blocking**
   - Solution: Allow port 8080 in firewall
   - Windows: Windows Defender Firewall settings
   - Linux: `sudo ufw allow 8080`

### Pitfall 2: "User Not Found" Error

**Symptom:**
```json
{"error": "User not found", "user_id": "user123"}
```

**Causes & Solutions:**
1. **User doesn't exist in database**
   - Solution: Create user first via `/api/v1/auth/register`
   - Or: Use an existing user ID from your database

2. **Wrong user ID**
   - Solution: Check database for correct user IDs
   - Query: `SELECT id FROM users;`

### Pitfall 3: Not Receiving Updates

**Symptom:** WebSocket connects but no messages appear when data changes.

**Causes & Solutions:**

1. **RealtimeService not injected**
   - Check: `routes.go` should have `deps.OrderService.SetRealtimeService(realtimeService)`
   - Solution: Add the injection line

2. **Service method not calling RealtimeService**
   - Check: `UpdateOrderStatus` should call `SendOrderUpdate()`
   - Solution: Add the call after database update

3. **Nil pointer check failing**
   - Check: `if s.realtimeService != nil` - might be nil
   - Solution: Make sure injection happens before any service calls

4. **Wrong user ID**
   - Check: WebSocket connected with userID "user123" but update is for "user456"
   - Solution: Use same user ID for connection and updates

**Debugging Steps:**
```go
// Add debug logs in your service:
func (s *orderService) UpdateOrderStatus(orderID string, status models.OrderStatus) error {
    // ... update code ...

    log.Printf("DEBUG: realtimeService is nil? %v", s.realtimeService == nil)

    if s.realtimeService != nil {
        log.Printf("DEBUG: Sending update to user %s", order.UserID)
        s.realtimeService.SendOrderUpdate(...)
    }

    return nil
}
```

### Pitfall 4: CORS Errors

**Symptom:**
```
Access to WebSocket has been blocked by CORS policy
```

**Solution:**
In `realtime_handler.go`, the upgrader has:
```go
var realtimeUpgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true  // This allows all origins - OK for development
    },
}
```

**For Production:**
```go
var realtimeUpgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        // Only allow your frontend domains
        return origin == "https://yourapp.com" ||
               origin == "https://www.yourapp.com"
    },
}
```

### Pitfall 5: Connection Drops Immediately

**Symptom:** Connection closes right after connecting.

**Causes & Solutions:**

1. **Authentication middleware blocking**
   - Check: Routes have `.Use(middleware.AuthMiddleware(...))`
   - Solution: Make sure you're sending auth token in request
   ```javascript
   const ws = new WebSocket('ws://localhost:8080/api/v1/realtime/users/user123/watch', {
       headers: {
           'Authorization': 'Bearer YOUR_TOKEN_HERE'
       }
   });
   ```

2. **Panic in HandleConnection**
   - Check: Server logs for panic messages
   - Solution: Add error handling in HandleConnection

### Pitfall 6: Messages Not JSON Formatted

**Symptom:** Receiving `[object Object]` or malformed data.

**Cause:** Data not being properly serialized.

**Solution:** Always use JSON tags and proper marshaling:
```go
type MyData struct {
    Name  string `json:"name"`   // ✓ Correct
    Value int    // ✗ Wrong - no json tag
}
```

### Pitfall 7: Memory Leak (Connections Not Closing)

**Symptom:** Server memory grows over time.

**Causes & Solutions:**

1. **Not removing connections**
   - Check: `RemoveConnection()` is called in `defer`
   - Solution:
   ```go
   func (s *RealtimeService) HandleConnection(userID string, conn *websocket.Conn) {
       s.AddConnection(userID, conn)
       defer s.RemoveConnection(userID)  // ✓ Always remove on exit
       // ... rest of code ...
   }
   ```

2. **Multiple connections per user**
   - Our code already handles this by replacing old connections
   - Check: `AddConnection()` closes existing connection before adding new one

---

## Troubleshooting Guide

### Problem: Build Fails

**Check 1: Missing imports**
```bash
Error: undefined: websocket

Solution:
import "github.com/gorilla/websocket"
```

**Check 2: Undefined types**
```bash
Error: undefined: RealtimeService

Solution:
- Make sure the type is exported (starts with capital letter)
- Check file is in correct package
- Verify import paths
```

**Check 3: Circular imports**
```bash
Error: import cycle not allowed

Solution:
- Don't import service package in models
- Don't import handlers package in services
- Keep dependencies flowing downward: handlers → services → repository
```

### Problem: Runtime Panics

**Check 1: Nil pointer dereference**
```bash
panic: runtime error: invalid memory address or nil pointer dereference

Solution:
Always check if pointer is nil before using:
if s.realtimeService != nil {
    s.realtimeService.SendUpdate(...)
}
```

**Check 2: Race conditions**
```bash
fatal error: concurrent map read and map write

Solution:
- Always use mutex locks when accessing shared maps
- Our RealtimeService already has proper locking
```

### Problem: WebSocket Connects But Then Disconnects

**Debug Steps:**

1. **Check server logs:**
```bash
go run cmd/main.go 2>&1 | tee server.log
```

2. **Add debug logging:**
```go
func (s *RealtimeService) HandleConnection(userID string, conn *websocket.Conn) {
    log.Printf("DEBUG: HandleConnection called for user: %s", userID)

    s.AddConnection(userID, conn)
    log.Printf("DEBUG: Connection added successfully")

    defer func() {
        log.Printf("DEBUG: Connection closing for user: %s", userID)
        s.RemoveConnection(userID)
    }()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("DEBUG: Read error for user %s: %v", userID, err)
            break
        }
        log.Printf("DEBUG: Received message from user: %s", userID)
    }
}
```

3. **Check authentication:**
```bash
# Make sure auth middleware isn't rejecting connections
# Temporarily remove auth middleware to test:
realtime := v1.Group("/realtime")
// realtime.Use(middleware.AuthMiddleware(...))  // Comment out for testing
```

### Problem: Updates Work Sometimes But Not Others

**Likely Cause:** Race condition or timing issue.

**Debug Steps:**

1. **Add request ID logging:**
```go
func (s *orderService) UpdateOrderStatus(orderID string, status models.OrderStatus) error {
    requestID := uuid.New().String()
    log.Printf("[%s] Starting update for order: %s", requestID, orderID)

    err := s.orderRepo.UpdateStatus(orderID, status)
    log.Printf("[%s] Database update result: %v", requestID, err)

    if s.realtimeService != nil {
        log.Printf("[%s] Sending WebSocket update", requestID)
        s.realtimeService.SendOrderUpdate(...)
        log.Printf("[%s] WebSocket update sent", requestID)
    } else {
        log.Printf("[%s] RealtimeService is nil!", requestID)
    }

    return nil
}
```

2. **Check connection state:**
```go
func (s *RealtimeService) sendToUser(userID string, message models.WSMessage) {
    s.mutex.RLock()
    conn, exists := s.connections[userID]
    s.mutex.RUnlock()

    log.Printf("DEBUG: Sending to user %s - connection exists: %v", userID, exists)

    if !exists {
        log.Printf("WARNING: User %s not connected", userID)
        return
    }

    // ... rest of code ...
}
```

### Quick Reference: Log Analysis

**Good Logs:**
```
Realtime: User user123 connected
DEBUG: HandleConnection called for user: user123
DEBUG: Connection added successfully
[abc-123] Starting update for order: order456
[abc-123] Database update result: <nil>
[abc-123] Sending WebSocket update
DEBUG: Sending to user user123 - connection exists: true
Realtime: Sent order update to user123, order: order456, status: preparing
```

**Bad Logs (Problem):**
```
Realtime: User user123 connected
DEBUG: Connection closing for user: user123  ← Closed too early!
[abc-123] RealtimeService is nil!  ← Not injected!
DEBUG: Sending to user user123 - connection exists: false  ← Not connected!
```

---

## Summary Checklist

Before you finish, verify:

### Implementation Checklist

- [ ] Message types added to `websocket.go`
- [ ] Data structures created with JSON tags
- [ ] Service methods implemented in `realtime_service.go`
- [ ] Handler methods created in `realtime_handler.go`
- [ ] Routes registered in `routes.go`
- [ ] Service has `realtimeService` field
- [ ] Service has `SetRealtimeService()` method
- [ ] Service integration completed (calls SendUpdate methods)
- [ ] Dependency injection added in `routes.go`
- [ ] Build succeeds: `go build ./...`

### Testing Checklist

- [ ] Server starts without errors
- [ ] WebSocket connects successfully
- [ ] Updates trigger real-time messages
- [ ] Multiple connections work
- [ ] Disconnection handled cleanly
- [ ] Messages have correct JSON structure
- [ ] Log messages appear for debugging

### Production Readiness Checklist

- [ ] CORS properly configured (not just `return true`)
- [ ] Authentication working on WebSocket endpoints
- [ ] Error handling for all edge cases
- [ ] Connection limits implemented (optional)
- [ ] Debug logs removed or made conditional
- [ ] Load testing performed
- [ ] Documentation updated

---

## Need Help?

### Quick Debugging Commands

```bash
# Check if server is running
curl http://localhost:8080/api/v1/health

# Check WebSocket endpoint (should upgrade)
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: test" \
  http://localhost:8080/api/v1/realtime/users/user123/watch

# View server logs in real-time
go run cmd/main.go 2>&1 | grep -i "realtime\|websocket"

# Check active connections
curl http://localhost:8080/api/v1/realtime/stats
```

### Common Questions

**Q: Can multiple devices connect with same user ID?**
A: Yes! Our implementation supports multiple connections per user. All devices receive updates.

**Q: What happens if user is offline when update occurs?**
A: Update is lost. WebSocket is for real-time only. For guaranteed delivery, use push notifications.

**Q: How many concurrent connections can the server handle?**
A: Depends on server resources. Generally thousands. Monitor with `/realtime/stats` endpoint.

**Q: Do I need to send heartbeat/ping messages?**
A: No, our implementation automatically handles this with ReadMessage loop.

**Q: Can I send messages from client to server?**
A: Current implementation only receives to keep connection alive. To add client→server:
```go
for {
    _, message, err := conn.ReadMessage()
    if err != nil {
        break
    }
    // Handle client message here
    handleClientMessage(userID, message)
}
```

---

**Congratulations!** 🎉

You now understand how WebSocket endpoints work in this project. Use this guide whenever you need to add real-time features!

**Remember:**
1. Start with models (message types and data)
2. Add service methods (business logic)
3. Create handlers (entry points)
4. Register routes (URL mapping)
5. Integrate with existing services
6. Test thoroughly

Good luck with your implementation! 🚀
