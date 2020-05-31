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

func setupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/geoip", GeoIPInfo)
		api.GET("/ip", IPAddress)
		api.GET("/ip/:address", GeoIPInfo)
		api.POST("/update", UpdateMaxmindDB)
	}
}

func setupWebRoutes(r *gin.Engine) {
	webEnv := os.Getenv("WEB")
	if !(webEnv == "true" || webEnv == "") {
		return
	}

	web := r.Group("/web")
	{
		web.GET("/", DisplayGeoIPInfo)
		web.GET("/:address", DisplayGeoIPInfo)
		web.POST("/", SearchIPAddressInfo)
	}

	// Redirects
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/web"
		r.HandleContext(c)
	})

	// Serve static files via pkger's fs
	r.StaticFS("/static", pkger.Dir("/static"))

	// Load HTML templates
	t, err := utils.ParseTemplates("/templates")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
}

// SetupRouter returns a configured router
func SetupRouter() *gin.Engine {
	router := newRouter()
	setupAPIRoutes(router)
	setupWebRoutes(router)
	return router
}
