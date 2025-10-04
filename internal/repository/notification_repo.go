package repository

import (
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{
		db: database.DB,
	}
}

func (r *notificationRepository) GetByUserID(userID string, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&notifications).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch notifications", err)
	}
	return notifications, nil
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	if err := r.db.Create(notification).Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to create notification", err)
	}
	return nil
}

func (r *notificationRepository) MarkAsRead(id string) error {
	err := r.db.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to mark notification as read", err)
	}
	return nil
}

func (r *notificationRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.Notification{}).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to delete notification", err)
	}
	return nil
}
