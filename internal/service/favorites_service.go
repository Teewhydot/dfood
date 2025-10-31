package service

import (
	"net/http"
	"strings"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/pkg/errors"
)

type FavoritesService interface {
	GetFavoriteFoods(userID string) ([]models.Food, error)
	GetFavoriteRestaurants(userID string) ([]models.Restaurant, error)
	AddFavoriteFood(userID, foodID string) error
	RemoveFavoriteFood(userID, foodID string) error
	AddFavoriteRestaurant(userID, restaurantID string) error
	RemoveFavoriteRestaurant(userID, restaurantID string) error
	CheckFoodFavoriteStatus(userID, foodID string) (bool, error)
	CheckRestaurantFavoriteStatus(userID, restaurantID string) (bool, error)
	ToggleFoodFavorite(userID, foodID string) (bool, error)
	ToggleRestaurantFavorite(userID, restaurantID string) (bool, error)
	ClearAllFavorites(userID string) error
	GetFavoritesStats(userID string) (map[string]int, error)
	SetWebSocketService(wsService WebSocketService)
	GetWebSocketService() WebSocketService
	SetRealtimeService(realtimeService *RealtimeService)
}

type favoritesService struct {
	favoritesRepo   repository.FavoritesRepository
	userRepo        repository.UserRepository
	foodRepo        repository.FoodRepository
	restaurantRepo  repository.RestaurantRepository
	wsService       WebSocketService
	realtimeService *RealtimeService
}

func NewFavoritesService(favoritesRepo repository.FavoritesRepository, userRepo repository.UserRepository, foodRepo repository.FoodRepository, restaurantRepo repository.RestaurantRepository) FavoritesService {
	return &favoritesService{
		favoritesRepo:  favoritesRepo,
		userRepo:       userRepo,
		foodRepo:       foodRepo,
		restaurantRepo: restaurantRepo,
	}
}

// WebSocket service methods
func (s *favoritesService) SetWebSocketService(wsService WebSocketService) {
	s.wsService = wsService
}

func (s *favoritesService) GetWebSocketService() WebSocketService {
	return s.wsService
}

func (s *favoritesService) SetRealtimeService(realtimeService *RealtimeService) {
	s.realtimeService = realtimeService
}

func (s *favoritesService) GetFavoriteFoods(userID string) ([]models.Food, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.favoritesRepo.GetFavoriteFoods(userID)
}

func (s *favoritesService) GetFavoriteRestaurants(userID string) ([]models.Restaurant, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.favoritesRepo.GetFavoriteRestaurants(userID)
}

