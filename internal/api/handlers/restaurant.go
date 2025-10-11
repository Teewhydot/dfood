package handlers

import (
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"
	"strconv"

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
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			restaurants, err := h.restaurantService.GetAll(limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"restaurants": restaurants,
				"limit":       limit,
				"offset":      offset,
			}, nil
		},
		"getting all restaurants",
	)
	result.RespondWithJSON(c)
}

func (h *RestaurantHandler) GetRestaurantByID(c *gin.Context) {
	restaurantID := c.Param("id")

	result := errors.HandleError(
		func() (interface{}, error) {
			restaurant, err := h.restaurantService.GetByID(restaurantID)
			if err != nil {
				return nil, err
			}
			return restaurant, nil
		},
		"getting restaurant by ID",
	)
	result.RespondWithJSON(c)
}

func (h *RestaurantHandler) GetPopularRestaurants(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			restaurants, err := h.restaurantService.GetPopular(limit)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"restaurants": restaurants,
				"limit":       limit,
			}, nil
		},
		"getting popular restaurants",
	)
	result.RespondWithJSON(c)
}

func (h *RestaurantHandler) GetNearbyRestaurants(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5.0")
	limitStr := c.DefaultQuery("limit", "20")

	if latStr == "" || lngStr == "" {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Latitude and longitude are required", nil)
			},
			"validating location parameters",
		)
		result.RespondWithJSON(c)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid latitude", err)
			},
			"parsing latitude",
		)
		result.RespondWithJSON(c)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid longitude", err)
			},
			"parsing longitude",
		)
		result.RespondWithJSON(c)
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		radius = 5.0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 20
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			restaurants, err := h.restaurantService.GetNearby(lat, lng, radius, limit)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"restaurants": restaurants,
				"latitude":    lat,
				"longitude":   lng,
				"radius":      radius,
				"limit":       limit,
			}, nil
		},
		"getting nearby restaurants",
	)
	result.RespondWithJSON(c)
}

func (h *RestaurantHandler) SearchRestaurants(c *gin.Context) {
	query := c.Query("query")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	if query == "" {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Search query is required", nil)
			},
			"validating search query",
		)
		result.RespondWithJSON(c)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			restaurants, err := h.restaurantService.Search(query, limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"restaurants": restaurants,
				"query":       query,
				"limit":       limit,
				"offset":      offset,
			}, nil
		},
		"searching restaurants",
	)
	result.RespondWithJSON(c)
}

func (h *RestaurantHandler) GetRestaurantsByCategory(c *gin.Context) {
	category := c.Param("category")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			restaurants, err := h.restaurantService.GetByCategory(category, limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"restaurants": restaurants,
				"category":    category,
				"limit":       limit,
				"offset":      offset,
			}, nil
		},
		"getting restaurants by category",
	)
	result.RespondWithJSON(c)
}

func (h *RestaurantHandler) GetRestaurantMenu(c *gin.Context) {
	restaurantID := c.Param("id")

	result := errors.HandleError(
		func() (interface{}, error) {
			menu, err := h.restaurantService.GetMenu(restaurantID)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"restaurant_id": restaurantID,
				"menu":          menu,
			}, nil
		},
		"getting restaurant menu",
	)
	result.RespondWithJSON(c)
}
