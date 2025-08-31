package routes

import (
	"matching_system/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", handlers.HealthCheck)

	return router
}
