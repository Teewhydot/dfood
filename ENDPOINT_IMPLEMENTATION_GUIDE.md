# Endpoint Implementation Guide

This guide provides a complete step-by-step process for adding a new endpoint to the dfood application, from model definition to route registration.

## Architecture Overview

The dfood application follows a **Clean Architecture** pattern with four main layers:

```
Route → Handler → Service → Repository → Database
```

Each layer has specific responsibilities:
- **Route**: HTTP endpoint mapping
- **Handler**: Request/response handling, validation
- **Service**: Business logic, validation, coordination
- **Repository**: Data access, database queries
- **Database**: GORM ORM, SQLite

## Complete Implementation Steps

Let's implement a new feature: **Review Management** for restaurants.

### Step 1: Define the Model

**File**: `internal/models/review.go`

```go
package models

import (
	"time"
)

type Review struct {
	ID           string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID       string    `json:"userId" gorm:"type:varchar(255);not null"`
	RestaurantID string    `json:"restaurantId" gorm:"type:varchar(255);not null"`
	Rating       float64   `json:"rating" gorm:"type:real;not null"` // 1.0 - 5.0
	Comment      string    `json:"comment" gorm:"type:text"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relationships
	User       User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Restaurant Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}

// Table name
func (Review) TableName() string {
	return "reviews"
}

// Request/Response DTOs
type CreateReviewRequest struct {
	RestaurantID string  `json:"restaurantId" binding:"required"`
	Rating       float64 `json:"rating" binding:"required,min=1,max=5"`
	Comment      string  `json:"comment" binding:"max=500"`
}

