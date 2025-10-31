package handlers

import (
	"dfood/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RealtimeHandler handles all WebSocket connections
type RealtimeHandler struct {
	realtimeService     *service.RealtimeService
	userService         service.UserService
	addressService      service.AddressService
	favoritesService    service.FavoritesService
	notificationService service.NotificationService
}

// NewRealtimeHandler creates a new realtime handler
func NewRealtimeHandler(
	realtimeService *service.RealtimeService,
	userService service.UserService,
	addressService service.AddressService,
	favoritesService service.FavoritesService,
	notificationService service.NotificationService,
) *RealtimeHandler {
	return &RealtimeHandler{
		realtimeService:     realtimeService,
		userService:         userService,
		addressService:      addressService,
		favoritesService:    favoritesService,
		notificationService: notificationService,
	}
}

// WebSocket upgrader
var realtimeUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure properly for production
	},
}

// User Profile WebSocket Endpoints

// WatchUserProfile handles user profile updates WebSocket
func (h *RealtimeHandler) WatchUserProfile(c *gin.Context) {
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

// Address WebSocket Endpoints

// WatchAddresses handles address updates WebSocket
func (h *RealtimeHandler) WatchAddresses(c *gin.Context) {
	userID := c.Param("userId")

	// Validate user exists and has addresses
	_, err := h.addressService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found or no addresses",
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

// Favorites WebSocket Endpoints

// WatchFavoriteFoods handles favorite foods updates WebSocket
func (h *RealtimeHandler) WatchFavoriteFoods(c *gin.Context) {
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

// WatchFavoriteRestaurants handles favorite restaurants updates WebSocket
func (h *RealtimeHandler) WatchFavoriteRestaurants(c *gin.Context) {
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

// Notification WebSocket Endpoints

// WatchNotifications handles notifications updates WebSocket
func (h *RealtimeHandler) WatchNotifications(c *gin.Context) {
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

// Statistics Endpoint

// GetConnectionStats returns WebSocket connection statistics
func (h *RealtimeHandler) GetConnectionStats(c *gin.Context) {
	count := h.realtimeService.GetConnectionCount()

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"connected_users": count,
		"message":         "WebSocket connection statistics",
	})
}
