package repository

import (
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type favoritesRepository struct {
	db *gorm.DB
}

func NewFavoritesRepository() FavoritesRepository {
	return &favoritesRepository{
		db: database.DB,
	}
}

func (r *favoritesRepository) GetFavoriteFoods(userID string) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.Table("foods").
		Joins("INNER JOIN favorite_foods ON foods.id = favorite_foods.food_id").
		Where("favorite_foods.user_id = ?", userID).
		Find(&foods).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch favorite foods", err)
	}
	return foods, nil
}

func (r *favoritesRepository) GetFavoriteRestaurants(userID string) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Table("restaurants").
		Joins("INNER JOIN favorite_restaurants ON restaurants.id = favorite_restaurants.restaurant_id").
		Where("favorite_restaurants.user_id = ?", userID).
		Find(&restaurants).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch favorite restaurants", err)
	}
	return restaurants, nil
}

func (r *favoritesRepository) AddFavoriteFood(userID, foodID string) error {
	favorite := models.FavoriteFood{
		UserID: userID,
		FoodID: foodID,
	}
	if err := r.db.Create(&favorite).Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to add favorite food", err)
	}
	return nil
}

func (r *favoritesRepository) RemoveFavoriteFood(userID, foodID string) error {
	err := r.db.Where("user_id = ? AND food_id = ?", userID, foodID).Delete(&models.FavoriteFood{}).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to remove favorite food", err)
	}
	return nil
}

func (r *favoritesRepository) AddFavoriteRestaurant(userID, restaurantID string) error {
	favorite := models.FavoriteRestaurant{
		UserID:       userID,
		RestaurantID: restaurantID,
	}
	if err := r.db.Create(&favorite).Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to add favorite restaurant", err)
	}
	return nil
}

func (r *favoritesRepository) RemoveFavoriteRestaurant(userID, restaurantID string) error {
	err := r.db.Where("user_id = ? AND restaurant_id = ?", userID, restaurantID).Delete(&models.FavoriteRestaurant{}).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to remove favorite restaurant", err)
	}
	return nil
}

func (r *favoritesRepository) IsFoodFavorite(userID, foodID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.FavoriteFood{}).Where("user_id = ? AND food_id = ?", userID, foodID).Count(&count).Error
	if err != nil {
		return false, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to check food favorite status", err)
	}
	return count > 0, nil
}

func (r *favoritesRepository) IsRestaurantFavorite(userID, restaurantID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.FavoriteRestaurant{}).Where("user_id = ? AND restaurant_id = ?", userID, restaurantID).Count(&count).Error
	if err != nil {
		return false, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to check restaurant favorite status", err)
	}
	return count > 0, nil
}

func (r *favoritesRepository) ClearAllFavorites(userID string) error {
	// Start transaction
	tx := r.db.Begin()

	// Delete all favorite foods
	if err := tx.Where("user_id = ?", userID).Delete(&models.FavoriteFood{}).Error; err != nil {
		tx.Rollback()
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to clear favorite foods", err)
	}

	// Delete all favorite restaurants
	if err := tx.Where("user_id = ?", userID).Delete(&models.FavoriteRestaurant{}).Error; err != nil {
		tx.Rollback()
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to clear favorite restaurants", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction", err)
	}

	return nil
}

func (r *favoritesRepository) GetFavoritesStats(userID string) (map[string]int, error) {
	stats := make(map[string]int)

	// Count favorite foods
	var foodCount int64
	if err := r.db.Model(&models.FavoriteFood{}).Where("user_id = ?", userID).Count(&foodCount).Error; err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to count favorite foods", err)
	}
	stats["foods"] = int(foodCount)

	// Count favorite restaurants
	var restaurantCount int64
	if err := r.db.Model(&models.FavoriteRestaurant{}).Where("user_id = ?", userID).Count(&restaurantCount).Error; err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to count favorite restaurants", err)
	}
	stats["restaurants"] = int(restaurantCount)

	return stats, nil
}
