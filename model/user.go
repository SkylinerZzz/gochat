package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
}

func AddUser(value interface{}) User {
	var u User
	u.Username = value.(map[string]interface{})["username"].(string)
	u.Password = value.(map[string]interface{})["password"].(string)
	ChatDB.Create(&u)
	return u
}

func FindUserById(id uint) User {
	var u User
	ChatDB.Where("id = ?", id).First(&u)
	return u
}

func FindUserByName(name string) User {
	var u User
	ChatDB.Where("username = ?", name).First(&u)
	return u
}
