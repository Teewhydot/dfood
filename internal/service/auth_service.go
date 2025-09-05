package service

import (
	"net/http"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/internal/utils"
	"dfood/pkg/errors"
)

type AuthService interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	UpdatePassword(email, currentPassword, newPassword string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
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

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.NewHTTPError(http.StatusInternalServerError, "Failed to secure password", err)
	}

	user.Password = hashedPassword
	return s.userRepo.Create(user)
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
