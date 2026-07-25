package main

import (
	"website-downloader/handlers"
	"website-downloader/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"http://localhost:5176",
		"http://localhost:3000",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		api.POST("/scrape", handlers.StartScrapeHandler)
		api.GET("/status/:jobId", handlers.GetStatusHandler)
		api.GET("/download/:jobId", handlers.DownloadHandler)
		api.GET("/progress/:jobId", handlers.WebSocketHandler) // WebSocket endpoint
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	utils.LogInfo("Server starting on :8080")

	// Start server
	if err := router.Run(":8080"); err != nil {
		utils.LogError("Failed to start server: %v", err)
	}
}
