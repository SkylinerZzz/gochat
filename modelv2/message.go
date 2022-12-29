package modelv2

import (
	"gochat/util"
	"gorm.io/gorm"
)

// Message definition
type Message struct {
	gorm.Model `json:"-"`
	UserId     string `json:"user_id" validate:"required"`
	Username   string `json:"username" validate:"required" gorm:"-"`
	RoomId     string `json:"room_id" validate:"required"`
	Content    string `json:"content,omitempty"`
	ImageId    string `json:"image_id,omitempty"`
}

func (Message) TableName() string {
	return "chat_history"
}

func SaveMessage(m Message) error {
	return util.DB.Create(&m).Error
}
