package routes

import (
	"rest-geoip/maxmind"

	"github.com/gin-gonic/gin"
)

// GeoIPInfo replies with the client's geo ip info
func GeoIPInfo(c *gin.Context) {
	record, err := maxmind.Info(c.ClientIP())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, record)
}

// DisplayGeoIPInfo displays geoip info via HTML
func DisplayGeoIPInfo(c *gin.Context) {
	record, err := maxmind.Info(c.ClientIP())

	// TODO: Show a custom 404 error on frontend
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(200, "/templates/ipAddressInfo.tmpl", record)
}

// UpdateMaxmindDB updates the db
func UpdateMaxmindDB(c *gin.Context) {
	if err := maxmind.DownloadAndUpdate(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{})
}
