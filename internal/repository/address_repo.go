package repository

import (
	"errors"
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository() AddressRepository {
	return &addressRepository{
		db: database.DB,
	}
}

func (r *addressRepository) GetByUserID(userID string) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user addresses", err)
	}
	return addresses, nil
}

func (r *addressRepository) Create(address *models.Address) error {
	if err := r.db.Create(address).Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to create address", err)
	}
	return nil
}

func (r *addressRepository) Update(id string, updates map[string]interface{}) error {
	err := r.db.Model(&models.Address{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update address", err)
	}
	return nil
}

func (r *addressRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.Address{}).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to delete address", err)
	}
	return nil
}

func (r *addressRepository) GetByID(id string) (*models.Address, error) {
	var address models.Address
	err := r.db.Where("id = ?", id).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "Address not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch address", err)
	}
	return &address, nil
}

func (r *addressRepository) GetDefaultByUserID(userID string) (*models.Address, error) {
	var address models.Address
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "Default address not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch default address", err)
	}
	return &address, nil
}

func (r *addressRepository) SetDefault(userID, addressID string) error {
	// Start transaction
	tx := r.db.Begin()

	// Unset all default addresses for the user
	if err := tx.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to unset default addresses", err)
	}

	// Set the new default address
	if err := tx.Model(&models.Address{}).Where("id = ? AND user_id = ?", addressID, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to set default address", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction", err)
	}

	return nil
}
