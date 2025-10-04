package repository

import (
	"errors"
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository() RestaurantRepository {
	return &restaurantRepository{
		db: database.DB,
	}
}

func (r *restaurantRepository) GetAll(limit, offset int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Limit(limit).Offset(offset).Find(&restaurants).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch restaurants", err)
	}
	return restaurants, nil
}

func (r *restaurantRepository) GetByID(id string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := r.db.Where("id = ?", id).First(&restaurant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "Restaurant not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch restaurant", err)
	}
	return &restaurant, nil
}

func (r *restaurantRepository) GetPopular(limit int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Where("rating >= ?", 4.0).Order("rating DESC").Limit(limit).Find(&restaurants).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch popular restaurants", err)
	}
	return restaurants, nil
}

func (r *restaurantRepository) GetNearby(latitude, longitude, radius float64, limit int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	// Using Haversine formula for distance calculation
	query := `
		SELECT *, (
			6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)
		) AS distance 
		FROM restaurants 
		HAVING distance < ? 
		ORDER BY distance 
		LIMIT ?
	`
	err := r.db.Raw(query, latitude, longitude, latitude, radius, limit).Scan(&restaurants).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch nearby restaurants", err)
	}
	return restaurants, nil
}

func (r *restaurantRepository) Search(query string, limit, offset int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	searchPattern := "%" + query + "%"
	err := r.db.Where("name LIKE ? OR description LIKE ? OR location LIKE ?",
		searchPattern, searchPattern, searchPattern).
		Limit(limit).Offset(offset).Find(&restaurants).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to search restaurants", err)
	}
	return restaurants, nil
}

func (r *restaurantRepository) GetByCategory(category string, limit, offset int) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Where("JSON_CONTAINS(categories, ?)", `"`+category+`"`).
		Limit(limit).Offset(offset).Find(&restaurants).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch restaurants by category", err)
	}
	return restaurants, nil
}
