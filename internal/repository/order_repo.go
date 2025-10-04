package repository

import (
	"errors"
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{
		db: database.DB,
	}
}

func (r *orderRepository) Create(order *models.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to create order", err)
	}
	return nil
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "Order not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch order", err)
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID string, limit, offset int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&orders).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user orders", err)
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(id string, status models.OrderStatus) error {
	err := r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update order status", err)
	}
	return nil
}

func (r *orderRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.Order{}).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to delete order", err)
	}
	return nil
}
