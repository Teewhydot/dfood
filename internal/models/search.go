package models

import (
	"time"
)

// RecentKeyword represents recent search keyword entity
type RecentKeyword struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string    `json:"userId" gorm:"column:user_id;not null;index"`
	Keyword   string    `json:"keyword" gorm:"column:keyword;not null;index"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SearchResult represents search result entity (for caching)
type SearchResult struct {
	Query       string      `json:"query"`
	Type        string      `json:"type"` // food, restaurant, mixed
	Results     interface{} `json:"results"`
	ResultCount int         `json:"resultCount"`
	Timestamp   time.Time   `json:"timestamp"`
}