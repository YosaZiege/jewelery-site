package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/controllers"
)



func CartRouter(router *gin.Engine) {
	// incomingRoutes.Use(middleware.Authenticate())
	cartGroup := router.Group("/cart")
	{
		cartGroup.POST("/add" , controllers.AddProductToCart())
		cartGroup.GET("/:user_id" , controllers.ViewCartProducts())
		cartGroup.PUT("/update" , controllers.UpdateQuantity())
		cartGroup.DELETE("/remove" , controllers.RemoveProduct())
	}
}