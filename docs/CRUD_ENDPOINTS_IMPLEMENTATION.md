# CRUD Endpoints Implementation

## Overview
This document outlines the complete CRUD (Create, Read, Update, Delete) endpoints implementation for the dfood application's core entities.

## ✅ Implemented Endpoints

### 1. User Profile Management (`/users/*`)

#### Read Operations
- `GET /users/:userId` - Get user profile by ID
  - Returns user profile with sensitive data removed
  - Validates user ID parameter

#### Update Operations  
- `PUT /users/:userId` - Update entire user profile
  - Accepts JSON payload with profile updates
  - Prevents updating sensitive fields (password, email, id, created_at)
  - Returns success message

- `PATCH /users/:userId/:field` - Update specific profile field
  - Updates single field with validation
  - Allowed fields: first_name, last_name, phone_number, bio, profile_image_url, fcm_token
  - Returns success message

#### Delete Operations
- `DELETE /users/:userId/profile-image` - Remove profile image
  - Sets profile_image_url to null
  - Returns success message

### 2. Restaurant Management (`/restaurants/*`)

#### Read Operations
- `GET /restaurants` - Get all restaurants (paginated)
  - Query params: limit (default: 20, max: 100), offset (default: 0)
  - Returns restaurants array with pagination info

- `GET /restaurants/:id` - Get restaurant by ID
  - Returns single restaurant details
  - Validates restaurant ID parameter

- `GET /restaurants/popular` - Get popular restaurants
  - Query params: limit (default: 10, max: 50)
  - Returns highly-rated restaurants

- `GET /restaurants/nearby` - Get nearby restaurants
  - Query params: lat (required), lng (required), radius (default: 5km), limit (default: 20)
  - Validates coordinates and radius
  - Returns restaurants within specified radius

- `GET /restaurants/search` - Search restaurants
  - Query params: query (required), limit (default: 20), offset (default: 0)
  - Searches by name, category, description
  - Returns matching restaurants with search metadata

- `GET /restaurants/category/:category` - Get restaurants by category
  - Path param: category (required)
  - Query params: limit (default: 20), offset (default: 0)
  - Returns restaurants in specified category

- `GET /restaurants/:id/menu` - Get restaurant menu
  - Returns all food items for the restaurant
  - Includes menu organization and categories

### 3. Food/Menu Management (`/foods/*`)

#### Read Operations
- `GET /foods` - Get all foods (paginated)
  - Query params: limit (default: 20, max: 100), offset (default: 0)
  - Returns food items with pagination info

- `GET /foods/:id` - Get food item by ID
  - Returns single food item details
  - Includes restaurant information

- `GET /foods/popular` - Get popular food items
  - Query params: limit (default: 10, max: 50)
  - Returns highly-rated food items

- `GET /foods/recommended` - Get recommended foods
  - Query params: limit (default: 10, max: 50)
  - Currently returns popular foods (placeholder for ML recommendations)

- `GET /foods/category/:category` - Get foods by category
  - Path param: category (required)
  - Query params: limit (default: 20), offset (default: 0)
  - Returns foods in specified category

- `GET /foods/restaurant/:restaurantId` - Get foods by restaurant
  - Path param: restaurantId (required)
  - Query params: limit (default: 20), offset (default: 0)
  - Returns all menu items for specified restaurant

- `GET /foods/search` - Search food items
  - Query params: query (required), limit (default: 20), offset (default: 0)
  - Searches by name, description, category, restaurant
  - Returns matching foods with search metadata

### 4. Order Management (`/orders/*`)

#### Create Operations
- `POST /orders` - Create new order
  - Accepts CreateOrderRequest JSON payload
  - Validates required fields (restaurant, items, totals, address, payment)
  - Returns created order with generated ID

#### Read Operations
- `GET /orders/user/:userId` - Get user's order history
  - Path param: userId (required)
  - Query params: limit (default: 20), offset (default: 0)
  - Returns user's orders with pagination

- `GET /orders/:orderId` - Get order details
  - Path param: orderId (required)
  - Returns complete order information including items and status

- `GET /orders/:orderId/track` - Track order status
  - Path param: orderId (required)
  - Returns real-time order tracking information

#### Update Operations
- `PUT /orders/:orderId/status` - Update order status
  - Path param: orderId (required)
  - Accepts UpdateOrderStatusRequest JSON payload
  - Updates order status and delivery information

