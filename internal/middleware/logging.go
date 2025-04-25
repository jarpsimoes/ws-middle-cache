package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// LoggingMiddleware logs the incoming requests and their responses.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		duration := time.Since(startTime)
		log.Printf("Request: %s %s | Duration: %v | Status: %d",
			c.Request.Method, c.Request.URL.Path, duration, c.Writer.Status())
	}
}