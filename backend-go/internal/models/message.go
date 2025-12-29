package models

import "github.com/google/uuid"

type Message struct {
	BaseModel
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	GameRoomID  *uuid.UUID `gorm:"type:uuid;index" json:"game_room_id,omitempty"` // null = private message
	RecipientID *uuid.UUID `gorm:"type:uuid;index" json:"recipient_id,omitempty"` // for private messages
	Content     string     `gorm:"not null" json:"content"`
	MessageType string     `gorm:"default:text" json:"message_type"` // text, emoji, system

	// Relationships
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GameRoom  *GameRoom `gorm:"foreignKey:GameRoomID" json:"game_room,omitempty"`
	Recipient *User     `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
