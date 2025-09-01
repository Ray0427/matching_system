package routes

import (
	"matching_system/internal/api/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", handlers.HealthCheck)

	matchHandler := handlers.NewMatchHandler()

	router.POST("/add-single-person-and-match", matchHandler.AddSinglePersonAndMatch)
	router.DELETE("/remove-single-person/:id", matchHandler.RemoveSinglePerson)
	router.GET("/query-single-people", matchHandler.QuerySinglePeople)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
