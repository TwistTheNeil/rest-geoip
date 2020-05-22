package routes

import "github.com/gin-gonic/gin"

// SetupAndServe the gin router
func SetupAndServe() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// API routes
	api := router.Group("/api")
	{
		api.GET("/geoip", GeoIPInfo)
		api.GET("/ip", IPAddress)
	}

	// Web routes
	router.GET("/", DisplayGeoIPInfo)

	router.Run()
}
