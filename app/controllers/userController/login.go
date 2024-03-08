package userController

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"live-chat/app/apiException"
	"live-chat/app/midwares"
	"live-chat/app/services/userService"
	"live-chat/app/utils"
)

type LoginData struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录
func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//判断用户是否存在
	user, err := userService.IsUserExist(data.Account)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.UserExist)
		return
	}
	//判断密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.PasswordWrong)
		return
	}
	//生成token
	token, err := midwares.GenerateToken(user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var user_type string
	if user.Type == 0 {
		user_type = "普通用户"
	} else if user.Type == 3 {
		user_type = "管理员"
	}
	utils.JsonSuccessResponse(c, gin.H{
		"user_type": user_type,
		"token":     token,
	})
}
