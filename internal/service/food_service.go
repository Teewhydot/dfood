package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type FoodService interface {
	GetAll(limit, offset int) ([]models.Food, error)
	GetByID(id string) (*models.Food, error)
	GetPopular(limit int) ([]models.Food, error)
	GetRecommended(limit int) ([]models.Food, error)
	GetByCategory(category string, limit, offset int) ([]models.Food, error)
	GetByRestaurant(restaurantID string, limit, offset int) ([]models.Food, error)
	Search(query string, limit, offset int) ([]models.Food, error)
}

type foodService struct {
	foodRepo repository.FoodRepository
}

func NewFoodService(foodRepo repository.FoodRepository) FoodService {
	return &foodService{
		foodRepo: foodRepo,
	}
}

func (s *foodService) GetAll(limit, offset int) ([]models.Food, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.foodRepo.GetAll(limit, offset)
}

func (s *foodService) GetByID(id string) (*models.Food, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Food ID is required", nil)
	}

	return s.foodRepo.GetByID(id)
}

func (s *foodService) GetPopular(limit int) ([]models.Food, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	return s.foodRepo.GetPopular(limit)
}

func (s *foodService) GetRecommended(limit int) ([]models.Food, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	// For now, return popular foods as recommendations
	// In a real app, this would use ML algorithms, user preferences, order history, etc.
	return s.foodRepo.GetPopular(limit)
}

func (s *foodService) GetByCategory(category string, limit, offset int) ([]models.Food, error) {
	if strings.TrimSpace(category) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Category is required", nil)
	}
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.foodRepo.GetByCategory(category, limit, offset)
}

func (s *foodService) GetByRestaurant(restaurantID string, limit, offset int) ([]models.Food, error) {
	if strings.TrimSpace(restaurantID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 200 {
		limit = 200 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.foodRepo.GetByRestaurant(restaurantID, limit, offset)
}

func (s *foodService) Search(query string, limit, offset int) ([]models.Food, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Search query is required", nil)
	}
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.foodRepo.Search(query, limit, offset)
}
