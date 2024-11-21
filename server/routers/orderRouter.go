package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/controllers"
)

func OrderRouter(router *gin.Engine) {
    orderRoutes := router.Group("/orders")
    {
        orderRoutes.POST("/", controllers.CreateOrderApi())
        orderRoutes.GET("/:order_id", controllers.GetOrderById())
        orderRoutes.PUT("/:order_id", controllers.UpdateOrder())
        orderRoutes.DELETE("/:order_id", controllers.DeleteOrder())
    }
}
