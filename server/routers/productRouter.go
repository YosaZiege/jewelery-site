package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/controllers"
	// "github.com/yosaZiege/jewelery-website/middleware"
)



func ProductRouter(router *gin.Engine) {
	// incomingRoutes.Use(middleware.Authenticate())
	productGroup := router.Group("/products") 
	{
		productGroup.GET("/", controllers.GetAllProducts())
		productGroup.GET("/productpage/:product_id" , controllers.FetchProductPageDetails())
		productGroup.GET("/:product_id", controllers.GetProductByIdApi())
		productGroup.POST("/", controllers.AddProduct())
		productGroup.PUT("/:product_id", controllers.UpdateProduct())
		productGroup.DELETE("/:product_id", controllers.DeleteProduct())
		productGroup.GET("/best-sellers", controllers.BestSellingProducts())
	}	
	
}