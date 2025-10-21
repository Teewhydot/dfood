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
	addressService      service.AddressService
	favoritesService    service.FavoritesService
}

func NewWebSocketHandler(
	wsService service.WebSocketService,
	userService service.UserService,
	orderService service.OrderService,
	notificationService service.NotificationService,
	restaurantService service.RestaurantService,
	addressService service.AddressService,
	favoritesService service.FavoritesService,
) *WebSocketHandler {
	return &WebSocketHandler{
		wsService:           wsService,
		userService:         userService,
		orderService:        orderService,
		notificationService: notificationService,
		restaurantService:   restaurantService,
		addressService:      addressService,
		favoritesService:    favoritesService,
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

// Address WebSocket endpoints
func (h *WebSocketHandler) WatchAddresses(c *gin.Context) {
	userID := c.Param("userId")
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeAddress, userID, userID)
}

// Favorites WebSocket endpoints
func (h *WebSocketHandler) WatchFavoriteFoods(c *gin.Context) {
	userID := c.Param("userId")
	resourceKey := userID + "_foods"
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeFavorites, userID, resourceKey)
}

func (h *WebSocketHandler) WatchFavoriteRestaurants(c *gin.Context) {
	userID := c.Param("userId")
	resourceKey := userID + "_restaurants"
	h.handleWebSocketUpgrade(c, models.WSConnectionTypeFavorites, userID, resourceKey)
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

	case models.WSConnectionTypeAddress:
		// Verify user exists and can access their own addresses
		_, err := h.userService.GetByID(userID)
		return err

	case models.WSConnectionTypeFavorites:
		// Verify user exists and can access their own favorites
		_, err := h.userService.GetByID(userID)
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
