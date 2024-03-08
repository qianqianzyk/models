package userController

import (
	"github.com/gin-gonic/gin"
	"live-chat/app/apiException"
	"live-chat/app/models"
	"live-chat/app/services/userService"
	"live-chat/app/utils"
	"regexp"
	"time"
)

type UserInfo struct {
	Name         string      `json:"name"`
	Nickname     string      `json:"nickname"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	Introduction interface{} `json:"introduction"`
	Avatar       interface{} `json:"avatar"`
}

// 修改个人信息
func UpdateUserInfo(c *gin.Context) {
	var data UserInfo
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	user, err := userService.GetUserByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var nickname = "无"
	if data.Nickname != user.Nickname && data.Nickname != "" {
		nickname = data.Nickname
		err = userService.IsNicknameExist(nickname)
		if err == nil {
			utils.JsonErrorResponse(c, apiException.NicknameExist)
			return
		}
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", nickname)
		if !matched {
			utils.JsonErrorResponse(c, apiException.NicknameValid)
			return
		}
	} else if data.Nickname == "" {
		nickname = user.Nickname
	}
	var email = "无"
	var emailType uint8
	if data.Email != user.Email && data.Email != "" {
		email = data.Email
		emailType = 1
		err = userService.IsEmailExist(email)
		if err == nil {
			utils.JsonErrorResponse(c, apiException.EmailExist)
			return
		}
		matched, _ := regexp.MatchString(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`, email)
		if !matched {
			utils.JsonErrorResponse(c, apiException.EmailValid)
			return
		}
	} else if data.Email == "" {
		email = user.Email
	}
	var phone = "无"
	if data.Phone != user.Phone && data.Phone != "" {
		phone = data.Phone
		err = userService.IsPhoneExist(phone)
		if err == nil {
			utils.JsonErrorResponse(c, apiException.PhoneExist)
			return
		}
		matched, _ := regexp.MatchString(`^1[3456789]\d{9}$`, phone)
		if !matched {
			utils.JsonErrorResponse(c, apiException.PhoneValid)
			return
		}
	} else if data.Phone == "" {
		phone = user.Phone
	}
	var introduction = "无"
	if data.Introduction != nil {
		introduction = data.Introduction.(string)
	}
	var avatar = "无"
	if data.Avatar != nil {
		avatar = data.Avatar.(string)
	}
	err = userService.UpdateUserInfoByUserID(user.ID, models.User{
		Nickname:     nickname,
		Email:        email,
		EmailType:    emailType,
		Name:         data.Name,
		Phone:        phone,
		Introduction: introduction,
		Avatar:       avatar,
	})
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type UserInfoData struct {
	Name         string          `json:"name" binding:"required"`
	Nickname     string          `json:"nickname" binding:"required"`
	Type         models.UserType `json:"type" binding:"required"`
	Email        string          `json:"email" binding:"required"`
	EmailType    uint8           `json:"email_type" binding:"required"`
	Phone        string          `json:"phone" binding:"required"`
	Introduction string          `json:"introduction" binding:"required"`
	Avatar       string          `json:"avatar" binding:"required"`
	CreateTime   time.Time       `json:"create_time" binding:"required"`
}

func GetUserInfo(c *gin.Context) {
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	user, err := userService.GetUserByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	response := UserInfoData{
		Name:         user.Name,
		Nickname:     user.Nickname,
		Type:         user.Type,
		Email:        user.Email,
		EmailType:    user.EmailType,
		Phone:        user.Phone,
		Introduction: user.Introduction,
		Avatar:       user.Avatar,
		CreateTime:   user.CreateTime,
	}

	utils.JsonSuccessResponse(c, response)
}
