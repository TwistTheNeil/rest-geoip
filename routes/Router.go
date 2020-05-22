package routes

import "github.com/gin-gonic/gin"

// SetupAndServe the gin router
func SetupAndServe() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/geoip", GeoIPInfo)
		api.GET("/ip", IPAddress)
	}

	router.Run()
}
