package main

import (
	"log"

	"auth-backend/internal/config"
	"auth-backend/internal/database"
	"auth-backend/internal/handlers"
	"auth-backend/internal/middleware"
	"auth-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize services
	emailService := services.NewEmailService(cfg.EmailWebhookURL)
	authService := services.NewAuthService(db, emailService, cfg.JWTSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/request-verification", authHandler.RequestVerification)
		auth.POST("/verify-code", authHandler.VerifyCode)
		auth.POST("/set-password", authHandler.SetPassword)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/profile", authHandler.GetProfile)
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
