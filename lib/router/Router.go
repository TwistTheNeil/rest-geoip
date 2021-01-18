package router

import (
	"fmt"
	"net/http"
	"rest-geoip/lib/utils"
	"rest-geoip/routes"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"gopkg.in/unrolled/secure.v1"
)

func newRouter() (*gin.Engine, *http.Server) {
	listenAddress := viper.GetString("LISTEN_ADDRESS") + ":" + viper.GetString("LISTEN_PORT")
	var r *gin.Engine

	if !viper.GetBool("LOGGING") {
		r = gin.Default()
	} else {
		r = gin.New()
		r.Use(gin.Recovery())
	}

	s := &http.Server{
		Addr:           listenAddress,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return r, s
}

func validateAccess(c *gin.Context) {
	actualKey := c.Request.Header.Get("X-API-KEY")
	expectedKey := viper.GetString("API_KEY")
	if actualKey == expectedKey {
		c.Next()
		return
	}
	c.AbortWithStatus(401)
}

func setupAPIRoutes(r *gin.Engine) {
	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	// Limit API requests
	// No one really needs to request their deails every second
	r.Use(tollbooth_gin.LimitHandler(limiter))

	api := r.Group("/api")
	{
		api.GET("/ip", routes.IPAddress)
		api.GET("/geoip", routes.GeoIPInfo)
		api.GET("/geoip/:address", routes.GeoIPInfo)
		api.POST("/update", validateAccess, routes.UpdateMaxmindDB)
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
	if !viper.GetBool("WEB") {
		return
	}

	addSecurityHeaders(r)

	// Explicitly include /templates
	pkger.Include("/templates") //nolint
	// Load HTML templates
	t, err := utils.ParseTemplates("/templates")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	web := r.Group("/web")
	{
		web.GET("/", routes.DisplayGeoIPInfo)
		web.GET("/:address", routes.DisplayGeoIPInfo)
		web.POST("/", routes.SearchIPAddressInfo)
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
func SetupRouter() *http.Server {
	// Set up an API key if there is none set
	if viper.GetString("API_KEY") == "" {
		key, err := utils.GenerateKey(512)
		if err != nil {
			fmt.Println(err)
		} else {
			viper.Set("API_KEY", key)
			fmt.Printf("API KEY: %s\n", key)
		}
	}

	router, server := newRouter()
	setupWebRoutes(router)
	setupAPIRoutes(router)
	return server
}
