package models

import (
	"time"
)

// User represents the main user entity - optimized for MariaDB
type User struct {
	ID              string    `json:"id" gorm:"primaryKey;column:id;type:varchar(36)"`
	FirstName       string    `json:"firstName" gorm:"column:first_name;type:varchar(100);not null"`
	LastName        string    `json:"lastName" gorm:"column:last_name;type:varchar(100);not null"`
	Email           string    `json:"email" gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	PhoneNumber     string    `json:"phoneNumber" gorm:"column:phone_number;type:varchar(20);not null;index"`
	Password        string    `json:"password,omitempty" gorm:"column:password;type:varchar(255);not null"`
	ProfileImageURL *string   `json:"profileImageUrl,omitempty" gorm:"column:profile_image_url;type:text"`
	Bio             *string   `json:"bio,omitempty" gorm:"column:bio;type:text"`
	FirstTimeLogin  bool      `json:"firstTimeLogin" gorm:"column:first_time_login;type:tinyint(1);default:1"`
	EmailVerified   bool      `json:"emailVerified" gorm:"column:email_verified;type:tinyint(1);default:0"`
	FCMToken        *string   `json:"fcmToken,omitempty" gorm:"column:fcm_token;type:text"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	AccessToken     string    `json:"access_token,omitempty" gorm:"-"`
	RefreshToken    string    `json:"refresh_token,omitempty" gorm:"-"`
}


// Address represents user address entity - optimized for MariaDB
type Address struct {
	ID        string    `json:"id" gorm:"primaryKey;column:id;type:varchar(36)"`
	UserID    string    `json:"userId" gorm:"column:user_id;type:varchar(36);not null;index"`
	Street    string    `json:"street" gorm:"column:street;type:varchar(255);not null"`
	City      string    `json:"city" gorm:"column:city;type:varchar(100);not null"`
	State     string    `json:"state" gorm:"column:state;type:varchar(100);not null"`
	ZipCode   string    `json:"zipCode" gorm:"column:zip_code;type:varchar(20);not null"`
	Type      string    `json:"type" gorm:"column:type;type:enum('home','work','other');default:'home'"`
	Address   string    `json:"address" gorm:"column:address;type:varchar(500);not null"`
	Apartment string    `json:"apartment" gorm:"column:apartment;type:varchar(100);not null"`
	Title     *string   `json:"title,omitempty" gorm:"column:title;type:varchar(100)"`
	Latitude  *float64  `json:"latitude,omitempty" gorm:"column:latitude;type:decimal(10,8)"`
	Longitude *float64  `json:"longitude,omitempty" gorm:"column:longitude;type:decimal(11,8)"`
	IsDefault bool      `json:"isDefault" gorm:"column:is_default;type:tinyint(1);default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Permission represents permission entity - optimized for MariaDB
type Permission struct {
	PermissionName string    `json:"permissionName" gorm:"primaryKey;column:permission_name;type:varchar(100)"`
	UserID         string    `json:"userId" gorm:"column:user_id;type:varchar(36);not null;index"`
	IsGranted      bool      `json:"isGranted" gorm:"column:is_granted;type:tinyint(1);not null"`
	LastUpdated    time.Time `json:"lastUpdated" gorm:"column:last_updated;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	User           User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// LocationData represents location information - optimized for MariaDB
type LocationData struct {
	Latitude  float64 `json:"latitude" gorm:"column:latitude;type:decimal(10,8);not null"`
	Longitude float64 `json:"longitude" gorm:"column:longitude;type:decimal(11,8);not null"`
	Address   string  `json:"address" gorm:"column:address;type:varchar(500);not null"`
	City      string  `json:"city" gorm:"column:city;type:varchar(100);not null"`
	Country   string  `json:"country" gorm:"column:country;type:varchar(100);not null"`
}
