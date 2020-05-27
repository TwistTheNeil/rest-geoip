package routes

import (
	"rest-geoip/maxmind"

	"net/http"

	"github.com/gin-gonic/gin"
)

func getRelevantIPAddress(c *gin.Context) string {
	ipQuery := c.ClientIP()
	if c.Param("address") != "" {
		ipQuery = c.Param("address")
	}

	return ipQuery
}

// GeoIPInfo replies with the client's geo ip info
func GeoIPInfo(c *gin.Context) {
	record, err := maxmind.Info(getRelevantIPAddress(c))
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
	record, err := maxmind.Info(getRelevantIPAddress(c))

	// TODO: Show a custom 404 error on frontend
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "ipAddressInfo", record)
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

// SearchIPAddressInfo constructs the correct URL for the posted IP Address and redirects to it
func SearchIPAddressInfo(c *gin.Context) {
	type searchForm struct {
		IPAddress string `form:"search-input"`
	}

	var form searchForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
	}
	c.Redirect(http.StatusFound, "/web/"+form.IPAddress)
}
