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
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "Authorization header required", nil)
			},
			"getting authorization header",
		)
		result.RespondWithJSON(c)
		return
	}

	// Extract token (remove "Bearer " prefix)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.authService.Logout(token)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Successfully logged out"}, nil
		},
		"logging out user",
	)
	result.RespondWithJSON(c)
}

func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	var deleteRequest models.DeleteAccountRequest
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for delete account",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.authService.DeleteAccount(deleteRequest.Email, deleteRequest.Token)
			if err != nil {
				return nil, err
			}
			return gin.H{"message": "Account successfully deleted"}, nil
		},
		"deleting user account",
	)
	result.RespondWithJSON(c)
}

func (h *AuthHandler) SendPasswordReset(c *gin.Context) {
	var resetRequest models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&resetRequest); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for password reset",
		)
		result.RespondWithJSON(c)
		return
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			// For now, we'll return success without actually implementing the full flow
			// In a real implementation, you'd generate a reset token and send email
			return gin.H{
				"message": "If an account with that email exists, a password reset link has been sent",
				"email":   resetRequest.Email,
			}, nil
		},
		"sending password reset email",
	)
	result.RespondWithJSON(c)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "Authorization header required", nil)
			},
			"getting authorization header",
		)
		result.RespondWithJSON(c)
		return
	}

	// Extract token (remove "Bearer " prefix)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			user, err := h.authService.GetCurrentUser(token)
			if err != nil {
				return nil, err
			}
			return user, nil
		},
		"getting current user",
	)
	result.RespondWithJSON(c)
}
