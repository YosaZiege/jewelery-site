package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/config"
	"github.com/yosaZiege/jewelery-website/db"
	"github.com/yosaZiege/jewelery-website/routers"
)

func main() {
	// Initialize database connection
	db.InitDB()

	// Load environment variables and initialize Stripe
	config.LoadEnv()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	router := gin.Default()
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Apply routes from external router configurations
	routers.AuthRouter(router)
	routers.UserRouter(router)
	routers.ProductRouter(router)
	routers.CartRouter(router)
	routers.OrderRouter(router)
	routers.PaymentRoutes(router) // Register payment routes

	// Define unique API routes
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Access granted for API"})
	})
	router.GET("/api/v2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Access granted for v2 API"})
	})

	// Start server
	log.Fatal(router.Run(":" + port))
}
