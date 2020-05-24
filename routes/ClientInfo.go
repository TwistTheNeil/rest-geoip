package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IPAddress replies with the client's ip address
func IPAddress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"IP": c.ClientIP(),
	})
}
