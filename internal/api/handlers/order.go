package handlers

import (
	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Order Management
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderRequest models.CreateOrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for create order",
		)
		result.RespondWithJSON(c)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not authenticated", nil)
			},
			"getting user ID from context",
		)
		result.RespondWithJSON(c)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusInternalServerError, "Invalid user ID format", nil)
			},
			"converting user ID to string",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			order, err := h.orderService.CreateOrder(userIDStr, &orderRequest)
			if err != nil {
				return nil, err
			}
			return order, nil
		},
		"creating new order",
	)
	result.RespondWithJSON(c)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.Param("userId")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			orders, err := h.orderService.GetUserOrders(userID, limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"orders":  orders,
				"user_id": userID,
				"limit":   limit,
				"offset":  offset,
			}, nil
		},
		"getting user orders",
	)
	result.RespondWithJSON(c)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("orderId")

	result := errors.HandleError(
		func() (interface{}, error) {
			order, err := h.orderService.GetByID(orderID)
			if err != nil {
				return nil, err
			}
			return order, nil
		},
		"getting order by ID",
	)
	result.RespondWithJSON(c)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderId")

	var statusRequest models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for order status update",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.orderService.UpdateStatus(orderID, statusRequest.Status)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Order status updated successfully"}, nil
		},
		"updating order status",
	)
	result.RespondWithJSON(c)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.orderService.CancelOrder(orderID)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Order cancelled successfully"}, nil
		},
		"cancelling order",
	)
	result.RespondWithJSON(c)
}

func (h *OrderHandler) TrackOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	result := errors.HandleError(
		func() (interface{}, error) {
			tracking, err := h.orderService.TrackOrder(orderID)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"order_id": orderID,
				"tracking": tracking,
			}, nil
		},
		"tracking order",
	)
	result.RespondWithJSON(c)
}
