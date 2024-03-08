package router

import (
	"github.com/gin-gonic/gin"
	"live-chat/app/controllers/userController"
)

func Init(r *gin.Engine) {

	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/user/reg", userController.Register)
		api.POST("/user/login", userController.Login)
		userRouterInit(api)
	}
}
