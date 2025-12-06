package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Plus-Jia/todo-app-go/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 解析 JWT 的秘钥
var jwtKey = []byte("your_secret_key")

// 从 token 里获取 user_id
func getUserIDFromToken(c *gin.Context) (uint, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return 0, jwt.ErrSignatureInvalid
	}

	// 支持带 "Bearer " 前缀的 Authorization header，移除前缀并去除空白
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["user_id"].(float64)), nil
	}

	return 0, jwt.ErrSignatureInvalid
}

// ------------------------------
// 1. 创建任务
// ------------------------------

func CreateTask(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 绑定用户
	task.UserID = userID
	task.CreatedAt = time.Now()

	if err := models.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// ------------------------------
// 2. 获取任务列表（当前用户）
// ------------------------------

func GetTasks(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var tasks []models.Task
	if err := models.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// ------------------------------
// 3. 更新任务
// ------------------------------

func UpdateTask(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	id := c.Param("id")

	var task models.Task
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 绑定更新数据
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Title = input.Title
	task.Desc = input.Desc
	task.Done = input.Done
	task.Deadline = input.Deadline
	task.UpdatedAt = time.Now()

	models.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// ------------------------------
// 4. 删除任务
// ------------------------------

func DeleteTask(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	id := c.Param("id")

	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// 高级任务列表：筛选、分页、排序
func GetTasksAdvanced(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 解析 query 参数
	status := c.Query("status") // all / completed / pending
	search := c.Query("search") // 搜索关键字
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sort := c.DefaultQuery("sort", "desc") // asc / desc

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构造查询
	var tasks []models.Task
	var total int64
	query := models.DB.Where("user_id = ?", userID)

	// 状态筛选
	switch status {
	case "completed":
		query = query.Where("done = ?", true)
	case "pending":
		query = query.Where("done = ?", false)
	}

	// 搜索（标题 + 内容）
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("title LIKE ? OR desc LIKE ?", like, like)
	}

	// 获取总数（用于分页）
	if err := query.Model(&models.Task{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Count failed"})
		return
	}

	// 排序
	if sort == "asc" {
		query = query.Order("created_at asc")
	} else {
		query = query.Order("created_at desc")
	}

	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// 执行查询
	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"tasks":     tasks,
	})
}
