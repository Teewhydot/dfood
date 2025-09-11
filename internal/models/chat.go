package models

import (
	"time"
)

// Chat represents chat entity
type Chat struct {
	ID              string    `json:"id" gorm:"primaryKey;column:id"`
	SenderID        string    `json:"senderId" gorm:"column:sender_id;not null;index"`
	ReceiverID      string    `json:"receiverId" gorm:"column:receiver_id;not null;index"`
	OrderID         *string   `json:"orderId,omitempty" gorm:"column:order_id;index"`
	Name            string    `json:"name" gorm:"column:name;not null"`
	LastMessage     string    `json:"lastMessage" gorm:"column:last_message"`
	ImageURL        string    `json:"imageUrl" gorm:"column:image_url"`
	LastMessageTime time.Time `json:"lastMessageTime" gorm:"column:last_message_time;autoUpdateTime"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	Sender          User      `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Receiver        User      `json:"receiver,omitempty" gorm:"foreignKey:ReceiverID"`
	Order           *Order    `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Messages        []Message `json:"messages,omitempty" gorm:"foreignKey:ChatID"`
}

// Message represents message entity
type Message struct {
	ID          string    `json:"id" gorm:"primaryKey;column:id"`
	ChatID      string    `json:"chatId" gorm:"column:chat_id;not null;index"`
	Content     string    `json:"content" gorm:"column:content;not null"`
	SenderID    string    `json:"senderId" gorm:"column:sender_id;not null;index"`
	ReceiverID  string    `json:"receiverId" gorm:"column:receiver_id;not null;index"`
	IsRead      bool      `json:"isRead" gorm:"column:is_read;default:false"`
	MessageType string    `json:"messageType" gorm:"column:message_type;default:'text'"` // text, image, file
	FileURL     *string   `json:"fileUrl,omitempty" gorm:"column:file_url"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	Chat        Chat      `json:"chat,omitempty" gorm:"foreignKey:ChatID"`
	Sender      User      `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Receiver    User      `json:"receiver,omitempty" gorm:"foreignKey:ReceiverID"`
}
