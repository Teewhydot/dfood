# Step-by-Step WebSocket Endpoint Implementation

## 🎯 Example: Restaurant Menu Updates WebSocket

We'll implement a WebSocket endpoint that notifies clients when a restaurant's menu changes in real-time.

**Endpoint:** `/ws/restaurants/:restaurantId/menu/watch`
**Purpose:** Notify clients when menu items are added, removed, or updated

---

## Step 1: Define Message Types and Payloads

**Purpose:** Establish the contract for what messages will be sent over WebSocket

### 1.1 Add Message Types (`internal/models/websocket.go`)

```go
// Add to existing WSMessageType constants
const (
    // ... existing message types ...
    
    // Menu related messages - NEW
    WSMessageTypeMenuItemAdd    WSMessageType = "menu_item_add"    // New item added to menu
    WSMessageTypeMenuItemUpdate WSMessageType = "menu_item_update" // Existing item modified
    WSMessageTypeMenuItemRemove WSMessageType = "menu_item_remove" // Item removed from menu
    WSMessageTypeMenuUpdate     WSMessageType = "menu_update"     // General menu changes
)
```

**Why this step:** Clients need to know what type of message they're receiving so they can handle it appropriately (e.g., add item to UI vs remove item from UI).

### 1.2 Create Payload Structure

```go
// Add to existing payload types
type MenuUpdatePayload struct {
    RestaurantID string      `json:"restaurant_id"`           // Which restaurant's menu changed
    MenuItem     *Food       `json:"menu_item,omitempty"`     // The food item (for add/update)
    ItemID       string      `json:"item_id,omitempty"`       // Item ID (for remove operations)
    Action       string      `json:"action"`                  // "add", "update", "remove"
    Changes      map[string]interface{} `json:"changes,omitempty"` // What changed (for updates)
    UpdatedBy    string      `json:"updated_by"`              // Who made the change
}
```

**Why this step:** Defines the exact data structure clients will receive, ensuring consistency and type safety.

---

## Step 2: Add Connection Type (if needed)

**Purpose:** Categorize this type of WebSocket connection for proper routing

### 2.1 Add Connection Type (`internal/models/websocket.go`)

```go
// Add to existing WSConnectionType constants (if not already present)
const (
    // ... existing connection types ...
    WSConnectionTypeMenu WSConnectionType = "menu" // Menu updates
)
```

**Why this step:** Allows us to group connections by what they're watching and route messages efficiently.

---

## Step 3: Add Broadcast Method to WebSocket Service

**Purpose:** Create a convenient method for broadcasting menu updates

### 3.1 Add Interface Method (`internal/service/websocket_interface.go`)

```go
// Add to WebSocketService interface
type WebSocketService interface {
    // ... existing methods ...
    
    // BroadcastMenuUpdate notifies when restaurant menu changes
    BroadcastMenuUpdate(restaurantID string, menuItem *models.Food, action string, changes map[string]interface{})
}
```

### 3.2 Implement Method (`internal/service/websocket_service.go`)

```go
// Add to webSocketService implementation
func (s *webSocketService) BroadcastMenuUpdate(restaurantID string, menuItem *models.Food, action string, changes map[string]interface{}) {
    // Create the payload with menu change details
    payload := models.MenuUpdatePayload{
        RestaurantID: restaurantID,
        MenuItem:     menuItem,
        Action:       action,           // "add", "update", "remove"
        Changes:      changes,          // What specifically changed
        UpdatedBy:    "restaurant",     // Who made the change
    }
    
    // If it's a remove action, we only need the item ID
    if action == "remove" && menuItem != nil {
        payload.ItemID = menuItem.ID
        payload.MenuItem = nil // Don't send full object for removals
    }
    
    // Convert payload to JSON for transmission
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshaling menu update: %v", err)
        return
    }
    
    // Determine message type based on action
    var msgType models.WSMessageType
    switch action {
    case "add":
        msgType = models.WSMessageTypeMenuItemAdd
    case "update":
        msgType = models.WSMessageTypeMenuItemUpdate
    case "remove":
        msgType = models.WSMessageTypeMenuItemRemove
    default:
        msgType = models.WSMessageTypeMenuUpdate
    }
    
    // Create the WebSocket message
    msg := models.WSMessage{
        Type:      msgType,
        Data:      data,
        Timestamp: time.Now(),
    }
    
    // Send to all connections watching this restaurant's menu
    s.BroadcastToResource(models.WSConnectionTypeMenu, restaurantID, msg)
}
```

