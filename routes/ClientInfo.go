package routes

import (
	"github.com/gin-gonic/gin"
)

// IPAddress replies with the client's ip address
func IPAddress(c *gin.Context) {
	c.JSON(200, gin.H{
		"IP": c.ClientIP(),
	})
}
