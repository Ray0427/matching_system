package main

import (
	"log"
	"matching_system/internal/api/routes"
	"matching_system/internal/config"
	"matching_system/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @title Matching System API
// @version 1.0
// @description This is a matching system API server with in-memory storage
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := routes.Setup()

	// Start server
	logger.Info("Starting server on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
