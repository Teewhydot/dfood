package models

import (
	"time"
)

// FavoriteFood represents favorite food entity
type FavoriteFood struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"column:user_id;not null;index"`
	FoodID    string    `json:"food_id" gorm:"column:food_id;not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Food      Food      `json:"food,omitempty" gorm:"foreignKey:FoodID"`
}

// FavoriteRestaurant represents favorite restaurant entity
type FavoriteRestaurant struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       string     `json:"user_id" gorm:"column:user_id;not null;index"`
	RestaurantID string     `json:"restaurant_id" gorm:"column:restaurant_id;not null;index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	User         User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Restaurant   Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
}
