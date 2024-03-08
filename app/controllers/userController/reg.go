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

type RegData struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
	var data RegData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//判断邮箱是否存在
	err = userService.IsEmailExist(data.Email)
	if err == nil {
		utils.JsonErrorResponse(c, apiException.EmailExist)
		return
	}
	//判断昵称是否存在
	err = userService.IsNicknameExist(data.Nickname)
	if err == nil {
		utils.JsonErrorResponse(c, apiException.NicknameExist)
		return
	}
	//判断昵称是否合法
	matched1, _ := regexp.MatchString(`^[\p{Han}a-zA-Z0-9_]+$`, data.Nickname)
	if !matched1 {
		utils.JsonErrorResponse(c, apiException.NicknameValid)
		return
	}
	//判断邮箱是否合法
	matched2, _ := regexp.MatchString(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`, data.Email)
	if !matched2 {
		utils.JsonErrorResponse(c, apiException.EmailValid)
		return
	}
	//判断密码是否符合规定
	//matched3, _ := regexp.MatchString(`^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`, data.Password)
	flag := userService.IsPasswordValid(data.Password)
	if !flag {
		utils.JsonErrorResponse(c, apiException.PasswordValid)
		return
	}
	//创建用户
	err = userService.CreateUser(models.User{
		Nickname:     data.Nickname,
		Email:        data.Email,
		EmailType:    1,
		Password:     data.Password,
		Type:         models.Person,
		Name:         "无",
		Phone:        "无",
		Introduction: "无",
		Avatar:       "无",
		CreateTime:   time.Now(),
	})
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
