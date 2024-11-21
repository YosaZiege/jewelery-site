package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/controllers"
	// "github.com/yosaZiege/jewelery-website/middleware"
)

func UserRouter(incomingRoutes *gin.Engine){
	// incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users" , controllers.GetUsers())
	incomingRoutes.GET("/users/id/:user_id" , controllers.GetUserById())
	incomingRoutes.GET("/users/:user_email" , controllers.GetUserByEmail())
}