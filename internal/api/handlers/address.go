package handlers

import (
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type AddressHandler struct {
	addressService service.AddressService
}

func NewAddressHandler(addressService service.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// User Addresses
func (h *AddressHandler) GetUserAddresses(c *gin.Context) {
	userID := c.Param("userId") // ✅ CORRECT - matches route parameter
	result := errors.HandleError(
		func() (interface{}, error) {
			addresses, err := h.addressService.GetUserAddresses(userID) // ✅ Better variable name
			if err != nil {
				return nil, err
			}
			return addresses, nil
		},
		"getting user addresses", // ✅ Better error message
	)
	result.RespondWithJSON(c)
}

func (h *AddressHandler) SaveAddress(c *gin.Context) {

}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	// TODO: Implement update address
	c.JSON(200, gin.H{"message": "Update address - TODO"})
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	// TODO: Implement delete address
	c.JSON(200, gin.H{"message": "Delete address - TODO"})
}

func (h *AddressHandler) GetDefaultAddress(c *gin.Context) {
	// TODO: Implement get default address
	c.JSON(200, gin.H{"message": "Get default address - TODO"})
}

func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	// TODO: Implement set default address
	c.JSON(200, gin.H{"message": "Set default address - TODO"})
}

func (h *AddressHandler) GetAddressStream(c *gin.Context) {
	userID := c.Param("userId")

	// Verify user exists
	_, err := h.addressService.GetByUserID(userID)
	if err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, err
			},
			"verifying user addresses for WebSocket connection",
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

	// WebSocket functionality not implemented for addresses yet
	conn.Close()
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Address WebSocket streaming not implemented yet",
	})
}
