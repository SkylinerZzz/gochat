package modelv2

import (
	"gochat/util"
	"gorm.io/gorm"
)

// Message definition
type Message struct {
	gorm.Model
	UserId      string `json:"user_id" validate:"required"`
	RoomId      string `json:"room_id" validate:"required"`
	Content     string `json:"content,omitempty"`
	ImageUrl    string `json:"image_url,omitempty"`
	ImageBase64 string `json:"image_base64,omitempty"`
}

func (Message) TableName() string {
	return "chat_history"
}

func SaveMessage(m Message) error {
	return util.DB.Create(&m).Error
}
