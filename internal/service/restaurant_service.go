package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type RestaurantService interface {
	GetAllRestaurants(limit, offset int) ([]models.Restaurant, error)
	GetRestaurantByID(id string) (*models.Restaurant, error)
	GetPopularRestaurants(limit int) ([]models.Restaurant, error)
	GetNearbyRestaurants(latitude, longitude float64, radius float64, limit int) ([]models.Restaurant, error)
	SearchRestaurants(query string, limit, offset int) ([]models.Restaurant, error)
	GetRestaurantsByCategory(category string, limit, offset int) ([]models.Restaurant, error)
	GetRestaurantMenu(restaurantID string, limit, offset int) ([]models.Food, error)
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

func (s *restaurantService) GetAllRestaurants(limit, offset int) ([]models.Restaurant, error) {
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

func (s *restaurantService) GetRestaurantByID(id string) (*models.Restaurant, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	return s.restaurantRepo.GetByID(id)
}

func (s *restaurantService) GetPopularRestaurants(limit int) ([]models.Restaurant, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 50 {
		limit = 50 // Max limit
	}

	return s.restaurantRepo.GetPopular(limit)
}

func (s *restaurantService) GetNearbyRestaurants(latitude, longitude float64, radius float64, limit int) ([]models.Restaurant, error) {
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

func (s *restaurantService) SearchRestaurants(query string, limit, offset int) ([]models.Restaurant, error) {
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

func (s *restaurantService) GetRestaurantsByCategory(category string, limit, offset int) ([]models.Restaurant, error) {
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

func (s *restaurantService) GetRestaurantMenu(restaurantID string, limit, offset int) ([]models.Food, error) {
	if strings.TrimSpace(restaurantID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Validate restaurant exists
	_, err := s.restaurantRepo.GetByID(restaurantID)
	if err != nil {
		return nil, err
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
