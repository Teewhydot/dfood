package models

import (
	"time"
)

// User represents the main user entity - SQLite compatible
type User struct {
	ID              string    `json:"id" gorm:"primaryKey;column:id"`
	FirstName       string    `json:"first_name" gorm:"column:first_name;not null"`
	LastName        string    `json:"last_name" gorm:"column:last_name;not null"`
	Email           string    `json:"email" gorm:"column:email;uniqueIndex;not null"`
	PhoneNumber     string    `json:"phone_number" gorm:"column:phone_number;not null;index"`
	Password        string    `json:"password,omitempty" gorm:"column:password;not null"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty" gorm:"column:profile_image_url"`
	Bio             *string   `json:"bio,omitempty" gorm:"column:bio"`
	FirstTimeLogin  bool      `json:"first_time_login" gorm:"column:first_time_login;default:1"`
	EmailVerified   bool      `json:"email_verified" gorm:"column:email_verified;default:0"`
	FCMToken        *string   `json:"fcm_token,omitempty" gorm:"column:fcm_token"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at"`
	AccessToken     string    `json:"access_token,omitempty" gorm:"-"`
	RefreshToken    string    `json:"refresh_token,omitempty" gorm:"-"`
}

// Address represents user address entity - SQLite compatible
type Address struct {
	ID        string    `json:"id" gorm:"primaryKey;column:id"`
	UserID    string    `json:"user_id" gorm:"column:user_id;not null;index"`
	Street    string    `json:"street" gorm:"column:street;not null"`
	City      string    `json:"city" gorm:"column:city;not null"`
	State     string    `json:"state" gorm:"column:state;not null"`
	ZipCode   string    `json:"zip_code" gorm:"column:zip_code;not null"`
	Type      string    `json:"type" gorm:"column:type;default:'home'"`
	Address   string    `json:"address" gorm:"column:address;not null"`
	Apartment string    `json:"apartment" gorm:"column:apartment;not null"`
	Title     *string   `json:"title,omitempty" gorm:"column:title"`
	Latitude  *float64  `json:"latitude,omitempty" gorm:"column:latitude"`
	Longitude *float64  `json:"longitude,omitempty" gorm:"column:longitude"`
	IsDefault bool      `json:"is_default" gorm:"column:is_default;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Permission represents permission entity - SQLite compatible
type Permission struct {
	PermissionName string    `json:"permission_name" gorm:"primaryKey;column:permission_name"`
	UserID         string    `json:"user_id" gorm:"column:user_id;not null;index"`
	IsGranted      bool      `json:"is_granted" gorm:"column:is_granted;not null"`
	LastUpdated    time.Time `json:"last_updated" gorm:"column:last_updated;autoUpdateTime"`
	User           User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// LocationData represents location information - SQLite compatible
type LocationData struct {
	Latitude  float64 `json:"latitude" gorm:"column:latitude;not null"`
	Longitude float64 `json:"longitude" gorm:"column:longitude;not null"`
	Address   string  `json:"address" gorm:"column:address;not null"`
	City      string  `json:"city" gorm:"column:city;not null"`
	Country   string  `json:"country" gorm:"column:country;not null"`
}