func (s *favoritesService) AddFavoriteFood(userID, foodID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(foodID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Food ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Validate food exists
	_, err = s.foodRepo.GetByID(foodID)
	if err != nil {
		return err
	}

	// Check if already favorite
	isFavorite, err := s.favoritesRepo.IsFoodFavorite(userID, foodID)
	if err != nil {
		return err
	}
	if isFavorite {
		return errors.NewHTTPError(http.StatusConflict, "Food is already in favorites", nil)
	}

	err = s.favoritesRepo.AddFavoriteFood(userID, foodID)
	if err != nil {
		return err
	}

	// Send realtime update
	if s.realtimeService != nil {
		// Get the food details to send in the update
		food, foodErr := s.foodRepo.GetByID(foodID)
		if foodErr == nil {
			s.realtimeService.SendFavoriteAdd(userID, "food", foodID, food)
		}
	}

	return nil
}

func (s *favoritesService) RemoveFavoriteFood(userID, foodID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(foodID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Food ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if is favorite
	isFavorite, err := s.favoritesRepo.IsFoodFavorite(userID, foodID)
	if err != nil {
		return err
	}
	if !isFavorite {
		return errors.NewHTTPError(http.StatusNotFound, "Food is not in favorites", nil)
	}

	// Get the food details before removing for realtime update
	var food *models.Food
	if s.realtimeService != nil {
		food, _ = s.foodRepo.GetByID(foodID)
	}

	err = s.favoritesRepo.RemoveFavoriteFood(userID, foodID)
	if err != nil {
		return err
	}

	// Send realtime update
	if s.realtimeService != nil && food != nil {
		s.realtimeService.SendFavoriteRemove(userID, "food", foodID, food)
	}

	return nil
}

func (s *favoritesService) AddFavoriteRestaurant(userID, restaurantID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(restaurantID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Validate restaurant exists
	_, err = s.restaurantRepo.GetByID(restaurantID)
	if err != nil {
		return err
	}

	// Check if already favorite
	isFavorite, err := s.favoritesRepo.IsRestaurantFavorite(userID, restaurantID)
	if err != nil {
		return err
	}
	if isFavorite {
		return errors.NewHTTPError(http.StatusConflict, "Restaurant is already in favorites", nil)
	}

	err = s.favoritesRepo.AddFavoriteRestaurant(userID, restaurantID)
	if err != nil {
		return err
	}

	// Send realtime update
	if s.realtimeService != nil {
		// Get the restaurant details to send in the update
		restaurant, restaurantErr := s.restaurantRepo.GetByID(restaurantID)
		if restaurantErr == nil {
			s.realtimeService.SendFavoriteAdd(userID, "restaurant", restaurantID, restaurant)
		}
	}

	return nil
}

func (s *favoritesService) RemoveFavoriteRestaurant(userID, restaurantID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(restaurantID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if is favorite
	isFavorite, err := s.favoritesRepo.IsRestaurantFavorite(userID, restaurantID)
	if err != nil {
		return err
	}
	if !isFavorite {
		return errors.NewHTTPError(http.StatusNotFound, "Restaurant is not in favorites", nil)
	}

	// Get the restaurant details before removing for realtime update
	var restaurant *models.Restaurant
	if s.realtimeService != nil {
		restaurant, _ = s.restaurantRepo.GetByID(restaurantID)
	}

	err = s.favoritesRepo.RemoveFavoriteRestaurant(userID, restaurantID)
	if err != nil {
		return err
	}

	// Send realtime update
	if s.realtimeService != nil && restaurant != nil {
		s.realtimeService.SendFavoriteRemove(userID, "restaurant", restaurantID, restaurant)
	}

	return nil
}

func (s *favoritesService) CheckFoodFavoriteStatus(userID, foodID string) (bool, error) {
	if strings.TrimSpace(userID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(foodID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "Food ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	return s.favoritesRepo.IsFoodFavorite(userID, foodID)
}

func (s *favoritesService) CheckRestaurantFavoriteStatus(userID, restaurantID string) (bool, error) {
	if strings.TrimSpace(userID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(restaurantID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	return s.favoritesRepo.IsRestaurantFavorite(userID, restaurantID)
}

func (s *favoritesService) ToggleFoodFavorite(userID, foodID string) (bool, error) {
	if strings.TrimSpace(userID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(foodID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "Food ID is required", nil)
	}

	// Check current status
	isFavorite, err := s.CheckFoodFavoriteStatus(userID, foodID)
	if err != nil {
		return false, err
	}

	if isFavorite {
		// Remove from favorites
		err = s.RemoveFavoriteFood(userID, foodID)
		return false, err
	} else {
		// Add to favorites
		err = s.AddFavoriteFood(userID, foodID)
		return true, err
	}
}

func (s *favoritesService) ToggleRestaurantFavorite(userID, restaurantID string) (bool, error) {
	if strings.TrimSpace(userID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}
	if strings.TrimSpace(restaurantID) == "" {
		return false, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Check current status
	isFavorite, err := s.CheckRestaurantFavoriteStatus(userID, restaurantID)
	if err != nil {
		return false, err
	}

	if isFavorite {
		// Remove from favorites
		err = s.RemoveFavoriteRestaurant(userID, restaurantID)
		return false, err
	} else {
		// Add to favorites
		err = s.AddFavoriteRestaurant(userID, restaurantID)
		return true, err
	}
}

func (s *favoritesService) ClearAllFavorites(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	err = s.favoritesRepo.ClearAllFavorites(userID)
	if err != nil {
		return err
	}

	// Send realtime update
	if s.realtimeService != nil {
		s.realtimeService.SendFavoritesClear(userID)
	}

	return nil
}

func (s *favoritesService) GetFavoritesStats(userID string) (map[string]int, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.favoritesRepo.GetFavoritesStats(userID)
}
