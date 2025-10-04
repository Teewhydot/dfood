package repository

import (
	"errors"
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type foodRepository struct {
	db *gorm.DB
}

func NewFoodRepository() FoodRepository {
	return &foodRepository{
		db: database.DB,
	}
}

func (r *foodRepository) GetAll(limit, offset int) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.Where("is_available = ?", true).Limit(limit).Offset(offset).Find(&foods).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch foods", err)
	}
	return foods, nil
}

func (r *foodRepository) GetByID(id string) (*models.Food, error) {
	var food models.Food
	err := r.db.Where("id = ?", id).First(&food).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "Food not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch food", err)
	}
	return &food, nil
}

func (r *foodRepository) GetPopular(limit int) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.Where("rating >= ? AND is_available = ?", 4.0, true).
		Order("rating DESC").Limit(limit).Find(&foods).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch popular foods", err)
	}
	return foods, nil
}

func (r *foodRepository) GetByCategory(category string, limit, offset int) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.Where("category = ? AND is_available = ?", category, true).
		Limit(limit).Offset(offset).Find(&foods).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch foods by category", err)
	}
	return foods, nil
}

func (r *foodRepository) GetByRestaurant(restaurantID string, limit, offset int) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.Where("restaurant_id = ? AND is_available = ?", restaurantID, true).
		Limit(limit).Offset(offset).Find(&foods).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch foods by restaurant", err)
	}
	return foods, nil
}

func (r *foodRepository) Search(query string, limit, offset int) ([]models.Food, error) {
	var foods []models.Food
	searchPattern := "%" + query + "%"
	err := r.db.Where("(name LIKE ? OR description LIKE ? OR restaurant_name LIKE ?) AND is_available = ?",
		searchPattern, searchPattern, searchPattern, true).
		Limit(limit).Offset(offset).Find(&foods).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to search foods", err)
	}
	return foods, nil
}
