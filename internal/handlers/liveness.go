package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LivenessCheck responds with a 200 OK status to indicate the application is alive.
func LivenessCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}