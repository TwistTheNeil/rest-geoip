package routes

import (
	"rest-geoip/maxmind"

	"net/http"

	"github.com/gin-gonic/gin"
)

// GeoIPInfo replies with the client's geo ip info
func GeoIPInfo(c *gin.Context) {
	record, err := maxmind.Info(c.ClientIP())
	// TODO: use service unavaialble if db is not present
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DisplayGeoIPInfo displays geoip info via HTML
func DisplayGeoIPInfo(c *gin.Context) {
	record, err := maxmind.Info(c.ClientIP())

	// TODO: Show a custom 404 error on frontend
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "/templates/ipAddressInfo.tmpl", record)
}

// UpdateMaxmindDB updates the db
func UpdateMaxmindDB(c *gin.Context) {
	if err := maxmind.DownloadAndUpdate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
