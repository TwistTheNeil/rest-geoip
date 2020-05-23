package routes

import (
	"rest-geoip/templates"

	"github.com/gin-gonic/gin"
)

// SetupAndServe the gin router
func SetupAndServe() {
	router := gin.Default()
	t, err := templates.LoadTemplate()
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(t)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/geoip", GeoIPInfo)
		api.GET("/ip", IPAddress)
	}

	// Web routes
	router.GET("/", DisplayGeoIPInfo)

	router.Run() // #nosec
}
