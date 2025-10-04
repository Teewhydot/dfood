package service

import (
	"net/http"

	"dfood/pkg/errors"
)

type UploadService interface {
	UploadProfileImage(userID string, imageData []byte, filename string) (string, error)
	UploadFoodImage(foodID string, imageData []byte, filename string) (string, error)
	UploadRestaurantImage(restaurantID string, imageData []byte, filename string) (string, error)
	DeleteImage(imageID string) error
}

type uploadService struct {
	// TODO: Add storage dependencies (AWS S3, local storage, etc.)
}

func NewUploadService() UploadService {
	return &uploadService{}
}

func (s *uploadService) UploadProfileImage(userID string, imageData []byte, filename string) (string, error) {
	// TODO: Implement profile image upload logic
	// Should handle image validation, resizing, and storage
	return "", errors.NewHTTPError(http.StatusNotImplemented, "Upload profile image not implemented", nil)
}

func (s *uploadService) UploadFoodImage(foodID string, imageData []byte, filename string) (string, error) {
	// TODO: Implement food image upload logic
	return "", errors.NewHTTPError(http.StatusNotImplemented, "Upload food image not implemented", nil)
}

func (s *uploadService) UploadRestaurantImage(restaurantID string, imageData []byte, filename string) (string, error) {
	// TODO: Implement restaurant image upload logic
	return "", errors.NewHTTPError(http.StatusNotImplemented, "Upload restaurant image not implemented", nil)
}

func (s *uploadService) DeleteImage(imageID string) error {
	// TODO: Implement image deletion logic
	return errors.NewHTTPError(http.StatusNotImplemented, "Delete image not implemented", nil)
}
