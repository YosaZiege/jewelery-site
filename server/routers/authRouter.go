package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/controllers"
)

func AuthRouter (incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup" , controllers.Signup())
	incomingRoutes.POST("users/signin" , controllers.Login())
}