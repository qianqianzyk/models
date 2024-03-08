package userController

import (
	"github.com/gin-gonic/gin"
	"live-chat/app/apiException"
	"live-chat/app/services/userService"
	"live-chat/app/utils"
)

type SendData struct {
	SendType int `json:"send_type" binding:"required"` //1: email ; 2: phone
}

func SendCode(c *gin.Context) {
	var data SendData
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
	if data.SendType == 1 && user.Email != "" {
		err = userService.SendMsgToEmail(user.Email)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	} else if data.SendType == 2 && user.Phone != "" {
		err = userService.SendMsgToPhone(user.Phone)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	} else {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type EmailCodeData struct {
	Code int `json:"code" binding:"required"`
}

func VerifyEmail(c *gin.Context) {
	var data EmailCodeData
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
	if user.EmailType == 1 && user.Email != "" {
		code, err := userService.GetVerificationCode(user.Email)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if code != data.Code {
			utils.JsonErrorResponse(c, apiException.CodeWrong)
			return
		}
		err = userService.VerifyEmail(user.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	} else {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type FindData struct {
	FindType int    `json:"find_type" binding:"required"` //1: email ; 2: phone
	Code     int    `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func FindPassword(c *gin.Context) {
	var data FindData
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
	if data.FindType == 1 && user.EmailType == 2 {
		code, err := userService.GetVerificationCode(user.Email)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if code != data.Code {
			utils.JsonErrorResponse(c, apiException.CodeWrong)
			return
		}
	} else if data.FindType == 2 {
		code, err := userService.GetVerificationCode(user.Phone)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if code != data.Code {
			utils.JsonErrorResponse(c, apiException.CodeWrong)
			return
		}
	} else {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	flag := userService.IsPasswordValid(data.Password)
	if !flag {
		utils.JsonErrorResponse(c, apiException.PasswordValid)
		return
	}
	err = userService.FindPassword(user.ID, data.Password)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
