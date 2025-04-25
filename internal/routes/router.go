package routes

import (
	"ws-middle-cache/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) *gin.Engine {

	router.GET("/health", handlers.HealthCheck)
	router.GET("/liveness", handlers.LivenessCheck)

	api := router.Group("/api")
	{
		api.GET("/*any", handlers.CacheHandler)
	}

	return router
}
