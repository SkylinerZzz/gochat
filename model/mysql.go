package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ChatDB *gorm.DB

func init() {
	dsn := "root:123456@/chatdb?charset=utf8mb4&parseTime=true"
	ChatDB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
