package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type RestaurantService interface {
	GetAll(limit, offset int) ([]models.Restaurant, error)
	GetByID(id string) (*models.Restaurant, error)
	GetPopular(limit int) ([]models.Restaurant, error)
	GetNearby(latitude, longitude float64, radius float64, limit int) ([]models.Restaurant, error)
	Search(query string, limit, offset int) ([]models.Restaurant, error)
	GetByCategory(category string, limit, offset int) ([]models.Restaurant, error)
	GetMenu(restaurantID string) ([]models.Food, error)
}

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
	foodRepo       repository.FoodRepository
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepository, foodRepo repository.FoodRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
		foodRepo:       foodRepo,
	}
}

func (s *restaurantService) GetAll(limit, offset int) ([]models.Restaurant, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.restaurantRepo.GetAll(limit, offset)
}

func (s *restaurantService) GetByID(id string) (*models.Restaurant, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	return s.restaurantRepo.GetByID(id)
}

func (s *restaurantService) GetPopular(limit int) ([]models.Restaurant, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	return s.restaurantRepo.GetPopular(limit)
}

func (s *restaurantService) GetNearby(latitude, longitude float64, radius float64, limit int) ([]models.Restaurant, error) {
	if latitude < -90 || latitude > 90 {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid latitude", nil)
	}
	if longitude < -180 || longitude > 180 {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid longitude", nil)
	}
	if radius <= 0 {
		radius = 10 // Default 10km radius
	}
	if radius > 100 {
		radius = 100 // Max 100km radius
	}
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	return s.restaurantRepo.GetNearby(latitude, longitude, radius, limit)
}

func (s *restaurantService) Search(query string, limit, offset int) ([]models.Restaurant, error) {
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

	return s.restaurantRepo.Search(query, limit, offset)
}

func (s *restaurantService) GetByCategory(category string, limit, offset int) ([]models.Restaurant, error) {
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

	return s.restaurantRepo.GetByCategory(category, limit, offset)
}

func (s *restaurantService) GetMenu(restaurantID string) ([]models.Food, error) {
	if strings.TrimSpace(restaurantID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Validate restaurant exists
	_, err := s.restaurantRepo.GetByID(restaurantID)
	if err != nil {
		return nil, err
	}

	// Get all menu items for the restaurant
	return s.foodRepo.GetByRestaurant(restaurantID, 200, 0)
}
