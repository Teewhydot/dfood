package repository

import (
	"errors"
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to register user", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "User not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user by email", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "User not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user details", err)
	}
	return &user, nil
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to check if user exists", err)
	}
	return count > 0, nil
}

func (r *userRepository) UpdatePassword(email, hashedPassword string) error {
	err := r.db.Model(&models.User{}).Where("email = ?", email).Update("password", hashedPassword).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update user password", err)
	}
	return nil
}

func (r *userRepository) Update(id string, updates map[string]interface{}) error {
	err := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update user", err)
	}
	return nil
}

func (r *userRepository) UpdateField(id, field string, value interface{}) error {
	err := r.db.Model(&models.User{}).Where("id = ?", id).Update(field, value).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update user field", err)
	}
	return nil
}

func (r *userRepository) UpdateFCMToken(id, token string) error {
	err := r.db.Model(&models.User{}).Where("id = ?", id).Update("fcm_token", token).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update FCM token", err)
	}
	return nil
}
