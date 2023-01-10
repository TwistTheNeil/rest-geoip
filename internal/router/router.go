package router

import (
	"io"
	"rest-geoip/internal/templating"
	"rest-geoip/routes"
	"text/template"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"gopkg.in/unrolled/secure.v1"
)

var e *Echo

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func initSecurityHeaders() {
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

	e.Use(echo.WrapMiddleware(secureMiddleware.Handler))
}

func initRateLimiting() {
	limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	e.Use(tollbooth_echo.LimitHandler(limiter))
}

func initFrontend() {
	if !viper.GetBool("WEB") {
		return
	}

	initSecurityHeaders()

	pkger.Include("/templates")
	// Load HTML templates
	parsedTemplates, err := templating.ParseTemplates("/templates")
	if err != nil {
		panic(err)
	}

	t := &Template{
		templates: parsedTemplates,
	}

	e.Renderer = t

	frontend := e.Group("/")
	frontend.GET("/", routes.DisplayGeoIPInfo)
	frontend.GET("/:address", routes.DisplayGeoIPInfo)
	frontend.POST("/", routes.SearchIPAddressInfo)

	// https://echo.labstack.com/middleware/static/
	// Serve static files via pkger's fs
	e.Use(middleware.Static(pkger.Dir("/static")))
}

func InitRouter() {
	e = echo.New()

	if viper.GetBool("LOGGING") {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{viper.GetString("ALLOW_CORS_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	api := e.Group("/api")
	api.GET("/ip", ip)
	api.GET("/geoip", geoip)
	api.GET("/geoip/:address", address)
	api.PUT("/update", validateAuth, updateDB)

	initFrontend()

	e.Logger.Fatal(e.Start(":1323"))
}
