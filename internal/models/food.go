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
	ImageURL        string      `json:"imageUrl" gorm:"column:image_url;not null"`
	Category        string      `json:"category" gorm:"column:category;not null;index"`
	RestaurantID    string      `json:"restaurantId" gorm:"column:restaurant_id;not null;index"`
	RestaurantName  string      `json:"restaurantName" gorm:"column:restaurant_name;not null"`
	Ingredients     StringArray `json:"ingredients" gorm:"column:ingredients;type:json"`
	IsAvailable     bool        `json:"isAvailable" gorm:"column:is_available;default:true"`
	PreparationTime string      `json:"preparationTime" gorm:"column:preparation_time"`
	Calories        int         `json:"calories" gorm:"column:calories;default:0"`
	Quantity        int         `json:"quantity" gorm:"column:quantity;default:1"`
	IsVegetarian    bool        `json:"isVegetarian" gorm:"column:is_vegetarian;default:false"`
	IsVegan         bool        `json:"isVegan" gorm:"column:is_vegan;default:false"`
	IsGlutenFree    bool        `json:"isGlutenFree" gorm:"column:is_gluten_free;default:false"`
	CreatedAt       time.Time   `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time   `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	Restaurant      Restaurant  `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}
