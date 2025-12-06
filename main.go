package main

import (
	"github.com/Plus-Jia/todo-app-go/controllers"
	"github.com/Plus-Jia/todo-app-go/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 定义一个用户结构体，用于数据库操作
type User_test struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func main() {
	// 初始化数据库连接
	models.InitDB()

	//创建 gin 路由
	r := gin.Default()

	// 启用 CORS，允许 Vite 前端 (http://localhost:5173) 发起请求并处理预检
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ping 路由 测试用
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//功能路由，待完成---
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	// task路由
	r.POST("/tasks", controllers.CreateTask)
	r.GET("/tasks", controllers.GetTasks)
	r.PUT("/tasks/:id", controllers.UpdateTask)
	r.DELETE("/tasks/:id", controllers.DeleteTask)
	r.GET("/advanced", controllers.GetTasksAdvanced)

	// 中间件安全路由测试
	// r.GET("/protected", TokenAuthMiddleware(), func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
	// })
	// 启动 web 服务
	r.Run(":8080")
}

// var jwtKey = []byte("your_secret_key")

// // JWT 验证中间件
// func TokenAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
// 			c.Abort()
// 			return
// 		}

// 		// 去掉 "Bearer " 前缀
// 		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

// 		// 解析 JWT
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// 确保 JWT 的签名方法是我们期望的
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// 将用户 ID 放入上下文
// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			c.Set("user_id", claims["user_id"])
// 		}

// 		c.Next()
// 	}
// }
