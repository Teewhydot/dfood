package models

import (
	"time"
)

// Restaurant represents the restaurant entity - optimized for MariaDB
type Restaurant struct {
	ID             string                   `json:"id" gorm:"primaryKey;column:id;type:varchar(36)"`
	Name           string                   `json:"name" gorm:"column:name;type:varchar(255);not null;index"`
	Description    string                   `json:"description" gorm:"column:description;type:text;not null"`
	Location       string                   `json:"location" gorm:"column:location;type:varchar(500);not null"`
	Distance       float64                  `json:"distance" gorm:"column:distance;type:decimal(8,2);default:0.00"`
	Rating         float64                  `json:"rating" gorm:"column:rating;type:decimal(3,2);default:0.00;index"`
	DeliveryTime   string                   `json:"deliveryTime" gorm:"column:delivery_time;type:varchar(50);not null"`
	DeliveryFee    float64                  `json:"deliveryFee" gorm:"column:delivery_fee;type:decimal(8,2);not null"`
	ImageURL       string                   `json:"imageUrl" gorm:"column:image_url;type:text;not null"`
	Categories     StringArray              `json:"categories" gorm:"column:categories;type:json"`
	IsOpen         bool                     `json:"isOpen" gorm:"column:is_open;type:tinyint(1);default:1;index"`
	Latitude       float64                  `json:"latitude" gorm:"column:latitude;type:decimal(10,8);not null"`
	Longitude      float64                  `json:"longitude" gorm:"column:longitude;type:decimal(11,8);not null"`
	CreatedAt      time.Time                `json:"createdAt" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time                `json:"updatedAt" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Foods          []Food                   `json:"foods,omitempty" gorm:"foreignKey:RestaurantID"`
	FoodCategories []RestaurantFoodCategory `json:"foodCategories,omitempty" gorm:"foreignKey:RestaurantID"`
}

// RestaurantFoodCategory represents food categories within a restaurant
type RestaurantFoodCategory struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	RestaurantID string     `json:"restaurantId" gorm:"column:restaurant_id;not null;index"`
	Category     string     `json:"category" gorm:"column:category;not null"`
	ImageURL     string     `json:"imageUrl" gorm:"column:image_url;not null"`
	Restaurant   Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Foods        []Food     `json:"foods,omitempty" gorm:"foreignKey:RestaurantID,Category;references:RestaurantID,Category"`
}
