package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlePing responds to heartbeat pings
func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}