type UpdateReviewRequest struct {
	Rating  float64 `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment string  `json:"comment" binding:"omitempty,max=500"`
}
```

**Key Points**:
- Use descriptive struct tags for JSON and GORM
- Define validation rules in request DTOs using `binding` tags
- Include created/updated timestamps
- Define relationships for joins if needed

---

### Step 2: Define Repository Interface

**File**: `internal/repository/interfaces.go`

Add the interface to the existing file:

```go
type ReviewRepository interface {
	Create(review *models.Review) error
	GetByID(id string) (*models.Review, error)
	GetByRestaurantID(restaurantID string, limit, offset int) ([]models.Review, error)
	GetByUserID(userID string, limit, offset int) ([]models.Review, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error
	GetAverageRating(restaurantID string) (float64, error)
}
```

**Key Points**:
- Define all data access methods
- Use standard naming: Create, GetByID, Update, Delete
- Include pagination parameters (limit, offset) for list operations
- Return error as second value

---

### Step 3: Implement Repository

**File**: `internal/repository/review_repo.go`

```go
package repository

import (
	"errors"
	"net/http"

	"dfood/internal/database"
	"dfood/internal/models"
	pkgErrors "dfood/pkg/errors"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() ReviewRepository {
	return &reviewRepository{
		db: database.DB,
	}
}

func (r *reviewRepository) Create(review *models.Review) error {
	err := r.db.Create(review).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to create review", err)
	}
	return nil
}

func (r *reviewRepository) GetByID(id string) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("id = ?", id).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewHTTPError(http.StatusNotFound, "Review not found", err)
		}
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch review", err)
	}
	return &review, nil
}

func (r *reviewRepository) GetByRestaurantID(restaurantID string, limit, offset int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Preload("User"). // Load user data
		Find(&reviews).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch restaurant reviews", err)
	}
	return reviews, nil
}

func (r *reviewRepository) GetByUserID(userID string, limit, offset int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Preload("Restaurant"). // Load restaurant data
		Find(&reviews).Error
	if err != nil {
		return nil, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user reviews", err)
	}
	return reviews, nil
}

func (r *reviewRepository) Update(id string, updates map[string]interface{}) error {
	err := r.db.Model(&models.Review{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to update review", err)
	}
	return nil
}

func (r *reviewRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.Review{}).Error
	if err != nil {
		return pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to delete review", err)
	}
	return nil
}

func (r *reviewRepository) GetAverageRating(restaurantID string) (float64, error) {
	var avgRating float64
	err := r.db.Model(&models.Review{}).
		Where("restaurant_id = ?", restaurantID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	if err != nil {
		return 0, pkgErrors.NewHTTPError(http.StatusInternalServerError, "Failed to calculate average rating", err)
	}
	return avgRating, nil
}
```

**Key Points**:
- Always wrap database errors with `pkgErrors.NewHTTPError`
- Check for `gorm.ErrRecordNotFound` and return 404
- Use `Preload` for eager loading relationships
- Use `Updates(map[string]interface{})` for flexible updates
- Initialize repository with `database.DB` from the database package

---

### Step 4: Define Service Interface

**File**: Create `internal/service/review_service.go`

```go
package service

import (
	"net/http"
	"strings"
	"time"

	"dfood/internal/models"
	"dfood/internal/repository"
	"dfood/internal/utils"
	"dfood/pkg/errors"
)

type ReviewService interface {
	CreateReview(userID string, req models.CreateReviewRequest) (*models.Review, error)
	GetReviewByID(id string) (*models.Review, error)
	GetRestaurantReviews(restaurantID string, limit, offset int) ([]models.Review, error)
	GetUserReviews(userID string, limit, offset int) ([]models.Review, error)
	UpdateReview(id, userID string, req models.UpdateReviewRequest) (*models.Review, error)
	DeleteReview(id, userID string) error
	GetRestaurantAverageRating(restaurantID string) (float64, error)
}

type reviewService struct {
	reviewRepo     repository.ReviewRepository
	restaurantRepo repository.RestaurantRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, restaurantRepo repository.RestaurantRepository) ReviewService {
	return &reviewService{
		reviewRepo:     reviewRepo,
		restaurantRepo: restaurantRepo,
	}
}
```

**Key Points**:
- Define interface above implementation
- Service may need multiple repositories
- Match repository methods but add business logic

---

### Step 5: Implement Service Methods

**File**: Continue in `internal/service/review_service.go`

```go
func (s *reviewService) CreateReview(userID string, req models.CreateReviewRequest) (*models.Review, error) {
	// Validate user ID
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Validate restaurant exists
	_, err := s.restaurantRepo.GetByID(req.RestaurantID)
	if err != nil {
		return nil, err // Error already wrapped by repository
	}

	// Validate rating
	if req.Rating < 1.0 || req.Rating > 5.0 {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Rating must be between 1.0 and 5.0", nil)
	}

	// Create review
	review := &models.Review{
		ID:           utils.GenerateID("review"),
		UserID:       userID,
		RestaurantID: req.RestaurantID,
		Rating:       req.Rating,
		Comment:      req.Comment,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.reviewRepo.Create(review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) GetReviewByID(id string) (*models.Review, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Review ID is required", nil)
	}

	return s.reviewRepo.GetByID(id)
}

func (s *reviewService) GetRestaurantReviews(restaurantID string, limit, offset int) ([]models.Review, error) {
	if strings.TrimSpace(restaurantID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	// Apply default and max limits
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.reviewRepo.GetByRestaurantID(restaurantID, limit, offset)
}

func (s *reviewService) GetUserReviews(userID string, limit, offset int) ([]models.Review, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Apply default and max limits
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.reviewRepo.GetByUserID(userID, limit, offset)
}

func (s *reviewService) UpdateReview(id, userID string, req models.UpdateReviewRequest) (*models.Review, error) {
	// Validate inputs
	if strings.TrimSpace(id) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "Review ID is required", nil)
	}
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Get existing review
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if review.UserID != userID {
		return nil, errors.NewHTTPError(http.StatusForbidden, "You can only update your own reviews", nil)
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Rating > 0 {
		if req.Rating < 1.0 || req.Rating > 5.0 {
			return nil, errors.NewHTTPError(http.StatusBadRequest, "Rating must be between 1.0 and 5.0", nil)
		}
		updates["rating"] = req.Rating
	}

	if req.Comment != "" {
		updates["comment"] = req.Comment
	}

	updates["updated_at"] = time.Now()

	// Update review
	err = s.reviewRepo.Update(id, updates)
	if err != nil {
		return nil, err
	}

	// Return updated review
	return s.reviewRepo.GetByID(id)
}

func (s *reviewService) DeleteReview(id, userID string) error {
	// Validate inputs
	if strings.TrimSpace(id) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "Review ID is required", nil)
	}
	if strings.TrimSpace(userID) == "" {
		return errors.NewHTTPError(http.StatusBadRequest, "User ID is required", nil)
	}

	// Get existing review
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verify ownership
	if review.UserID != userID {
		return errors.NewHTTPError(http.StatusForbidden, "You can only delete your own reviews", nil)
	}

	return s.reviewRepo.Delete(id)
}

func (s *reviewService) GetRestaurantAverageRating(restaurantID string) (float64, error) {
	if strings.TrimSpace(restaurantID) == "" {
		return 0, errors.NewHTTPError(http.StatusBadRequest, "Restaurant ID is required", nil)
	}

	return s.reviewRepo.GetAverageRating(restaurantID)
}
```

**Key Points**:
- Always validate inputs (IDs, required fields)
- Implement business rules (e.g., ownership verification)
- Apply default limits and caps for pagination
- Generate IDs using `utils.GenerateID()`
- Set timestamps explicitly
- Coordinate between multiple repositories when needed

---

### Step 6: Create Handler

**File**: `internal/api/handlers/review.go`

```go
package handlers

import (
	"net/http"
	"strconv"

	"dfood/internal/models"
	"dfood/internal/service"
	"dfood/pkg/errors"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// CreateReview handles POST /api/v1/reviews
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not authenticated", nil)
			},
			"getting user ID from context",
		)
		result.RespondWithJSON(c)
		return
	}

	// Parse request body
	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for create review",
		)
		result.RespondWithJSON(c)
		return
	}

	// Call service
	result := errors.HandleErrorWithStatusCode(
		func() (interface{}, error) {
			review, err := h.reviewService.CreateReview(userID.(string), req)
			if err != nil {
				return nil, err
			}
			return review, nil
		},
		"creating review",
		http.StatusCreated,
	)
	result.RespondWithJSON(c)
}

// GetReview handles GET /api/v1/reviews/:id
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id := c.Param("id")

	result := errors.HandleError(
		func() (interface{}, error) {
			review, err := h.reviewService.GetReviewByID(id)
			if err != nil {
				return nil, err
			}
			return review, nil
		},
		"getting review",
	)
	result.RespondWithJSON(c)
}

// GetRestaurantReviews handles GET /api/v1/restaurants/:restaurantId/reviews
func (h *ReviewHandler) GetRestaurantReviews(c *gin.Context) {
	restaurantID := c.Param("restaurantId")

	// Parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result := errors.HandleError(
		func() (interface{}, error) {
			reviews, err := h.reviewService.GetRestaurantReviews(restaurantID, limit, offset)
			if err != nil {
				return nil, err
			}
			return reviews, nil
		},
		"getting restaurant reviews",
	)
	result.RespondWithJSON(c)
}

// GetUserReviews handles GET /api/v1/users/:userId/reviews
func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userID := c.Param("userId")

	// Parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	result := errors.HandleError(
		func() (interface{}, error) {
			reviews, err := h.reviewService.GetUserReviews(userID, limit, offset)
			if err != nil {
				return nil, err
			}
			return reviews, nil
		},
		"getting user reviews",
	)
	result.RespondWithJSON(c)
}

// UpdateReview handles PUT /api/v1/reviews/:id
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not authenticated", nil)
			},
			"getting user ID from context",
		)
		result.RespondWithJSON(c)
		return
	}

	// Parse request body
	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON payload", err)
			},
			"binding JSON for update review",
		)
		result.RespondWithJSON(c)
		return
	}

	// Call service
	result := errors.HandleError(
		func() (interface{}, error) {
			review, err := h.reviewService.UpdateReview(id, userID.(string), req)
			if err != nil {
				return nil, err
			}
			return review, nil
		},
		"updating review",
	)
	result.RespondWithJSON(c)
}

// DeleteReview handles DELETE /api/v1/reviews/:id
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		result := errors.HandleError(
			func() (interface{}, error) {
				return nil, errors.NewHTTPError(http.StatusUnauthorized, "User not authenticated", nil)
			},
			"getting user ID from context",
		)
		result.RespondWithJSON(c)
		return
	}

	// Call service
	result := errors.HandleError(
		func() (interface{}, error) {
			err := h.reviewService.DeleteReview(id, userID.(string))
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
		"deleting review",
	)
	result.RespondWithJSON(c)
}

// GetRestaurantAverageRating handles GET /api/v1/restaurants/:restaurantId/rating
func (h *ReviewHandler) GetRestaurantAverageRating(c *gin.Context) {
	restaurantID := c.Param("restaurantId")

	result := errors.HandleError(
		func() (interface{}, error) {
			rating, err := h.reviewService.GetRestaurantAverageRating(restaurantID)
			if err != nil {
				return nil, err
			}
			return gin.H{"restaurantId": restaurantID, "averageRating": rating}, nil
		},
		"getting restaurant average rating",
	)
	result.RespondWithJSON(c)
}
```

**Key Points**:
- Use `errors.HandleError()` for all operations
- Use `errors.HandleErrorWithStatusCode()` for custom status codes (e.g., 201 for creation)
- Extract parameters using `c.Param()` for path params
- Extract query params using `c.Query()` or `c.DefaultQuery()`
- Extract user context from authenticated requests
- Always validate JSON binding
- Use `result.RespondWithJSON(c)` to send response

---

### Step 7: Register Routes

**File**: `internal/api/routes/routes.go`

#### 7.1 Add to Dependencies struct:

```go
type Dependencies struct {
	AuthService         service.AuthService
	UserService         service.UserService
	RestaurantService   service.RestaurantService
	ReviewService       service.ReviewService  // Add this
	// ... other services
}
```

#### 7.2 Initialize handler in SetupRoutes:

```go
func SetupRoutes(deps *Dependencies) *gin.Engine {
	router := gin.New()

	// Global Middleware
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())
	router.Use(middleware.RateLimitMiddleware(10, time.Minute))

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(deps.AuthService)
	userHandler := handlers.NewUserHandler(deps.UserService)
	restaurantHandler := handlers.NewRestaurantHandler(deps.RestaurantService)
	reviewHandler := handlers.NewReviewHandler(deps.ReviewService) // Add this
	// ... other handlers
```

#### 7.3 Register routes:

```go
	// API v1 Routes
	v1 := router.Group("/api/v1")
	{
		// ... existing routes ...

		// Restaurant Reviews
		restaurants := v1.Group("/restaurants")
		{
			// ... existing restaurant routes ...
			restaurants.GET("/:restaurantId/reviews", reviewHandler.GetRestaurantReviews)
			restaurants.GET("/:restaurantId/rating", reviewHandler.GetRestaurantAverageRating)
		}

		// Review Endpoints
		reviews := v1.Group("/reviews")
		{
			// Public routes
			reviews.GET("/:id", reviewHandler.GetReview)

			// Protected routes (require authentication)
			reviews.POST("", middleware.AuthMiddleware(), reviewHandler.CreateReview)
			reviews.PUT("/:id", middleware.AuthMiddleware(), reviewHandler.UpdateReview)
			reviews.DELETE("/:id", middleware.AuthMiddleware(), reviewHandler.DeleteReview)
		}

		// User Reviews
		users := v1.Group("/users")
		{
			// ... existing user routes ...
			users.GET("/:userId/reviews", reviewHandler.GetUserReviews)
		}
	}

	return router
}
```

**Key Points**:
- Group related routes together
- Use middleware for protected routes
- RESTful naming conventions
- Logical route organization

---

### Step 8: Dependency Injection

**File**: `cmd/main.go`

```go
func main() {
	// ... configuration loading ...

	// Initialize Repositories
	userRepo := repository.NewUserRepository()
	restaurantRepo := repository.NewRestaurantRepository()
	reviewRepo := repository.NewReviewRepository() // Add this
	// ... other repos ...

	// Initialize Services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	restaurantService := service.NewRestaurantService(restaurantRepo)
	reviewService := service.NewReviewService(reviewRepo, restaurantRepo) // Add this
	// ... other services ...

	// Dependency Injection
	deps := &routes.Dependencies{
		AuthService:       authService,
		UserService:       userService,
		RestaurantService: restaurantService,
		ReviewService:     reviewService, // Add this
		// ... other services ...
	}

	// Setup and run server
	router := routes.SetupRoutes(deps)
	logger.Info("Starting server", "port", cfg.Port)
	if err := router.Run(":" + fmt.Sprint(cfg.Port)); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
```

**Key Points**:
- Follow dependency order: repositories → services → dependencies
- Pass all required repositories to services
- All services are singleton instances

---

### Step 9: Create API Test File

**File**: `api-test/reviews.http`

```http
### Variables
@baseURL = http://localhost:8080/api/v1
@token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
@restaurantId = restaurant-123
@reviewId = review-123

### Create Review (Protected)
POST {{baseURL}}/reviews
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "restaurantId": "{{restaurantId}}",
  "rating": 4.5,
  "comment": "Great food and excellent service!"
}

### Get Review by ID
GET {{baseURL}}/reviews/{{reviewId}}

### Get Restaurant Reviews (with pagination)
GET {{baseURL}}/restaurants/{{restaurantId}}/reviews?limit=10&offset=0

### Get User Reviews
GET {{baseURL}}/users/user-123/reviews?limit=10&offset=0

### Update Review (Protected)
PUT {{baseURL}}/reviews/{{reviewId}}
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "rating": 5.0,
  "comment": "Updated: Amazing experience!"
}

### Delete Review (Protected)
DELETE {{baseURL}}/reviews/{{reviewId}}
Authorization: Bearer {{token}}

### Get Restaurant Average Rating
GET {{baseURL}}/restaurants/{{restaurantId}}/rating
```

**Key Points**:
- Use variables for reusable values
- Test all endpoints
- Include authentication tokens for protected routes
- Test pagination
- Test error cases

---

## Error Handling Pattern

All handlers use the standardized error handling pattern:

```go
result := errors.HandleError(
	func() (interface{}, error) {
		// Your operation here
		data, err := service.SomeMethod()
		if err != nil {
			return nil, err
		}
		return data, nil
	},
	"operation description",
)
result.RespondWithJSON(c)
```

For custom status codes (e.g., 201 Created):

```go
result := errors.HandleErrorWithStatusCode(
	func() (interface{}, error) {
		// Your operation here
		return data, nil
	},
	"operation description",
	http.StatusCreated,
)
result.RespondWithJSON(c)
```

---

## Common Patterns

### 1. ID Generation

```go
import "dfood/internal/utils"

id := utils.GenerateID("review") // Generates: review-uuid
```

### 2. Pagination

```go
// In Service
if limit <= 0 {
	limit = 20 // Default
}
if limit > 100 {
	limit = 100 // Max
}
if offset < 0 {
	offset = 0
}

// In Handler
limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
```

### 3. Input Validation

```go
// Required field
if strings.TrimSpace(id) == "" {
	return nil, errors.NewHTTPError(http.StatusBadRequest, "ID is required", nil)
}

// Range validation
if rating < 1.0 || rating > 5.0 {
	return nil, errors.NewHTTPError(http.StatusBadRequest, "Rating must be between 1.0 and 5.0", nil)
}
```

### 4. Ownership Verification

```go
// Get resource
resource, err := repo.GetByID(id)
if err != nil {
	return nil, err
}

// Verify ownership
if resource.UserID != userID {
	return nil, errors.NewHTTPError(http.StatusForbidden, "You can only modify your own resources", nil)
}
```

### 5. Flexible Updates

```go
updates := make(map[string]interface{})

if req.Field1 != "" {
	updates["field1"] = req.Field1
}

if req.Field2 > 0 {
	updates["field2"] = req.Field2
}

updates["updated_at"] = time.Now()

err = repo.Update(id, updates)
```

---

## Testing Checklist

- [ ] Model defined with proper tags
- [ ] Repository interface added
- [ ] Repository implementation complete
- [ ] Service interface defined
- [ ] Service implementation complete with validation
- [ ] Handler created with error handling
- [ ] Routes registered in routes.go
- [ ] Dependencies injected in main.go
- [ ] API test file created
- [ ] All endpoints tested manually
- [ ] Error cases tested

---

## File Structure Summary

For a new "Review" feature, create/modify these files:

```
internal/
├── models/
│   └── review.go                      # New model
├── repository/
│   ├── interfaces.go                  # Add interface
│   └── review_repo.go                 # New repository
├── service/
│   └── review_service.go              # New service
├── api/
│   ├── handlers/
│   │   └── review.go                  # New handler
│   └── routes/
│       └── routes.go                  # Modify to add routes
cmd/
└── main.go                            # Modify for DI
api-test/
└── reviews.http                       # New test file
```

---

## Quick Reference Commands

```bash
# Run application
go run cmd/main.go

# Build application
go build -o dfood cmd/main.go

# Run specific test
curl -X GET http://localhost:8080/api/v1/reviews/review-123

# Test with authentication
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"restaurantId":"rest-123","rating":4.5,"comment":"Great!"}'
```

---

## Best Practices

1. **Always validate inputs** at the service layer
2. **Use custom errors** from `pkg/errors` package
3. **Apply pagination limits** to prevent abuse
4. **Verify ownership** before updates/deletes
5. **Use transactions** for multi-step database operations
6. **Load relationships** with `Preload()` when needed
7. **Clear sensitive data** before returning (e.g., passwords)
8. **Use meaningful operation names** in error handlers
9. **Follow RESTful conventions** for route naming
10. **Document endpoints** in .http test files

---

This guide provides a complete, production-ready pattern for implementing new endpoints in the dfood application. Follow these steps for consistency and maintainability.
