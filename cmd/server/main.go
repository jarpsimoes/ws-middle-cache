package main

import (
	"os"
	"time"
	"ws-middle-cache/internal/handlers"
	"ws-middle-cache/internal/routes"
	"ws-middle-cache/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the shared in-memory cache
	cache := services.NewCache()

	// Start the cleanup task for expired cache items
	cache.StartCleanupTask(5 * time.Minute)

	// Pass the cache instance to handlers
	handlers.SetCacheInstance(cache)

	router := gin.Default()

	routes.SetupRouter(router)

	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	if err := router.Run(port); err != nil {
		panic(err)
	}
}
