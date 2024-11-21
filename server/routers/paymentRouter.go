package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/controllers"
)

func PaymentRoutes(router *gin.Engine) {
	router.POST("/create-payment-intent", controllers.CreatePaymentIntent)
}