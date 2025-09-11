package models

import (
	"time"
)

// Notification represents notification entity
type Notification struct {
	ID        string    `json:"id" gorm:"primaryKey;column:id"`
	UserID    string    `json:"userId" gorm:"column:user_id;not null;index"`
	Title     string    `json:"title" gorm:"column:title;not null"`
	Body      string    `json:"body" gorm:"column:body;not null"`
	Type      string    `json:"type" gorm:"column:type;not null"` // order, promotion, system, chat
	IsRead    bool      `json:"isRead" gorm:"column:is_read;default:false"`
	Data      *string   `json:"data,omitempty" gorm:"column:data;type:json"` // Additional data as JSON
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
