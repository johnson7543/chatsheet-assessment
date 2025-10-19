package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnson7543/chatsheet-assessment/internal/config"
	"github.com/johnson7543/chatsheet-assessment/internal/database"
	"github.com/johnson7543/chatsheet-assessment/internal/handlers"
	"github.com/johnson7543/chatsheet-assessment/internal/middleware"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := config.GetConfig()

	// Initialize database
	if err := database.InitDatabase(cfg.DatabasePath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS from config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.FrontendURL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "LinkedIn Connector API is running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// LinkedIn connection routes
			linkedin := protected.Group("/linkedin")
			{
				linkedin.POST("/connect/cookie", handlers.ConnectLinkedInWithCookie)
				linkedin.POST("/connect/credentials", handlers.ConnectLinkedInWithCredentials)
			}

			// Account management routes
			protected.GET("/accounts", handlers.GetAccounts)
			protected.DELETE("/accounts/:id", handlers.DeleteAccount)
		}
	}

	// Start server
	log.Printf("Server starting on port %d", cfg.Server.Port)

	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