**Why this step:** Provides a clean, reusable method that other services can call to broadcast menu updates. Handles message formatting and routing automatically.

---

## Step 4: Create WebSocket Handler

**Purpose:** Handle HTTP → WebSocket upgrade and connection management

### 4.1 Add Handler Method (`internal/api/handlers/websocket.go`)

```go
// Add to WebSocketHandler
func (h *WebSocketHandler) WatchRestaurantMenu(c *gin.Context) {
    restaurantID := c.Param("restaurantId")
    userID := c.GetString("user_id") // From JWT middleware (optional for public menus)
    
    // Verify restaurant exists
    _, err := h.restaurantService.GetByID(restaurantID)
    if err != nil {
        result := errors.HandleError(
            func() (interface{}, error) {
                return nil, err
            },
            "verifying restaurant for menu WebSocket connection",
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
    
    // Handle the connection (this blocks until connection closes)
    h.wsService.HandleConnection(models.WSConnectionTypeMenu, userID, restaurantID, conn)
}
```

**Why this step:** 
- **Validates** the restaurant exists before allowing connection
- **Upgrades** HTTP to WebSocket protocol
- **Delegates** connection management to the WebSocket service
- **Handles errors** gracefully with proper HTTP responses

---

## Step 5: Add Route

**Purpose:** Make the WebSocket endpoint accessible via URL

### 5.1 Add Route (`internal/api/routes/routes.go`)

```go
// Add to WebSocket routes section
ws := v1.Group("/ws")
ws.Use(middleware.AuthMiddleware(deps.UserRepository)) // Optional: remove if public
{
    // ... existing routes ...
    
    // Restaurant menu WebSocket endpoint
    ws.GET("/restaurants/:restaurantId/menu/watch", wsHandler.WatchRestaurantMenu)
}
```

**Why this step:** Makes the endpoint accessible at `/ws/restaurants/123/menu/watch`. The middleware is optional depending on whether menu watching should be public or require authentication.

---

## Step 6: Integrate with Business Logic

**Purpose:** Automatically broadcast WebSocket updates when menu data changes

### 6.1 Update Food Service (`internal/service/food_service.go`)

```go
// Add WebSocket service to FoodService interface
type FoodService interface {
    // ... existing methods ...
    SetWebSocketService(wsService WebSocketService)
    GetWebSocketService() WebSocketService
}

// Add to foodService struct
type foodService struct {
    // ... existing fields ...
    wsService WebSocketService
}

// Add WebSocket service methods
func (s *foodService) SetWebSocketService(wsService WebSocketService) {
    s.wsService = wsService
}

func (s *foodService) GetWebSocketService() WebSocketService {
    return s.wsService
}
```

### 6.2 Update Food CRUD Methods

```go
// Example: Update the CreateFood method
func (s *foodService) CreateFood(food *models.Food) error {
    // 1. Save to database
    err := s.foodRepo.Create(food)
    if err != nil {
        return err
    }
    
    // 2. Broadcast WebSocket update
    if s.wsService != nil {
        s.wsService.BroadcastMenuUpdate(
            food.RestaurantID,  // Which restaurant
            food,               // The new menu item
            "add",              // Action type
            nil,                // No changes for new items
        )
    }
    
    return nil
}

// Example: Update the UpdateFood method
func (s *foodService) UpdateFood(foodID string, updates map[string]interface{}) error {
    // 1. Update database
    err := s.foodRepo.Update(foodID, updates)
    if err != nil {
        return err
    }
    
    // 2. Broadcast WebSocket update
    if s.wsService != nil {
        // Get updated food item
        food, _ := s.foodRepo.GetByID(foodID)
        if food != nil {
            s.wsService.BroadcastMenuUpdate(
                food.RestaurantID,  // Which restaurant
                food,               // Updated menu item
                "update",           // Action type
                updates,            // What changed
            )
        }
    }
    
    return nil
}

// Example: Update the DeleteFood method
func (s *foodService) DeleteFood(foodID string) error {
    // 1. Get food item before deletion (for WebSocket broadcast)
    food, err := s.foodRepo.GetByID(foodID)
    if err != nil {
        return err
    }
    
    // 2. Delete from database
    err = s.foodRepo.Delete(foodID)
    if err != nil {
        return err
    }
    
    // 3. Broadcast WebSocket update
    if s.wsService != nil {
        s.wsService.BroadcastMenuUpdate(
            food.RestaurantID,  // Which restaurant
            food,               // Deleted menu item (for ID reference)
            "remove",           // Action type
            nil,                // No changes for deletions
        )
    }
    
    return nil
}
```

