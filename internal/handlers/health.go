package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck responds with a status indicating the application is healthy.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}