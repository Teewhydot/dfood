package models

import (
	"time"
)

// PaymentMethod represents payment method entity
type PaymentMethod struct {
	ID      string `json:"id" gorm:"primaryKey;column:id"`
	Name    string `json:"name" gorm:"column:name;not null"`
	Type    string `json:"type" gorm:"column:type;not null"`
	IconURL string `json:"icon_url" gorm:"column:icon_url;not null"`
}

// Card represents payment card entity
type Card struct {
	ID              string        `json:"id" gorm:"primaryKey;column:id"`
	UserID          string        `json:"user_id" gorm:"column:user_id;not null;index"`
	PaymentMethodID string        `json:"payment_method_id" gorm:"column:payment_method_id;not null"`
	PAN             string        `json:"pan" gorm:"column:pan;not null"` // Should be encrypted
	CVV             string        `json:"cvv" gorm:"column:cvv;not null"` // Should be encrypted
	ExpiryMonth     int           `json:"expiry_month" gorm:"column:expiry_month;not null"`
	ExpiryYear      int           `json:"expiry_year" gorm:"column:expiry_year;not null"`
	CardholderName  string        `json:"cardholder_name" gorm:"column:cardholder_name;not null"`
	IsDefault       bool          `json:"is_default" gorm:"column:is_default;default:false"`
	CreatedAt       time.Time     `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time     `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	User            User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PaymentMethod   PaymentMethod `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`
}

// PaymentTransaction represents payment transaction entity
type PaymentTransaction struct {
	ID              string     `json:"id" gorm:"primaryKey;column:id"`
	OrderID         string     `json:"order_id" gorm:"column:order_id;not null;index"`
	UserID          string     `json:"user_id" gorm:"column:user_id;not null;index"`
	PaymentMethodID string     `json:"payment_method_id" gorm:"column:payment_method_id;not null"`
	Amount          float64    `json:"amount" gorm:"column:amount;not null"`
	Currency        string     `json:"currency" gorm:"column:currency;default:'USD';not null"`
	Status          string     `json:"status" gorm:"column:status;not null"`                  // pending, completed, failed, refunded
	TransactionID   *string    `json:"transaction_id,omitempty" gorm:"column:transaction_id"` // External payment gateway transaction ID
	FailureReason   *string    `json:"failure_reason,omitempty" gorm:"column:failure_reason"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty" gorm:"column:processed_at"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Order           Order      `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	User            User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
