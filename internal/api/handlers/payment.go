package handlers

import (
	"dfood/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Payment Methods
func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	// TODO: Implement get available payment methods
	c.JSON(200, gin.H{"message": "Get payment methods - TODO"})
}

func (h *PaymentHandler) GetUserCards(c *gin.Context) {
	// TODO: Implement get user's saved cards
	c.JSON(200, gin.H{"message": "Get user cards - TODO"})
}

func (h *PaymentHandler) SaveCard(c *gin.Context) {
	// TODO: Implement save new payment card
	c.JSON(200, gin.H{"message": "Save card - TODO"})
}

func (h *PaymentHandler) DeleteCard(c *gin.Context) {
	// TODO: Implement delete saved card
	c.JSON(200, gin.H{"message": "Delete card - TODO"})
}

// Payment Processing
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	// TODO: Implement process payment for order
	c.JSON(200, gin.H{"message": "Process payment - TODO"})
}

func (h *PaymentHandler) GetTransactionDetails(c *gin.Context) {
	// TODO: Implement get payment transaction details
	c.JSON(200, gin.H{"message": "Get transaction details - TODO"})
}

func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	// TODO: Implement process refund
	c.JSON(200, gin.H{"message": "Process refund - TODO"})
}
