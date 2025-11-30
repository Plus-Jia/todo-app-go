package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Plus-Jia/todo-app-go/controllers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// 定义一个用户结构体，用于数据库操作
type User_test struct {
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
	db.AutoMigrate(&User_test{})
}
func main() {
	// 初始化数据库连接
	initDatabase()

	//创建 gin 路由
	r := gin.Default()

	// ping 路由 测试用
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//功能路由，待完成---
	r.POST("/api/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// 启动 web 服务
	r.Run(":8080")
}

var jwtKey = []byte("your_secret_key")

// JWT 验证中间件
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// 去掉 "Bearer " 前缀
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// 解析 JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 确保 JWT 的签名方法是我们期望的
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户 ID 放入上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
