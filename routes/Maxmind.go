package routes

import (
	"net"
	"rest-geoip/customerrors"
	"rest-geoip/maxmind"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func relevantIPAddress(c *gin.Context) (net.IP, error) {
	ipQuery := c.ClientIP()
	if c.Param("address") != "" {
		ipQuery = c.Param("address")
	}

	ip := net.ParseIP(ipQuery)

	if ip == nil {
		return nil, customerrors.ErrInvalidIPAddress
	}

	return ip, nil
}

// GeoIPInfo replies with the client's geo ip info
func GeoIPInfo(c *gin.Context) {
	ip, err := relevantIPAddress(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": err.Error(),
		})
		return
	}

	record, err := maxmind.Info(ip)
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
	ip, err := relevantIPAddress(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": err.Error(),
		})
		return
	}

	record, err := maxmind.Info(ip)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "ipAddressInfo", gin.H{
		"record":              record,
		"mapbox_access_token": viper.GetString("MAPBOX_ACCESS_TOKEN"),
	})
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
