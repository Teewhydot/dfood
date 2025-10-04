package models

import (
	"time"
)

// Restaurant represents the restaurant entity - SQLite compatible
type Restaurant struct {
	ID             string                   `json:"id" gorm:"primaryKey;column:id"`
	Name           string                   `json:"name" gorm:"column:name;not null;index"`
	Description    string                   `json:"description" gorm:"column:description;not null"`
	Location       string                   `json:"location" gorm:"column:location;not null"`
	Distance       float64                  `json:"distance" gorm:"column:distance;default:0.0"`
	Rating         float64                  `json:"rating" gorm:"column:rating;default:0.0;index"`
	DeliveryTime   string                   `json:"delivery_time" gorm:"column:delivery_time;not null"`
	DeliveryFee    float64                  `json:"delivery_fee" gorm:"column:delivery_fee;not null"`
	ImageURL       string                   `json:"image_url" gorm:"column:image_url;not null"`
	Categories     StringArray              `json:"categories" gorm:"column:categories"`
	IsOpen         bool                     `json:"is_open" gorm:"column:is_open;default:true;index"`
	Latitude       float64                  `json:"latitude" gorm:"column:latitude;not null"`
	Longitude      float64                  `json:"longitude" gorm:"column:longitude;not null"`
	CreatedAt      time.Time                `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time                `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Foods          []Food                   `json:"foods,omitempty" gorm:"foreignKey:RestaurantID"`
	FoodCategories []RestaurantFoodCategory `json:"food_categories,omitempty" gorm:"foreignKey:RestaurantID"`
}

// RestaurantFoodCategory represents food categories within a restaurant
type RestaurantFoodCategory struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	RestaurantID string     `json:"restaurant_id" gorm:"column:restaurant_id;not null;index"`
	Category     string     `json:"category" gorm:"column:category;not null"`
	ImageURL     string     `json:"image_url" gorm:"column:image_url;not null"`
	Restaurant   Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Foods        []Food     `json:"foods,omitempty" gorm:"foreignKey:RestaurantID,Category;references:RestaurantID,Category"`
}
