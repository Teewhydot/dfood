package service

import (
	"net/http"
	"time"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/internal/utils"
	"dfood/pkg/errors"
)

type AuthService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	UpdatePassword(email, currentPassword, newPassword string) error
	Logout(token string) error
	DeleteAccount(email, token string) error
	GetCurrentUser(token string) (*models.User, error)
}

type authService struct {
	userRepo     repository.UserRepository
	emailService EmailService
}

func NewAuthService(userRepo repository.UserRepository, emailService EmailService) AuthService {
	return &authService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *authService) UpdatePassword(email, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.NewHTTPError(http.StatusNotFound, "User not found", err)
	}

	// Check that the current password is valid for that account
	passwordIsValid := utils.CheckPasswordHash(user.Password, currentPassword)
	if !passwordIsValid {
		return errors.NewHTTPError(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	// Check if new password is not the same as the old password.
	passwordIsSame := utils.CheckPasswordHash(user.Password, newPassword)
	if passwordIsSame {
		return errors.NewHTTPError(http.StatusBadRequest, "New password must be different from the current password", nil)
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to hash new password", err)
	}

	err = s.userRepo.UpdatePassword(email, hashedPassword)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to update password", err)
	}

	return nil
}

func (s *authService) Register(user *models.User) error {
	exists, err := s.userRepo.EmailExists(user.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewHTTPError(http.StatusConflict, "User already exists", nil)
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = utils.GenerateID()
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to secure password", err)
	}

	user.Password = hashedPassword

	// Set default values
	user.FirstTimeLogin = true
	user.EmailVerified = false

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err = s.userRepo.Create(user)
	if err != nil {
		return err
	}

	// Send welcome email (don't fail registration if email fails)
	if s.emailService != nil {
		go func() {
			subject := "Welcome to dfood!"
			plainText := `Hi ` + user.FirstName + `,

Welcome to dfood! We're excited to have you join our community of food lovers.

You can now:
- Discover amazing restaurants
- Share your dining experiences through reviews
- Help others find great places to eat

Start exploring and sharing your food journey with us!

Best regards,
The dfood Team`

			htmlContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to dfood</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2c3e50;">Welcome to dfood!</h1>
        <p>Hi ` + user.FirstName + `,</p>
        <p>Welcome to dfood! We're excited to have you join our community of food lovers.</p>
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <h3 style="margin-top: 0;">You can now:</h3>
            <ul>
                <li>Discover amazing restaurants</li>
                <li>Share your dining experiences through reviews</li>
                <li>Help others find great places to eat</li>
            </ul>
        </div>
        <p>Start exploring and sharing your food journey with us!</p>
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`

			_ = s.emailService.SendEmail(user.Email, user.FirstName, subject, plainText, htmlContent)
		}()
	}

	return nil
}

func (s *authService) Login(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	passwordIsValid := utils.CheckPasswordHash(user.Password, password)
	if !passwordIsValid {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid credentials", nil)
	}

	user.Password = ""
	// Generate JWT token
	accessToken, err := utils.GenerateJwtToken(user.Email, false)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token", err)
	}
	refreshToken, err := utils.GenerateJwtToken(user.Email, true)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token", err)
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return user, nil
}

func (s *authService) Logout(token string) error {
	// Validate token first
	_, err := utils.ValidateToken(token)
	if err != nil {
		return errors.NewHTTPError(http.StatusUnauthorized, "Invalid token", err)
	}

	// Invalidate the token
	err = utils.InvalidateToken(token)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to invalidate token", err)
	}

	return nil
}

func (s *authService) DeleteAccount(email, token string) error {
	// Validate token first
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return errors.NewHTTPError(http.StatusUnauthorized, "Invalid token", err)
	}

	// Check if token belongs to the user
	if tokenEmail, ok := (*claims)["sub"].(string); !ok || tokenEmail != email {
		return errors.NewHTTPError(http.StatusForbidden, "Token does not belong to this user", nil)
	}

	// Check if user exists
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.NewHTTPError(http.StatusNotFound, "User not found", err)
	}

	// Invalidate all tokens for this user
	utils.InvalidateAllUserTokens(email)

	// TODO: Delete user from database
	// This would require implementing a Delete method in UserRepository
	_ = user

	return errors.NewHTTPError(http.StatusNotImplemented, "Account deletion not implemented", nil)
}

func (s *authService) GetCurrentUser(token string) (*models.User, error) {
	// Validate token first
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid token", err)
	}

	// Extract email from token
	email, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid token claims", nil)
	}

	// Get user from database
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusNotFound, "User not found", err)
	}

	// Clear sensitive information
	user.Password = ""

	return user, nil
}
