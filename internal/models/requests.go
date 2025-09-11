package models

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


// UpdatePasswordModel represents password update request
type UpdatePasswordRequest struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Token       string `json:"token"`
	UserProfile User   `json:"userProfile"`
}

// CreateOrderRequest represents create order request
type CreateOrderRequest struct {
	RestaurantID    string      `json:"restaurantId" binding:"required"`
	RestaurantName  string      `json:"restaurantName" binding:"required"`
	Items           []OrderItem `json:"items" binding:"required,min=1"`
	Subtotal        float64     `json:"subtotal" binding:"required,min=0"`
	DeliveryFee     float64     `json:"deliveryFee" binding:"required,min=0"`
	Tax             float64     `json:"tax" binding:"required,min=0"`
	Total           float64     `json:"total" binding:"required,min=0"`
	DeliveryAddress string      `json:"deliveryAddress" binding:"required"`
	PaymentMethodID string      `json:"paymentMethodId" binding:"required"`
	Notes           *string     `json:"notes,omitempty"`
}

// UpdateOrderStatusRequest represents update order status request
type UpdateOrderStatusRequest struct {
	Status              OrderStatus `json:"status" binding:"required"`
	DeliveryPersonName  *string     `json:"deliveryPersonName,omitempty"`
	DeliveryPersonPhone *string     `json:"deliveryPersonPhone,omitempty"`
	TrackingURL         *string     `json:"trackingUrl,omitempty"`
}

// CreateAddressRequest represents create address request
type CreateAddressRequest struct {
	Street    string   `json:"street" binding:"required"`
	City      string   `json:"city" binding:"required"`
	State     string   `json:"state" binding:"required"`
	ZipCode   string   `json:"zipCode" binding:"required"`
	Type      string   `json:"type" binding:"required"`
	Address   string   `json:"address" binding:"required"`
	Apartment string   `json:"apartment" binding:"required"`
	Title     *string  `json:"title,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	IsDefault bool     `json:"isDefault"`
}

// SendMessageRequest represents send message request
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"messageType"` // Default: "text"
}

// RestaurantSearchParams represents restaurant search parameters
type RestaurantSearchParams struct {
	Query     string  `form:"query"`
	Lat       float64 `form:"lat"`
	Lng       float64 `form:"lng"`
	Radius    float64 `form:"radius"` // in kilometers
	Category  string  `form:"category"`
	MinRating float64 `form:"minRating"`
	IsOpen    *bool   `form:"isOpen"`
	Page      int     `form:"page" binding:"min=1"`
	Limit     int     `form:"limit" binding:"min=1,max=100"`
}

// FoodSearchParams represents food search parameters
type FoodSearchParams struct {
	Query        string  `form:"query"`
	Category     string  `form:"category"`
	RestaurantID string  `form:"restaurantId"`
	MinRating    float64 `form:"minRating"`
	MaxPrice     float64 `form:"maxPrice"`
	IsVegetarian *bool   `form:"isVegetarian"`
	IsVegan      *bool   `form:"isVegan"`
	IsGlutenFree *bool   `form:"isGlutenFree"`
	Page         int     `form:"page" binding:"min=1"`
	Limit        int     `form:"limit" binding:"min=1,max=100"`
}