**Why this step:** 
- **Automatic updates**: WebSocket messages are sent whenever menu data changes via existing API endpoints
- **No code duplication**: Business logic remains unchanged, WebSocket is just an additional notification layer
- **Consistency**: All menu changes trigger WebSocket updates, regardless of how they're made

---

## Step 7: Initialize WebSocket Service in Main

**Purpose:** Wire up the WebSocket service with the food service

### 7.1 Update Dependencies (`cmd/main.go`)

```go
func main() {
    // ... existing initialization ...
    
    // Initialize services
    foodService := service.NewFoodService(foodRepo)
    // ... other services ...
    
    // Initialize WebSocket service
    wsService := service.NewWebSocketService(
        userRepo,
        orderRepo,
        notificationRepo,
        restaurantRepo,
        addressRepo,
        favoritesRepo,
    )
    
    // Set WebSocket service in services that need it
    foodService.SetWebSocketService(wsService)
    // ... other service integrations ...
    
    deps := &routes.Dependencies{
        // ... existing dependencies ...
        FoodService:      foodService,
        WebSocketService: wsService,
    }
    
    // ... rest of main ...
}
```

**Why this step:** Connects all the pieces together so the food service can broadcast WebSocket updates when menu items change.

---

## Step 8: Test the Implementation

**Purpose:** Verify the WebSocket endpoint works correctly

### 8.1 Test WebSocket Connection

```bash
# Using wscat (install with: npm install -g wscat)
wscat -c ws://localhost:8080/ws/restaurants/restaurant-123/menu/watch
```

**Expected response:**
```json
{"type": "connected", "timestamp": "2024-01-15T10:30:00Z", "connection_id": "abc-123"}
```

### 8.2 Test Menu Updates

```bash
# Add a new menu item via REST API
curl -X POST http://localhost:8080/api/v1/foods \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "New Pizza",
    "price": 15.99,
    "restaurant_id": "restaurant-123"
  }'
```

**Expected WebSocket message:**
```json
{
  "type": "menu_item_add",
  "data": {
    "restaurant_id": "restaurant-123",
    "menu_item": {
      "id": "food-456",
      "name": "New Pizza",
      "price": 15.99
    },
    "action": "add",
    "updated_by": "restaurant"
  },
  "timestamp": "2024-01-15T10:31:00Z"
}
```

### 8.3 Test Connection Statistics

```bash
# Check WebSocket connection stats
curl -X GET http://localhost:8080/ws/stats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Expected response:**
```json
{
  "success": true,
  "data": {
    "total_connections": 1,
    "connections_by_type": {
      "menu": 1
    },
    "timestamp": "2024-01-15T10:32:00Z"
  }
}
```

---

## 🎯 Summary of Each Step's Purpose

| Step | Purpose | What It Achieves |
|------|---------|------------------|
| **1. Define Messages** | Establish communication contract | Clients know what messages to expect |
| **2. Add Connection Type** | Categorize connections | Efficient message routing |
| **3. Add Broadcast Method** | Create reusable broadcasting | Clean API for sending updates |
| **4. Create Handler** | Handle WebSocket upgrade | Convert HTTP to WebSocket |
| **5. Add Route** | Make endpoint accessible | Clients can connect via URL |
| **6. Integrate Business Logic** | Automatic notifications | Updates sent when data changes |
| **7. Initialize Services** | Wire everything together | All components work as a system |
| **8. Test Implementation** | Verify functionality | Ensure everything works correctly |

## 🔄 The Complete Flow

```
1. Client connects to /ws/restaurants/123/menu/watch
2. Server upgrades HTTP → WebSocket
3. Server registers connection in memory
4. Restaurant adds new menu item via REST API
5. FoodService.CreateFood() saves to database
6. FoodService.CreateFood() calls wsService.BroadcastMenuUpdate()
7. WebSocket service finds all connections watching restaurant 123's menu
8. WebSocket service sends "menu_item_add" message to those connections
9. Client receives message and updates UI with new menu item
```

This pattern can be applied to any real-time feature: order tracking, user presence, chat messages, live notifications, etc. The key is following these steps systematically to ensure all pieces work together correctly.