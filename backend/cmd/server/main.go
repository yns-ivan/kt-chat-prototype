package main

import (
	"log"
	"os"

	"ktchat/backend/internal/chat"
	"ktchat/backend/internal/database"
	"ktchat/backend/internal/websocket"
	"ktchat/backend/pkg/config"
	"ktchat/backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize chat service
	chatService := chat.NewService(db, hub)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "ktchat-backend",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", chatService.Login)
			auth.POST("/register", chatService.Register)
			auth.POST("/refresh", chatService.RefreshToken)
		}

		// Chat routes (protected)
		chatRoutes := api.Group("/chat")
		chatRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			chatRoutes.GET("/rooms", chatService.GetRooms)
			chatRoutes.POST("/rooms", chatService.CreateRoom)
			chatRoutes.GET("/rooms/:roomID/messages", chatService.GetMessages)
			chatRoutes.POST("/rooms/:roomID/join", chatService.JoinRoom)
			chatRoutes.POST("/rooms/:roomID/leave", chatService.LeaveRoom)
		}

		// WebSocket endpoint
		api.GET("/ws", func(c *gin.Context) {
			websocket.HandleWebSocket(hub, c.Writer, c.Request)
		})
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 