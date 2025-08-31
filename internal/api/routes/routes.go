package routes

import (
	"matching_system/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", handlers.HealthCheck)

	matchHandler := handlers.NewMatchHandler()

	router.POST("/add-single-person-and-match", matchHandler.AddSinglePersonAndMatch)
	router.DELETE("/remove-single-person/:id", matchHandler.RemoveSinglePerson)
	router.GET("/query-single-people", matchHandler.QuerySinglePeople)

	return router
}