#### Delete Operations
- `DELETE /orders/:orderId` - Cancel order
  - Path param: orderId (required)
  - Cancels order if cancellation is allowed
  - Returns success message

## 🔧 Service Layer Architecture

### Service Interfaces
All endpoints are backed by service layer interfaces that provide:
- Input validation and sanitization
- Business logic enforcement
- Error handling with appropriate HTTP status codes
- Data transformation and formatting

### Key Service Methods

#### UserService
```go
type UserService interface {
    GetByID(userID string) (*models.User, error)
    Update(userID string, updates map[string]interface{}) error
    UpdateField(userID, field string, value interface{}) error
    // ... other methods
}
```

#### RestaurantService
```go
type RestaurantService interface {
    GetAll(limit, offset int) ([]models.Restaurant, error)
    GetByID(id string) (*models.Restaurant, error)
    GetPopular(limit int) ([]models.Restaurant, error)
    GetNearby(lat, lng, radius float64, limit int) ([]models.Restaurant, error)
    Search(query string, limit, offset int) ([]models.Restaurant, error)
    GetByCategory(category string, limit, offset int) ([]models.Restaurant, error)
    GetMenu(restaurantID string) ([]models.Food, error)
}
```

#### FoodService
```go
type FoodService interface {
    GetAll(limit, offset int) ([]models.Food, error)
    GetByID(id string) (*models.Food, error)
    GetPopular(limit int) ([]models.Food, error)
    GetRecommended(limit int) ([]models.Food, error)
    GetByCategory(category string, limit, offset int) ([]models.Food, error)
    GetByRestaurant(restaurantID string, limit, offset int) ([]models.Food, error)
    Search(query string, limit, offset int) ([]models.Food, error)
}
```

#### OrderService
```go
type OrderService interface {
    CreateOrder(orderRequest *models.CreateOrderRequest) (*models.Order, error)
    GetUserOrders(userID string, limit, offset int) ([]models.Order, error)
    GetByID(orderID string) (*models.Order, error)
    UpdateStatus(orderID string, status models.OrderStatus) error
    CancelOrder(orderID string) error
    TrackOrder(orderID string) (interface{}, error)
}
```

## 🛡️ Security & Validation

### Input Validation
- All endpoints validate required parameters
- Pagination limits are enforced (max 100 items per request)
- Coordinate validation for location-based queries
- Field validation for profile updates

### Error Handling
- Consistent error response format
- Appropriate HTTP status codes
- Detailed error messages for debugging
- Security-conscious error messages (no sensitive data leakage)

### Data Sanitization
- Sensitive data removal (passwords, tokens)
- Input trimming and validation
- SQL injection prevention through parameterized queries

## 📊 Response Formats

### Success Responses
```json
{
  "data": [...],
  "limit": 20,
  "offset": 0,
  "query": "search term"
}
```

### Error Responses
```json
{
  "error": "Error message",
  "details": "Additional context"
}
```

## 🧪 Testing

### Test Coverage
- Happy path scenarios for all endpoints
- Error cases and validation testing
- Pagination boundary testing
- Search functionality validation
- Location-based query testing

### Test Files
- `api-test/crud-endpoints.http` - Comprehensive CRUD endpoint testing
- Includes positive and negative test cases
- Covers all implemented endpoints with sample data

## 🚀 Performance Considerations

### Pagination
- Default limits prevent large data dumps
- Configurable limits with reasonable maximums
- Offset-based pagination for consistent results

### Search Optimization
- Query parameter validation
- Search result limiting
- Efficient database queries through repository layer

### Caching Opportunities
- Popular items can be cached
- Restaurant menus can be cached
- Search results can be cached for common queries

## 📈 Future Enhancements

### Advanced Features
1. **Filtering & Sorting**
   - Price range filtering
   - Rating-based sorting
   - Distance-based sorting
   - Dietary restriction filtering

2. **Real-time Features**
   - WebSocket support for live order tracking
   - Real-time menu updates
   - Live restaurant availability

3. **Analytics**
   - Search analytics
   - Popular item tracking
   - User behavior insights

4. **Recommendations**
   - ML-based food recommendations
   - Personalized restaurant suggestions
   - Order history-based recommendations

## 🔗 Integration Points

### Database Layer
- Repository pattern for data access
- GORM for ORM functionality
- SQLite for development/testing

### External Services
- Payment processing integration
- Geolocation services
- Image upload services
- Push notification services

This implementation provides a solid foundation for a food delivery application with comprehensive CRUD operations, proper validation, and scalable architecture.