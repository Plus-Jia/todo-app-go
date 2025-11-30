package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// 用户模型
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open("sqlite3", "todo-app.db")
	if err != nil {
		fmt.Println("连接数据库失败:", err)
	}
	DB.AutoMigrate(&User{})
}
