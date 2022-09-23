package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

var (
	loadCacheOnce = sync.Once{} // only load rooms from database to cache once
)

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
	// synchronize cache
	AddRoomToCache(r)
	logrus.Info("synchronize cache")
	return r
}
func ListAllRooms() []Room {
	var rs []Room
	// load rooms once
	loadCacheOnce.Do(func() {
		ChatDB.Find(&rs)
		for _, r := range rs {
			AddRoomToCache(r)
		}
		logrus.Info("load rooms from database into cache")
	})
	rs = ListAllRoomsFromCache()
	logrus.Info("hit cache")
	return rs
}
