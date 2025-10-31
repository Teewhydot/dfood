package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type NotificationService interface {
	GetByUserID(userID string, limit, offset int) ([]models.Notification, error)
	GetUserNotifications(userID string, limit, offset int) ([]models.Notification, error)
	SendNotification(notification *models.Notification) error
	MarkNotificationAsRead(notificationID string) error
	DeleteNotification(notificationID string) error
	SendPushNotification(userID, title, body string, data map[string]interface{}) error
	SetWebSocketService(wsService WebSocketService)
	GetWebSocketService() WebSocketService
	SetRealtimeService(realtimeService *RealtimeService)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	wsService        WebSocketService
	realtimeService  *RealtimeService
}

func NewNotificationService(notificationRepo repository.NotificationRepository, userRepo repository.UserRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

// WebSocket service methods
func (s *notificationService) SetWebSocketService(wsService WebSocketService) {
	s.wsService = wsService
}

func (s *notificationService) GetWebSocketService() WebSocketService {
	return s.wsService
}

func (s *notificationService) SetRealtimeService(realtimeService *RealtimeService) {
	s.realtimeService = realtimeService
}

func (s *notificationService) GetByUserID(userID string, limit, offset int) ([]models.Notification, error) {
	return s.GetUserNotifications(userID, limit, offset)
}

func (s *notificationService) GetUserNotifications(userID string, limit, offset int) ([]models.Notification, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.notificationRepo.GetByUserID(userID, limit, offset)
}

func (s *notificationService) SendNotification(notification *models.Notification) error {
	if notification == nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Notification is required", nil)
	}

	// Validate required fields
	if strings.TrimSpace(notification.UserID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(notification.Title) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Title is required", nil)
	}
	if strings.TrimSpace(notification.Body) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Body is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(notification.UserID)
	if err != nil {
		return err
	}

	// Validate notification type
	validTypes := map[string]bool{"order": true, "promotion": true, "system": true, "chat": true}
	if !validTypes[notification.Type] {
		notification.Type = "system" // Default to system
	}

	err = s.notificationRepo.Create(notification)
	if err != nil {
		return err
	}

	// Send realtime update
	if s.realtimeService != nil {
		s.realtimeService.SendNotificationNew(notification.UserID, notification)
	}

	return nil
}

func (s *notificationService) MarkNotificationAsRead(notificationID string) error {
	if strings.TrimSpace(notificationID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Notification ID is required", nil)
	}

	// Mark notification as read
	// Note: Realtime updates for mark as read are not implemented yet
	// as the repository doesn't provide GetByID method
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *notificationService) DeleteNotification(notificationID string) error {
	if strings.TrimSpace(notificationID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Notification ID is required", nil)
	}

	// Delete notification
	// Note: Realtime updates for delete are not implemented yet
	// as the repository doesn't provide GetByID method
	return s.notificationRepo.Delete(notificationID)
}

func (s *notificationService) SendPushNotification(userID, title, body string, data map[string]interface{}) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(title) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Title is required", nil)
	}
	if strings.TrimSpace(body) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Body is required", nil)
	}

	// Validate user exists and has FCM token
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user.FCMToken == nil || strings.TrimSpace(*user.FCMToken) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User does not have FCM token", nil)
	}

	// TODO: Implement actual FCM push notification logic
	// This would integrate with Firebase Cloud Messaging or similar service
	return errors.NewHTTPError(http.StatusNotImplemented, "Push notification sending not implemented", nil)
}
