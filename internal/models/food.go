package models

import (
	"time"
)

// Food represents the food entity
type Food struct {
	ID              string      `json:"id" gorm:"primaryKey;column:id"`
	Name            string      `json:"name" gorm:"column:name;not null;index"`
	Description     string      `json:"description" gorm:"column:description;not null"`
	Price           float64     `json:"price" gorm:"column:price;not null"`
	Rating          float64     `json:"rating" gorm:"column:rating;default:0.0"`
	ImageURL        string      `json:"image_url" gorm:"column:image_url;not null"`
	Category        string      `json:"category" gorm:"column:category;not null;index"`
	RestaurantID    string      `json:"restaurant_id" gorm:"column:restaurant_id;not null;index"`
	RestaurantName  string      `json:"restaurant_name" gorm:"column:restaurant_name;not null"`
	Ingredients     StringArray `json:"ingredients" gorm:"column:ingredients"`
	IsAvailable     bool        `json:"is_available" gorm:"column:is_available;default:true"`
	PreparationTime string      `json:"preparation_time" gorm:"column:preparation_time"`
	Calories        int         `json:"calories" gorm:"column:calories;default:0"`
	Quantity        int         `json:"quantity" gorm:"column:quantity;default:1"`
	IsVegetarian    bool        `json:"is_vegetarian" gorm:"column:is_vegetarian;default:false"`
	IsVegan         bool        `json:"is_vegan" gorm:"column:is_vegan;default:false"`
	IsGlutenFree    bool        `json:"is_gluten_free" gorm:"column:is_gluten_free;default:false"`
	CreatedAt       time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Restaurant      Restaurant  `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}
