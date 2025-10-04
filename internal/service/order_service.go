package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type OrderService interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	GetUserOrders(userID string, limit, offset int) ([]models.Order, error)
	GetOrderByID(orderID string) (*models.Order, error)
	UpdateOrderStatus(orderID string, status models.OrderStatus) error
	CancelOrder(orderID string) error
	TrackOrder(orderID string) (*models.Order, error)
}

type orderService struct {
	orderRepo      repository.OrderRepository
	userRepo       repository.UserRepository
	restaurantRepo repository.RestaurantRepository
	foodRepo       repository.FoodRepository
}

func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository, restaurantRepo repository.RestaurantRepository, foodRepo repository.FoodRepository) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		userRepo:       userRepo,
		restaurantRepo: restaurantRepo,
		foodRepo:       foodRepo,
	}
}

func (s *orderService) CreateOrder(order *models.Order) (*models.Order, error) {
	if order == nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order is required", nil)
	}

	// Validate required fields
	if strings.TrimSpace(order.UserID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(order.RestaurantID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}
	if len(order.Items) == 0 {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order items are required", nil)
	}
	if strings.TrimSpace(order.DeliveryAddress) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Delivery address is required", nil)
	}
	if strings.TrimSpace(order.PaymentMethod) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Payment method is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(order.UserID)
	if err != nil {
		return nil, err
	}

	// Validate restaurant exists
	restaurant, err := s.restaurantRepo.GetByID(order.RestaurantID)
	if err != nil {
		return nil, err
	}
	order.RestaurantName = restaurant.Name

	// Validate order items and calculate totals
	var subtotal float64
	for i, item := range order.Items {
		if strings.TrimSpace(item.FoodID) == "" {
			return nil, errors.NewHTTPError(http.StatusBadRequest, "Food ID is required for all items", nil)
		}
		if item.Quantity <= 0 {
			return nil, errors.NewHTTPError(http.StatusBadRequest, "Quantity must be greater than 0", nil)
		}

		// Validate food exists and is available
		food, err := s.foodRepo.GetByID(item.FoodID)
		if err != nil {
			return nil, err
		}
		if !food.IsAvailable {
			return nil, errors.NewHTTPError(http.StatusBadRequest, "Food item is not available: "+food.Name, nil)
		}
		if food.RestaurantID != order.RestaurantID {
			return nil, errors.NewHTTPError(http.StatusBadRequest, "All items must be from the same restaurant", nil)
		}

		// Update item details
		order.Items[i].FoodName = food.Name
		order.Items[i].Price = food.Price
		order.Items[i].Total = food.Price * float64(item.Quantity)
		subtotal += order.Items[i].Total
	}

	// Set order totals
	order.Subtotal = subtotal
	if order.DeliveryFee < 0 {
		order.DeliveryFee = restaurant.DeliveryFee
	}
	if order.Tax <= 0 {
		order.Tax = subtotal * 0.08 // 8% tax rate
	}
	order.Total = order.Subtotal + order.DeliveryFee + order.Tax
	order.Status = models.OrderStatusPending

	// Create order
	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID string, limit, offset int) ([]models.Order, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.orderRepo.GetByUserID(userID, limit, offset)
}

func (s *orderService) GetOrderByID(orderID string) (*models.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
	}

	return s.orderRepo.GetByID(orderID)
}

func (s *orderService) UpdateOrderStatus(orderID string, status models.OrderStatus) error {
	if strings.TrimSpace(orderID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
	}

	// Validate status
	validStatuses := map[models.OrderStatus]bool{
		models.OrderStatusPending:   true,
		models.OrderStatusConfirmed: true,
		models.OrderStatusPreparing: true,
		models.OrderStatusOnTheWay:  true,
		models.OrderStatusDelivered: true,
		models.OrderStatusCancelled: true,
	}
	if !validStatuses[status] {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid order status", nil)
	}

	// Validate order exists
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	// Validate status transition
	if order.Status == models.OrderStatusDelivered || order.Status == models.OrderStatusCancelled {
		return errors.NewHTTPError(http.StatusBadRequest, "Cannot update status of completed order", nil)
	}

	return s.orderRepo.UpdateStatus(orderID, status)
}

func (s *orderService) CancelOrder(orderID string) error {
	if strings.TrimSpace(orderID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
	}

	// Validate order exists
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	// Check if order can be cancelled
	if order.Status == models.OrderStatusDelivered {
		return errors.NewHTTPError(http.StatusBadRequest, "Cannot cancel delivered order", nil)
	}
	if order.Status == models.OrderStatusCancelled {
		return errors.NewHTTPError(http.StatusBadRequest, "Order is already cancelled", nil)
	}

	return s.orderRepo.UpdateStatus(orderID, models.OrderStatusCancelled)
}

func (s *orderService) TrackOrder(orderID string) (*models.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
	}

	return s.orderRepo.GetByID(orderID)
}
