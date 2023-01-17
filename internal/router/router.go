package router

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"rest-geoip/internal/maxmind"
	"sync"
	"text/template"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gopkg.in/unrolled/secure.v1"
)

var e *echo.Echo
var once sync.Once

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
	// TODO
	// limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter := tollbooth.NewLimiter(1, nil)
	e.Use(tollbooth_echo.LimitHandler(limiter))
}

// func initFrontend() {
// 	if !viper.GetBool("WEB") {
// 		return
// 	}
//
// 	initSecurityHeaders()
//
// 	pkger.Include("/templates")
// 	// Load HTML templates
// 	parsedTemplates, err := templating.ParseTemplates("/templates")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	t := &Template{
// 		templates: parsedTemplates,
// 	}
//
// 	e.Renderer = t
//
// 	frontend := e.Group("/")
// 	// frontend.GET("/", routes.DisplayGeoIPInfo)
// 	// frontend.GET("/:address", routes.DisplayGeoIPInfo)
// 	// frontend.POST("/", routes.SearchIPAddressInfo)
//
// 	// https://echo.labstack.com/middleware/static/
// 	// Serve static files via pkger's fs
// 	e.Use(middleware.Static(pkger.Dir("/static")))
// }

func ip(c echo.Context) error {
	type dto struct {
		IPAddress string `json:"ip_address"`
	}

	return c.JSON(http.StatusOK, dto{
		IPAddress: c.RealIP(),
	})
}

func geoip(c echo.Context) error {
	db := maxmind.GetInstance()
	r, _ := db.Lookup(net.ParseIP(c.RealIP()))

	return c.JSON(http.StatusOK, r)
}

func geoipForAddress(c echo.Context) error {
	db := maxmind.GetInstance()
	r, _ := db.Lookup(net.ParseIP(c.Param("ip_address")))

	return c.JSON(http.StatusOK, r)
}

func InitRouter() {
	once.Do(func() {
		e = echo.New()
	})

	listeningAddress := fmt.Sprintf("%s:%s", viper.GetString("LISTEN_ADDRESS"), viper.GetString("LISTEN_PORT"))

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
	api.GET("/geoip/:ip_address", geoipForAddress)
	// api.PUT("/update", validateAuth, updateDB)

	// initFrontend()

	e.Logger.Fatal(e.Start(listeningAddress))
}
