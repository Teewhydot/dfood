package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type UserService interface {
	GetProfile(userID string) (*models.User, error)
	UpdateProfile(userID string, updates map[string]interface{}) error
	UpdateProfileField(userID, field string, value interface{}) error
	UploadProfileImage(userID, imageURL string) error
	DeleteProfileImage(userID string) error
	SyncProfile(userID string) error
	UpdateFCMToken(userID, token string) error
	GetFCMToken(userID string) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(userID string) (*models.User, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Clear sensitive data
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateProfile(userID string, updates map[string]interface{}) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	if len(updates) == 0 {
		return errors.NewHTTPError(http.StatusBadRequest, "No updates provided", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Remove sensitive fields that shouldn't be updated this way
	delete(updates, "password")
	delete(updates, "id")
	delete(updates, "created_at")

	// Validate email if being updated
	if email, exists := updates["email"]; exists {
		emailStr, ok := email.(string)
		if !ok || strings.TrimSpace(emailStr) == "" {
			return errors.NewHTTPError(http.StatusBadRequest, "Invalid email format", nil)
		}

		// Check if email already exists for another user
		existingUser, _ := s.userRepo.GetByEmail(emailStr)
		if existingUser != nil && existingUser.ID != userID {
			return errors.NewHTTPError(http.StatusConflict, "Email already exists", nil)
		}
	}

	return s.userRepo.Update(userID, updates)
}

func (s *userService) UpdateProfileField(userID, field string, value interface{}) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	if strings.TrimSpace(field) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Field name is required", nil)
	}

	// Validate field is allowed to be updated
	allowedFields := map[string]bool{
		"first_name":        true,
		"last_name":         true,
		"phone_number":      true,
		"bio":               true,
		"profile_image_url": true,
	}

	if !allowedFields[field] {
		return errors.NewHTTPError(http.StatusBadRequest, "Field cannot be updated", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateField(userID, field, value)
}

func (s *userService) UploadProfileImage(userID, imageURL string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	if strings.TrimSpace(imageURL) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Image URL is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateField(userID, "profile_image_url", imageURL)
}

func (s *userService) DeleteProfileImage(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateField(userID, "profile_image_url", nil)
}

func (s *userService) SyncProfile(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// For sync, we just return the current profile (no actual sync logic needed for basic CRUD)
	// In a real app, this might sync with external services
	_ = user
	return nil
}

func (s *userService) UpdateFCMToken(userID, token string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	if strings.TrimSpace(token) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "FCM token is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateFCMToken(userID, token)
}

func (s *userService) GetFCMToken(userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", err
	}

	if user.FCMToken == nil {
		return "", errors.NewHTTPError(http.StatusNotFound, "FCM token not found", nil)
	}

	return *user.FCMToken, nil
}
