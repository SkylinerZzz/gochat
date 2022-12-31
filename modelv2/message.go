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
	ImageUrl   string `json:"image_url,omitempty"`
}

func (Message) TableName() string {
	return "chat_history"
}

// Validate return nil if validation passed, or return err
func (m Message) Validate() error {
	return validate.Struct(m)
}

func SaveMessage(m Message) error {
	return util.DB.Create(&m).Error
}
