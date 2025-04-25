package utils

import (
	"github.com/gin-gonic/gin"
)

// JSONResponse formats a JSON response with a given status code and data.
func JSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}

// ErrorResponse formats a JSON error response with a given status code and message.
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}
