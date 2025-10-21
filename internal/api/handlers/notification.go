package handlers

import (
	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// Notification Management
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	// TODO: Implement get user notifications (limit 50)
	c.JSON(200, gin.H{"message": "Get user notifications - TODO"})
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	// TODO: Implement send notification to user
	c.JSON(200, gin.H{"message": "Send notification - TODO"})
}

func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	// TODO: Implement mark notification as read
	c.JSON(200, gin.H{"message": "Mark notification as read - TODO"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	// TODO: Implement delete notification
	c.JSON(200, gin.H{"message": "Delete notification - TODO"})
}

func (h *NotificationHandler) GetNotificationsStream(c *gin.Context) {
	userID := c.Param("userId")

	// Verify user exists
	_, err := h.notificationService.GetByUserID(userID, 1, 0)
	if err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, err
			},
			"verifying user notifications for WebSocket connection",
		)
		result.RespondWithJSON(c)
		return
	}

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Configure properly for production
		},
	}

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

	// Get WebSocket service from notification service
	if wsService := h.notificationService.GetWebSocketService(); wsService != nil {
		wsService.HandleConnection(models.WSConnectionTypeNotification, userID, userID, conn)
	} else {
		conn.Close()
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusServiceUnavailable, "WebSocket service not available", nil)
			},
			"getting WebSocket service",
		)
		result.RespondWithJSON(c)
	}
}

// Push Notifications
func (h *NotificationHandler) UpdateFCMToken(c *gin.Context) {
	// TODO: Implement update FCM token
	c.JSON(200, gin.H{"message": "Update FCM token - TODO"})
}

func (h *NotificationHandler) GetFCMToken(c *gin.Context) {
	// TODO: Implement get FCM token
	c.JSON(200, gin.H{"message": "Get FCM token - TODO"})
}

func (h *NotificationHandler) SendPushNotification(c *gin.Context) {
	// TODO: Implement send push notification to user
	c.JSON(200, gin.H{"message": "Send push notification - TODO"})
}
