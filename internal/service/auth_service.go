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
	RefreshToken(refreshToken string) (*models.User, error)
	SendEmailVerification(email string) error
	VerifyEmail(token string) error
	SendPasswordReset(email string) error
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

func (s *authService) RefreshToken(refreshToken string) (*models.User, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token", err)
	}

	// Check if it's actually a refresh token
	isRefresh, ok := (*claims)["is_refresh"].(bool)
	if !ok || !isRefresh {
		return nil, errors.NewHTTPError(http.StatusUnauthorized, "Token is not a refresh token", nil)
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

	// Generate new tokens
	accessToken, err := utils.GenerateJwtToken(user.Email, false)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate access token", err)
	}

	newRefreshToken, err := utils.GenerateJwtToken(user.Email, true)
	if err != nil {
		return nil, errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate refresh token", err)
	}

	// Clear sensitive information
	user.Password = ""
	user.AccessToken = accessToken
	user.RefreshToken = newRefreshToken

	return user, nil
}

func (s *authService) SendEmailVerification(email string) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.NewHTTPError(http.StatusNotFound, "User not found", err)
	}

	// Check if already verified
	if user.EmailVerified {
		return errors.NewHTTPError(http.StatusBadRequest, "Email already verified", nil)
	}

	// Generate verification token (valid for 24 hours)
	verificationToken, err := utils.GenerateVerificationToken(email, 24*time.Hour)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to generate verification token", err)
	}

	// Send verification email
	if s.emailService != nil {
		subject := "Verify your dfood account"
		plainText := `Hi ` + user.FirstName + `,

Please verify your email address by clicking the link below:

http://localhost:8080/api/v1/auth/verify-email?token=` + verificationToken + `

This link will expire in 24 hours.

If you didn't create an account with dfood, please ignore this email.

Best regards,
The dfood Team`

		htmlContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify your dfood account</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2c3e50;">Verify your dfood account</h1>
        <p>Hi ` + user.FirstName + `,</p>
        <p>Please verify your email address by clicking the button below:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="http://localhost:8080/api/v1/auth/verify-email?token=` + verificationToken + `" 
               style="background-color: #3498db; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">
                Verify Email Address
            </a>
        </div>
        <p>This link will expire in 24 hours.</p>
        <p>If you didn't create an account with dfood, please ignore this email.</p>
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`

		err = s.emailService.SendEmail(user.Email, user.FirstName, subject, plainText, htmlContent)
		if err != nil {
			return errors.NewHTTPError(http.StatusInternalServerError, "Failed to send verification email", err)
		}
	}

	return nil
}

func (s *authService) VerifyEmail(token string) error {
	// Validate verification token
	claims, err := utils.ValidateVerificationToken(token)
	if err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid or expired verification token", err)
	}

	// Extract email from token
	email, ok := (*claims)["sub"].(string)
	if !ok {
		return errors.NewHTTPError(http.StatusBadRequest, "Invalid token claims", nil)
	}

	// Get user from database
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.NewHTTPError(http.StatusNotFound, "User not found", err)
	}

	// Check if already verified
	if user.EmailVerified {
		return errors.NewHTTPError(http.StatusBadRequest, "Email already verified", nil)
	}

	// Update user's email verification status
	err = s.userRepo.UpdateEmailVerification(email, true)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to update email verification status", err)
	}

	return nil
}

func (s *authService) SendPasswordReset(email string) error {
	// Check if user exists (but don't reveal if they don't for security)
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Return success even if user doesn't exist for security
		return nil
	}

	// Generate a new random password
	newPassword := utils.GenerateRandomPassword(12)

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to hash new password", err)
	}

	// Update the user's password in the database
	err = s.userRepo.UpdatePassword(email, hashedPassword)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to update password", err)
	}

	// Invalidate all existing tokens for this user
	utils.InvalidateAllUserTokens(email)

	// Send the new password via email
	if s.emailService != nil {
		subject := "Your new dfood password"
		plainText := `Hi ` + user.FirstName + `,

We received a request to reset your password for your dfood account.

Your new temporary password is: ` + newPassword + `

Please use this password to log in and then update it to something more memorable using the "Update Password" feature in your account settings.

For security reasons, we recommend changing this password as soon as possible after logging in.

If you didn't request a password reset, please contact our support team immediately.

Best regards,
The dfood Team`

		htmlContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your new dfood password</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2c3e50;">Your new dfood password</h1>
        <p>Hi ` + user.FirstName + `,</p>
        <p>We received a request to reset your password for your dfood account.</p>
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0; text-align: center;">
            <h3 style="margin-top: 0; color: #e74c3c;">Your new temporary password:</h3>
            <div style="font-family: 'Courier New', monospace; font-size: 18px; font-weight: bold; background-color: #fff; padding: 15px; border: 2px solid #e74c3c; border-radius: 5px; display: inline-block;">
                ` + newPassword + `
            </div>
        </div>
        <p><strong>Important:</strong> Please use this password to log in and then update it to something more memorable using the "Update Password" feature in your account settings.</p>
        <p>For security reasons, we recommend changing this password as soon as possible after logging in.</p>
        <p>If you didn't request a password reset, please contact our support team immediately.</p>
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`

		err = s.emailService.SendEmail(user.Email, user.FirstName, subject, plainText, htmlContent)
		if err != nil {
			return errors.NewHTTPError(http.StatusInternalServerError, "Failed to send password reset email", err)
		}
	}

	return nil
}
