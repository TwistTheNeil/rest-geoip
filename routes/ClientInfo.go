package routes

import (
	"rest-geoip/geoip"

	"github.com/gin-gonic/gin"
)

// GeoIPInfo replies with the client's geo ip info
func GeoIPInfo(c *gin.Context) {
	record, err := geoip.Info(c.ClientIP())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, record)
}

// IPAddress replies with the client's ip address
func IPAddress(c *gin.Context) {
	c.JSON(200, gin.H{
		"IP": c.ClientIP(),
	})
}

// DisplayGeoIPInfo displays geoip info via HTML
func DisplayGeoIPInfo(c *gin.Context) {
	record, err := geoip.Info(c.ClientIP())

	// TODO: Show a custom 404 error on frontend
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.HTML(200, "ipAddressInfo.tmpl", record)
}
