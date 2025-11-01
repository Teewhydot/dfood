package handlers

import (
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Profile Management
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.Param("userId")

	result := errors.HandleError(
		func() (interface{}, error) {
			user, err := h.userService.GetByID(userID)
			if err != nil {
				return nil, err
			}
			// Clear sensitive information
			user.Password = ""
			return user, nil
		},
		"getting user profile",
	)
	result.RespondWithJSON(c)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("userId")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for profile update",
		)
		result.RespondWithJSON(c)
		return
	}

	// Remove sensitive fields that shouldn't be updated via this endpoint
	delete(updateData, "password")
	delete(updateData, "email")
	delete(updateData, "id")
	delete(updateData, "created_at")

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.userService.Update(userID, updateData)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Profile updated successfully"}, nil
		},
		"updating user profile",
	)
	result.RespondWithJSON(c)
}

func (h *UserHandler) UpdateProfileField(c *gin.Context) {
	userID := c.Param("userId")
	field := c.Param("field")

	var fieldData struct {
		Value interface{} `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&fieldData); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for field update",
		)
		result.RespondWithJSON(c)
		return
	}

	// Validate allowed fields
	allowedFields := map[string]bool{
		"first_name":        true,
		"last_name":         true,
		"phone_number":      true,
		"bio":               true,
		"profile_image_url": true,
		"fcm_token":         true,
	}

	if !allowedFields[field] {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Field not allowed for update", nil)
			},
			"validating field name",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.userService.UpdateField(userID, field, fieldData.Value)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Field updated successfully"}, nil
		},
		"updating user field",
	)
	result.RespondWithJSON(c)
}

func (h *UserHandler) DeleteProfileImage(c *gin.Context) {
	userID := c.Param("userId")

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.userService.UpdateField(userID, "profile_image_url", nil)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Profile image deleted successfully"}, nil
		},
		"deleting profile image",
	)
	result.RespondWithJSON(c)
}

