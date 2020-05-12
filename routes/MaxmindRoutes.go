package routes

import (
	"rest-geoip/geoip"

	"github.com/gin-gonic/gin"
)

func GetIpAddressInfo(c *gin.Context) {
	var record geoip.MaxmindRecord
	var err error

	record, err = geoip.GetIpAddressInfo(c.ClientIP())
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(200, record)
}
