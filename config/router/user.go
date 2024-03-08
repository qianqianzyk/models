package router

import (
	"github.com/gin-gonic/gin"
	"live-chat/app/controllers/userController"
	"live-chat/app/midwares"
)

func userRouterInit(r *gin.RouterGroup) {
	user := r.Group("/user").Use(midwares.JWTAuthMiddleware())
	{
		user.GET("/info", userController.GetUserInfo)
		user.PUT("/info", userController.UpdateUserInfo)
		user.POST("/sendCode", userController.SendCode)
		user.PUT("/emailVerify", userController.VerifyEmail)
		user.PUT("/findPassword", userController.FindPassword)
		user.POST("/img", userController.UploadImage)
	}
}
