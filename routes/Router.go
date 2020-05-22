package routes

import "github.com/gin-gonic/gin"

// SetupAndServe the gin router
func SetupAndServe() {
	r := gin.Default()
	r.GET("/geoip", GeoIPInfo)
	r.GET("/ip", IPAddress)

	r.Run()
}
