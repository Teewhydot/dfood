package models

import (
	"time"
)

// PaymentMethod represents payment method entity
type PaymentMethod struct {
	ID      string `json:"id" gorm:"primaryKey;column:id"`
	Name    string `json:"name" gorm:"column:name;not null"`
	Type    string `json:"type" gorm:"column:type;not null"`
	IconURL string `json:"iconUrl" gorm:"column:icon_url;not null"`
}

// Card represents payment card entity
type Card struct {
	ID              string        `json:"id" gorm:"primaryKey;column:id"`
	UserID          string        `json:"userId" gorm:"column:user_id;not null;index"`
	PaymentMethodID string        `json:"paymentMethodId" gorm:"column:payment_method_id;not null"`
	PAN             string        `json:"pan" gorm:"column:pan;not null"` // Should be encrypted
	CVV             string        `json:"cvv" gorm:"column:cvv;not null"` // Should be encrypted
	ExpiryMonth     int           `json:"mExp" gorm:"column:expiry_month;not null"`
	ExpiryYear      int           `json:"yExp" gorm:"column:expiry_year;not null"`
	CardholderName  string        `json:"cardholderName" gorm:"column:cardholder_name;not null"`
	IsDefault       bool          `json:"isDefault" gorm:"column:is_default;default:false"`
	CreatedAt       time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	User            User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PaymentMethod   PaymentMethod `json:"paymentMethod,omitempty" gorm:"foreignKey:PaymentMethodID"`
}

// PaymentTransaction represents payment transaction entity
type PaymentTransaction struct {
	ID              string     `json:"id" gorm:"primaryKey;column:id"`
	OrderID         string     `json:"orderId" gorm:"column:order_id;not null;index"`
	UserID          string     `json:"userId" gorm:"column:user_id;not null;index"`
	PaymentMethodID string     `json:"paymentMethodId" gorm:"column:payment_method_id;not null"`
	Amount          float64    `json:"amount" gorm:"column:amount;not null"`
	Currency        string     `json:"currency" gorm:"column:currency;default:'USD';not null"`
	Status          string     `json:"status" gorm:"column:status;not null"`                 // pending, completed, failed, refunded
	TransactionID   *string    `json:"transactionId,omitempty" gorm:"column:transaction_id"` // External payment gateway transaction ID
	FailureReason   *string    `json:"failureReason,omitempty" gorm:"column:failure_reason"`
	ProcessedAt     *time.Time `json:"processedAt,omitempty" gorm:"column:processed_at"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	Order           Order      `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	User            User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
