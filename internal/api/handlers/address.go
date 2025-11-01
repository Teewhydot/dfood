package handlers

import (
	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/internal/utils"
	"dfood/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	userID := c.Param("userId")

	// Authorization: Verify the authenticated user matches the requested userId
	authenticatedUserID, exists := c.Get("user_id")
	if !exists {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not authenticated", nil)
			},
			"getting authenticated user ID",
		)
		result.RespondWithJSON(c)
		return
	}

	if authenticatedUserID.(string) != userID {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusForbidden, "You are not authorized to access this user's addresses", nil)
			},
			"validating user authorization",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			addresses, err := h.addressService.GetUserAddresses(userID)
			if err != nil {
				return nil, err
			}
			return addresses, nil
		},
		"getting user addresses",
	)
	result.RespondWithJSON(c)
}

func (h *AddressHandler) SaveAddress(c *gin.Context) {
	userID := c.Param("userId")

	// Authorization: Verify the authenticated user matches the requested userId
	authenticatedUserID, exists := c.Get("user_id")
	if !exists {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not authenticated", nil)
			},
			"getting authenticated user ID",
		)
		result.RespondWithJSON(c)
		return
	}

	if authenticatedUserID.(string) != userID {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusForbidden, "You are not authorized to save addresses for this user", nil)
			},
			"validating user authorization",
		)
		result.RespondWithJSON(c)
		return
	}

	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for address creation",
		)
		result.RespondWithJSON(c)
		return
	}

	// Set userID from URL parameter
	address.UserID = userID

	// Generate ID if not provided
	if address.ID == "" {
		address.ID = "addr-" + utils.GenerateID()
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.addressService.SaveAddress(&address)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"message": "Address saved successfully",
				"address": address,
			}, nil
		},
		"saving address",
	)
	result.RespondWithJSON(c)
}
