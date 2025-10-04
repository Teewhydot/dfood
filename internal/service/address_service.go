package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type AddressService interface {
	GetUserAddresses(userID string) ([]models.Address, error)
	SaveAddress(address *models.Address) error
	UpdateAddress(addressID string, updates map[string]interface{}) error
	DeleteAddress(addressID string) error
	GetDefaultAddress(userID string) (*models.Address, error)
	SetDefaultAddress(userID, addressID string) error
}

type addressService struct {
	addressRepo repository.AddressRepository
	userRepo    repository.UserRepository
}

func NewAddressService(addressRepo repository.AddressRepository, userRepo repository.UserRepository) AddressService {
	return &addressService{
		addressRepo: addressRepo,
		userRepo:    userRepo,
	}
}

func (s *addressService) GetUserAddresses(userID string) ([]models.Address, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.addressRepo.GetByUserID(userID)
}

func (s *addressService) SaveAddress(address *models.Address) error {
	if address == nil {
		return errors.NewHTTPError(http.StatusBadRequest, "Address is required", nil)
	}

	// Validate required fields
	if strings.TrimSpace(address.UserID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(address.Street) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Street is required", nil)
	}
	if strings.TrimSpace(address.City) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "City is required", nil)
	}
	if strings.TrimSpace(address.State) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "State is required", nil)
	}
	if strings.TrimSpace(address.ZipCode) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Zip code is required", nil)
	}
	if strings.TrimSpace(address.Address) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Address is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(address.UserID)
	if err != nil {
		return err
	}

	// Validate address type
	validTypes := map[string]bool{"home": true, "work": true, "other": true}
	if !validTypes[address.Type] {
		address.Type = "home" // Default to home
	}

	return s.addressRepo.Create(address)
}

func (s *addressService) UpdateAddress(addressID string, updates map[string]interface{}) error {
	if strings.TrimSpace(addressID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Address ID is required", nil)
	}

	if len(updates) == 0 {
		return errors.NewHTTPError(http.StatusBadRequest, "No updates provided", nil)
	}

	// Validate address exists
	_, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return err
	}

	// Remove fields that shouldn't be updated this way
	delete(updates, "id")
	delete(updates, "user_id")
	delete(updates, "created_at")

	// Validate address type if being updated
	if addressType, exists := updates["type"]; exists {
		typeStr, ok := addressType.(string)
		if !ok {
			return errors.NewHTTPError(http.StatusBadRequest, "Invalid address type format", nil)
		}
		validTypes := map[string]bool{"home": true, "work": true, "other": true}
		if !validTypes[typeStr] {
			return errors.NewHTTPError(http.StatusBadRequest, "Invalid address type", nil)
		}
	}

	return s.addressRepo.Update(addressID, updates)
}

func (s *addressService) DeleteAddress(addressID string) error {
	if strings.TrimSpace(addressID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Address ID is required", nil)
	}

	// Validate address exists
	_, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return err
	}

	return s.addressRepo.Delete(addressID)
}

func (s *addressService) GetDefaultAddress(userID string) (*models.Address, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.addressRepo.GetDefaultByUserID(userID)
}

func (s *addressService) SetDefaultAddress(userID, addressID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(addressID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Address ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Validate address exists and belongs to user
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return errors.NewHTTPError(http.StatusForbidden, "Address does not belong to user", nil)
	}

	return s.addressRepo.SetDefault(userID, addressID)
}
