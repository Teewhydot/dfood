package service

import (
	"net/http"
	"strings"
	"time"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/internal/utils"
	"dfood/pkg/errors"
)

type OrderService interface {
	CreateOrder(userID string, orderRequest *models.CreateOrderRequest) (*models.Order, error)
	GetUserOrders(userID string, limit, offset int) ([]models.Order, error)
	GetByID(orderID string) (*models.Order, error)
	UpdateStatus(orderID string, status models.OrderStatus) error
	CancelOrder(orderID string) error
	TrackOrder(orderID string) (interface{}, error)
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

func (s *orderService) CreateOrder(userID string, orderRequest *models.CreateOrderRequest) (*models.Order, error) {
	if orderRequest == nil {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order request is required", nil)
	}

	// Validate user ID
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Create order from request
	order := &models.Order{
		UserID:          userID,
		RestaurantID:    orderRequest.RestaurantID,
		RestaurantName:  orderRequest.RestaurantName,
		Items:           make(models.OrderItemsArray, len(orderRequest.Items)),
		Subtotal:        orderRequest.Subtotal,
		DeliveryFee:     orderRequest.DeliveryFee,
		Tax:             orderRequest.Tax,
		Total:           orderRequest.Total,
		DeliveryAddress: orderRequest.DeliveryAddress,
		PaymentMethod:   orderRequest.PaymentMethodID, // Map PaymentMethodID to PaymentMethod
		Notes:           orderRequest.Notes,
	}

	// Copy items from request
	for i, item := range orderRequest.Items {
		order.Items[i] = models.OrderItem{
			FoodID:              item.FoodID,
			FoodName:            item.Name,
			Price:               item.Price,
			Quantity:            item.Quantity,
			Total:               item.Price * float64(item.Quantity),
			SpecialInstructions: item.SpecialInstructions,
		}
	}

	// Validate required fields
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
		if food.RestaurantID != order.RestaurantID {
			return nil, errors.NewHTTPError(http.StatusBadRequest, "All items must be from the same restaurant", nil)
		}

		// Update item details
		order.Items[i].FoodName = food.Name
		order.Items[i].Price = food.Price
		order.Items[i].Total = food.Price * float64(item.Quantity)
		subtotal += order.Items[i].Total
	}

	// Set calculated totals
	order.Subtotal = subtotal
	if order.DeliveryFee <= 0 {
		order.DeliveryFee = 3.99 // Default delivery fee
	}
	if order.Tax <= 0 {
		order.Tax = subtotal * 0.08 // 8% tax rate
	}
	order.Total = order.Subtotal + order.DeliveryFee + order.Tax
	order.Status = models.OrderStatusPending

	// Generate order ID and set timestamps
	order.ID = utils.GenerateOrderID()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

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

func (s *orderService) GetByID(orderID string) (*models.Order, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
	}

	return s.orderRepo.GetByID(orderID)
}

func (s *orderService) UpdateStatus(orderID string, status models.OrderStatus) error {
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

func (s *orderService) TrackOrder(orderID string) (interface{}, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Order ID is required", nil)
	}

	return s.orderRepo.GetByID(orderID)
}
