package routes

import (
	"os"
	"rest-geoip/utils"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"gopkg.in/unrolled/secure.v1"
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
	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	// Limit API requests
	// No one really needs to request their deails every second
	r.Use(tollbooth_gin.LimitHandler(limiter))

	api := r.Group("/api")
	{
		api.GET("/ip", IPAddress)
		api.GET("/geoip", GeoIPInfo)
		api.GET("/geoip/:address", GeoIPInfo)
		api.POST("/update", UpdateMaxmindDB)
	}
}

func addSecurityHeaders(r *gin.Engine) {
	secureMiddleware := secure.New(secure.Options{
		STSSeconds:              31536000,
		STSIncludeSubdomains:    true,
		STSPreload:              false,
		ForceSTSHeader:          true,
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
		ReferrerPolicy:          "same-origin",
		FeaturePolicy:           "vibrate 'none';",
		CustomFrameOptionsValue: "SAMEORIGIN",
		FrameDeny:               true,
		ContentSecurityPolicy:   "default-src 'self' https://api.mapbox.com 'unsafe-inline'; img-src 'self' https://api.mapbox.com data:",
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	r.Use(secureFunc)
}

func setupWebRoutes(r *gin.Engine) {
	webEnv := os.Getenv("WEB")
	if !(webEnv == "true" || webEnv == "") {
		return
	}

	addSecurityHeaders(r)

	// Explicitly include /templates
	pkger.Include("/templates")
	// Load HTML templates
	t, err := utils.ParseTemplates("/templates")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

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
}

// SetupRouter returns a configured router
func SetupRouter() *gin.Engine {
	router := newRouter()
	setupWebRoutes(router)
	setupAPIRoutes(router)
	return router
}
