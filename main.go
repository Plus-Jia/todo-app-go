package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// 定义一个用户结构体，用于数据库操作
type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

var db *gorm.DB
var err error

// 初始化数据库
func initDatabase() {
	// 打开数据库
	db, err = gorm.Open("sqlite3", "todo-app.db")
	if err != nil {
		fmt.Println("数据库连接失败：", err)
	}
	// 自动迁移（创建表）
	db.AutoMigrate(&User{})
}
func main() {
	// 初始化数据库连接
	initDatabase()

	//创建 gin 路由
	r := gin.Default()

	// ping 路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动 web 服务
	r.Run(":8080")
}
