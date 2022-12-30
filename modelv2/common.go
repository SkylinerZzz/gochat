package modelv2

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gochat/util"
	"gorm.io/gorm"
)

// UserInfo describes user info from registration and login
type UserInfo struct {
	gorm.Model `json:"-"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func AddUser(u UserInfo) error {
	return util.DB.Create(&u).Error
}

func FindUserByName(username string) (UserInfo, error) {
	var u UserInfo
	err := util.DB.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u, nil
		}
		log.Errorf("[FindUserByName] failed to find user by name, err = %s", err)
		return u, err
	}
	return u, nil
}

// CheckUserExists returns false if username already exists, or return true
func CheckUserExists(username string) (bool, error) {
	err := util.DB.Where("username = ?", username).First(&UserInfo{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		log.Errorf("[CheckUserExists] failed to find user by name, err = %s", err)
		return false, err
	}
	return true, err
}

// RoomInfo describes room info while creating new room
type RoomInfo struct {
	gorm.Model
	RoomName string `form:"room_name"` // room name
	UserId   string `form:"user_id"`   // owner id
}

func (RoomInfo) TableName() string {
	return "room_info"
}

func AddRoom(r RoomInfo) error {
	return util.DB.Create(&r).Error
}

func FindRoomsByName(roomName string) ([]RoomInfo, error) {
	var r []RoomInfo
	err := util.DB.Where("room_name = ?", roomName).Find(&r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r, nil
		}
		log.Errorf("[FindRoomsByName] failed to find rooms by name, err = %s", err)
		return r, err
	}
	return r, nil
}
