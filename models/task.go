package models

import "time"

type Task struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	Deadline time.Time `json:"deadline"`
	Done     bool      `json:"done"`
	UserID   uint      `json:"user_id"` // 外键
}
