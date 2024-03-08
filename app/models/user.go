package models

import "time"

type User struct {
	ID           int       `json:"id"`                     // 登录编号
	Name         string    `json:"name"`                   //姓名
	Nickname     string    `json:"nickname" gorm:"unique"` //昵称
	Type         UserType  `json:"type"`
	Password     string    `json:"password"`     // 密码
	Email        string    `json:"email"`        //邮箱
	EmailType    uint8     `json:"email_type"`   //1: 未验证 ; 2: 已验证
	Phone        string    `json:"phone"`        //电话
	Introduction string    `json:"introduction"` //自我简介
	Avatar       string    `json:"avatar"`       //头像
	CreateTime   time.Time `json:"create_time"`
}

type UserType int

const (
	Person UserType = 0
	Admin  UserType = 3
)
