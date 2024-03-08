package models

import "time"

type Config struct {
	ID         int
	Title      string
	Content    string
	UpdateTime time.Time `gorm:"comment:'设置时间';type:timestamp;"`
}
