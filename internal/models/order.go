package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// OrderStatus represents order status enum
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusOnTheWay  OrderStatus = "onTheWay"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents individual items in an order
type OrderItem struct {
	FoodID              string  `json:"food_id"`
	FoodName            string  `json:"food_name"`
	Price               float64 `json:"price"`
	Quantity            int     `json:"quantity"`
	Total               float64 `json:"total"`
	SpecialInstructions *string `json:"special_instructions,omitempty"`
}

// OrderItemsArray is a custom type for handling order items array
type OrderItemsArray []OrderItem

func (oia OrderItemsArray) Value() (driver.Value, error) {
	return json.Marshal(oia)
}

func (oia *OrderItemsArray) Scan(value interface{}) error {
	if value == nil {
		*oia = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, oia)
}

// Order represents the order entity
type Order struct {
	ID                  string          `json:"id" gorm:"primaryKey;column:id"`
	UserID              string          `json:"user_id" gorm:"column:user_id;not null;index"`
	RestaurantID        string          `json:"restaurant_id" gorm:"column:restaurant_id;not null;index"`
	RestaurantName      string          `json:"restaurant_name" gorm:"column:restaurant_name;not null"`
	Items               OrderItemsArray `json:"items" gorm:"column:items;not null"`
	Subtotal            float64         `json:"subtotal" gorm:"column:subtotal;not null"`
	DeliveryFee         float64         `json:"delivery_fee" gorm:"column:delivery_fee;not null"`
	Tax                 float64         `json:"tax" gorm:"column:tax;not null"`
	Total               float64         `json:"total" gorm:"column:total;not null"`
	DeliveryAddress     string          `json:"delivery_address" gorm:"column:delivery_address;not null"`
	PaymentMethod       string          `json:"payment_method" gorm:"column:payment_method;not null"`
	Status              OrderStatus     `json:"status" gorm:"column:status;not null;default:'pending'"`
	DeliveryPersonName  *string         `json:"delivery_person_name,omitempty" gorm:"column:delivery_person_name"`
	DeliveryPersonPhone *string         `json:"delivery_person_phone,omitempty" gorm:"column:delivery_person_phone"`
	TrackingURL         *string         `json:"tracking_url,omitempty" gorm:"column:tracking_url"`
	Notes               *string         `json:"notes,omitempty" gorm:"column:notes"`
	CreatedAt           time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time       `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeliveredAt         *time.Time      `json:"delivered_at,omitempty" gorm:"column:delivered_at"`
	User                User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Restaurant          Restaurant      `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}

// Cart represents cart entity (typically handled in memory/session)
type Cart struct {
	UserID     string  `json:"user_id"`
	Items      []Food  `json:"items"`
	TotalPrice float64 `json:"total_price"`
	ItemCount  int     `json:"item_count"`
}
