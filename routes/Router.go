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
		api.GET("/ip/:address", GeoIPInfo)
		api.POST("/update", UpdateMaxmindDB)
	}

	// Web routes
	web := router.Group("/web")
	{
		web.GET("/", DisplayGeoIPInfo)
		web.GET("/:address", DisplayGeoIPInfo)
		web.POST("/", SearchIPAddressInfo)
	}

	// Redirects
	router.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/web"
		router.HandleContext(c)
	})

	router.Run() // #nosec
}
