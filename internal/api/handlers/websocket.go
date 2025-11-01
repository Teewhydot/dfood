package handlers

import (
	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	wsService   *service.WebSocketService
	userService service.UserService
}

func NewWebSocketHandler(wsService *service.WebSocketService, userService service.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		wsService:   wsService,
		userService: userService,
	}
}

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, check the origin properly
		return true
	},
}

// WatchUserProfile handles WebSocket connections for watching user profile updates
// Endpoint: GET /api/v1/users/:userId/watch
func (h *WebSocketHandler) WatchUserProfile(c *gin.Context) {
	userID := c.Param("userId")

	// Verify user exists before upgrading connection
	user, err := h.userService.GetByID(userID)
	if err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, err
			},
			"verifying user for WebSocket connection",
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

	// Send initial profile data immediately after connection
	initialMessage := models.WSMessage{
		Type: models.MessageTypeUserUpdate,
		Data: models.UserUpdateData{
			User:    user,
			Changes: map[string]interface{}{"initial": true},
		},
		Timestamp: time.Now(),
	}

	if err := conn.WriteJSON(initialMessage); err != nil {
		conn.Close()
		return
	}

	// Handle the connection (this blocks until connection closes)
	// The service will manage sending updates when profile changes
	h.wsService.HandleConnection(userID, conn)
}

// GetConnectionStats returns WebSocket connection statistics
// Endpoint: GET /api/v1/websocket/stats
func (h *WebSocketHandler) GetConnectionStats(c *gin.Context) {
	stats := gin.H{
		"total_connections": h.wsService.GetConnectionCount(),
		"timestamp":         time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
