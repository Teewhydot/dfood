package handlers

import (
	"dfood/internal/service"
	"dfood/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	foodService service.FoodService
}

func NewFoodHandler(foodService service.FoodService) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
	}
}

// Food Data
func (h *FoodHandler) GetAllFoods(c *gin.Context) {
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
			foods, err := h.foodService.GetAll(limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"foods":  foods,
				"limit":  limit,
				"offset": offset,
			}, nil
		},
		"getting all foods",
	)
	result.RespondWithJSON(c)
}

func (h *FoodHandler) GetFoodByID(c *gin.Context) {
	foodID := c.Param("id")

	result := errors.HandleError(
		func() (interface{}, error) {
			food, err := h.foodService.GetByID(foodID)
			if err != nil {
				return nil, err
			}
			return food, nil
		},
		"getting food by ID",
	)
	result.RespondWithJSON(c)
}

func (h *FoodHandler) GetPopularFoods(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			foods, err := h.foodService.GetPopular(limit)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"foods": foods,
				"limit": limit,
			}, nil
		},
		"getting popular foods",
	)
	result.RespondWithJSON(c)
}

func (h *FoodHandler) GetRecommendedFoods(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	result := errors.HandleError(
		func() (interface{}, error) {
			foods, err := h.foodService.GetRecommended(limit)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"foods": foods,
				"limit": limit,
			}, nil
		},
		"getting recommended foods",
	)
	result.RespondWithJSON(c)
}

func (h *FoodHandler) GetFoodsByCategory(c *gin.Context) {
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
			foods, err := h.foodService.GetByCategory(category, limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"foods":    foods,
				"category": category,
				"limit":    limit,
				"offset":   offset,
			}, nil
		},
		"getting foods by category",
	)
	result.RespondWithJSON(c)
}

func (h *FoodHandler) GetFoodsByRestaurant(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
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
			foods, err := h.foodService.GetByRestaurant(restaurantID, limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"foods":         foods,
				"restaurant_id": restaurantID,
				"limit":         limit,
				"offset":        offset,
			}, nil
		},
		"getting foods by restaurant",
	)
	result.RespondWithJSON(c)
}

func (h *FoodHandler) SearchFoods(c *gin.Context) {
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
			foods, err := h.foodService.Search(query, limit, offset)
			if err != nil {
				return nil, err
			}
			return gin.H{
				"foods":  foods,
				"query":  query,
				"limit":  limit,
				"offset": offset,
			}, nil
		},
		"searching foods",
	)
	result.RespondWithJSON(c)
}
