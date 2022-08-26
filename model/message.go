package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	UserId  uint
	RoomId  uint
	Content string
}

func SaveContent(value interface{}) Message {
	var m Message
	m.UserId = value.(map[string]interface{})["userId"].(uint)
	m.RoomId = value.(map[string]interface{})["roomId"].(uint)
	m.Content = value.(map[string]interface{})["content"].(string)
	ChatDB.Create(&m)
	return m
}
func ListRecordsByRoomId(roomId uint) []Message {
	var ms []Message
	ChatDB.Where("room_id = ?", roomId).Find(&ms)
	return ms
}
