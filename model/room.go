package model

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	UserId   uint
	RoomName string
}

func AddRoom(value interface{}) Room {
	var r Room
	r.UserId = value.(map[string]interface{})["userId"].(uint)
	r.RoomName = value.(map[string]interface{})["roomName"].(string)
	ChatDB.Create(&r)
	return r
}
func ListAllRooms() []Room {
	var rs []Room
	ChatDB.Find(&rs)
	return rs
}
