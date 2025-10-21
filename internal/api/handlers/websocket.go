package handlers

import (
	"dfood/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections
// Like a controller in Flutter - handles the presentation layer
type WebSocketHandler struct {
	wsService   *service.WebSocketService // WebSocket service for managing connections
	userService service.UserService       // User service for validation
}

// NewWebSocketHandler creates a new WebSocket handler
// Like dependency injection in Flutter
func NewWebSocketHandler(wsService *service.WebSocketService, userService service.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		wsService:   wsService,
		userService: userService,
	}
}

// WebSocket upgrader configuration
// This converts HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, you should validate the origin properly
		return true
	},
}

// WatchUserProfile handles the WebSocket endpoint for user profile updates
// Endpoint: GET /api/v1/users/:userId/watch
func (h *WebSocketHandler) WatchUserProfile(c *gin.Context) {
	// Extract user ID from URL parameter
	userID := c.Param("userId")

	// Validate that the user exists before allowing WebSocket connection
	_, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"user_id": userID,
		})
		return
	}

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to upgrade to WebSocket",
			"details": err.Error(),
		})
		return
	}

	// Handle the WebSocket connection
	// This call blocks until the connection closes
	h.wsService.HandleConnection(userID, conn)
}

// GetConnectionStats returns WebSocket connection statistics
// Useful for monitoring how many users are connected
func (h *WebSocketHandler) GetConnectionStats(c *gin.Context) {
	count := h.wsService.GetConnectionCount()

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"connected_users": count,
		"message":         "WebSocket connection statistics",
	})
}
