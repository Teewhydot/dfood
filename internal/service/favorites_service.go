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
}

type favoritesService struct {
	favoritesRepo  repository.FavoritesRepository
	userRepo       repository.UserRepository
	foodRepo       repository.FoodRepository
	restaurantRepo repository.RestaurantRepository
}

func NewFavoritesService(favoritesRepo repository.FavoritesRepository, userRepo repository.UserRepository, foodRepo repository.FoodRepository, restaurantRepo repository.RestaurantRepository) FavoritesService {
	return &favoritesService{
		favoritesRepo:  favoritesRepo,
		userRepo:       userRepo,
		foodRepo:       foodRepo,
		restaurantRepo: restaurantRepo,
	}
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

	return s.favoritesRepo.AddFavoriteFood(userID, foodID)
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

	return s.favoritesRepo.RemoveFavoriteFood(userID, foodID)
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

	return s.favoritesRepo.AddFavoriteRestaurant(userID, restaurantID)
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

	return s.favoritesRepo.RemoveFavoriteRestaurant(userID, restaurantID)
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

	return s.favoritesRepo.ClearAllFavorites(userID)
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
