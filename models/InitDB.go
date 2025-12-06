package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	db, err := gorm.Open(sqlite.Open("todo-app.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败:", err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Task{})

	DB = db.Debug()
}
