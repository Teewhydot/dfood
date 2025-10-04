package handlers

import (
	"dfood/internal/service"

	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	restaurantService service.RestaurantService
}

func NewRestaurantHandler(restaurantService service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

// Restaurant Data
func (h *RestaurantHandler) GetAllRestaurants(c *gin.Context) {
	// TODO: Implement get all restaurants (paginated)
	c.JSON(200, gin.H{"message": "Get all restaurants - TODO"})
}

func (h *RestaurantHandler) GetRestaurantByID(c *gin.Context) {
	// TODO: Implement get restaurant by ID
	c.JSON(200, gin.H{"message": "Get restaurant by ID - TODO"})
}

func (h *RestaurantHandler) GetPopularRestaurants(c *gin.Context) {
	// TODO: Implement get popular restaurants (rating >= 4.0)
	c.JSON(200, gin.H{"message": "Get popular restaurants - TODO"})
}

func (h *RestaurantHandler) GetNearbyRestaurants(c *gin.Context) {
	// TODO: Implement get nearby restaurants (lat/lng query params)
	c.JSON(200, gin.H{"message": "Get nearby restaurants - TODO"})
}

func (h *RestaurantHandler) SearchRestaurants(c *gin.Context) {
	// TODO: Implement search restaurants by name/category/description
	c.JSON(200, gin.H{"message": "Search restaurants - TODO"})
}

func (h *RestaurantHandler) GetRestaurantsByCategory(c *gin.Context) {
	// TODO: Implement get restaurants by category
	c.JSON(200, gin.H{"message": "Get restaurants by category - TODO"})
}

func (h *RestaurantHandler) GetRestaurantMenu(c *gin.Context) {
	// TODO: Implement get restaurant menu/categories with foods
	c.JSON(200, gin.H{"message": "Get restaurant menu - TODO"})
}
