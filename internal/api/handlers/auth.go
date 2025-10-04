package handlers

import (
	"net/http"

	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/pkg/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var newPasswordJson models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&newPasswordJson); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for login user",
		)
		result.RespondWithJSON(c)
		return
	}
	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.authService.UpdatePassword(newPasswordJson.Email, newPasswordJson.CurrentPassword, newPasswordJson.NewPassword)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
		"updating user password",
	)
	result.RespondWithJSON(c)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for new user",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleErrorWithStatusCode(
		func() (interface{}, error) {
			err := h.authService.Register(&newUser)
			if err != nil {
				return nil, err
			}
			newUser.Password = ""
			return newUser, nil
		},
		"creating new user",
		http.StatusCreated,
	)
	result.RespondWithJSON(c)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginUser models.User
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for login user",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			user, err := h.authService.Login(loginUser.Email, loginUser.Password)
			if err != nil {
				return nil, err
			}
			return user, nil
		},
		"logging in user",
	)
	result.RespondWithJSON(c)
}

// Additional Authentication Endpoints
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement user logout
	c.JSON(200, gin.H{"message": "User logout - TODO"})
}

func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	// TODO: Implement delete user account and associated data
	c.JSON(200, gin.H{"message": "Delete account - TODO"})
}

func (h *AuthHandler) SendPasswordReset(c *gin.Context) {
	// TODO: Implement send password reset email
	c.JSON(200, gin.H{"message": "Send password reset - TODO"})
}

func (h *AuthHandler) SendEmailVerification(c *gin.Context) {
	// TODO: Implement send email verification
	c.JSON(200, gin.H{"message": "Send email verification - TODO"})
}

func (h *AuthHandler) VerifyEmailStatus(c *gin.Context) {
	// TODO: Implement check email verification status
	c.JSON(200, gin.H{"message": "Verify email status - TODO"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// TODO: Implement get current authenticated user info
	c.JSON(200, gin.H{"message": "Get current user - TODO"})
}
