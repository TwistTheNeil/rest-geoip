package routes

import (
	"os"
	"rest-geoip/utils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
)

func newRouter() *gin.Engine {
	logging := os.Getenv("LOGGING")

	if logging == "" || logging == "true" {
		return gin.Default()
	}

	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

// SetupRouter returns a configured router
func SetupRouter() *gin.Engine {
	router := newRouter()
	t, err := utils.ParseTemplates("/templates")
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

	// Serve static files via pkger's fs
	router.StaticFS("/static", pkger.Dir("/static"))

	return router
}
