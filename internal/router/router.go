package router

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"rest-geoip/internal/config"
	"rest-geoip/internal/maxmind"
	"strings"
	"sync"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	e    *echo.Echo
	once sync.Once
)

//go:embed all:dist
var spaFS embed.FS

//func initSecurityHeaders() {
//	secureMiddleware := secure.New(secure.Options{
//		STSSeconds:              31536000,
//		STSIncludeSubdomains:    true,
//		STSPreload:              false,
//		ForceSTSHeader:          true,
//		ContentTypeNosniff:      true,
//		BrowserXssFilter:        true,
//		ReferrerPolicy:          "same-origin",
//		FeaturePolicy:           "vibrate 'none';",
//		CustomFrameOptionsValue: "SAMEORIGIN",
//		FrameDeny:               true,
//		ContentSecurityPolicy:   "default-src 'self' https://api.mapbox.com 'unsafe-inline'; img-src 'self' https://api.mapbox.com data:",
//	})
//
//	e.Use(echo.WrapMiddleware(secureMiddleware.Handler))
//}

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

func cliAgentHander(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// If the request comes from curl or httpie, just send them back
		// a string with the ip in it.
		// Otherwise, show them the SPA frontend (handled by a separate
		// middleware)
		if c.Request().URL.String() != "/" {
			return next(c)
		}
		cliUserAgents := map[string]struct{}{
			"curl":   {},
			"HTTPie": {},
			"xh":     {},
		}

		ua := c.Request().UserAgent()
		uaName := strings.Split(ua, "/")

		_, isKnownCliUserAgent := cliUserAgents[uaName[0]]

		if isKnownCliUserAgent {
			return c.String(http.StatusOK, c.RealIP())
		}

		return next(c)
	}
}

func InitRouter() {
	once.Do(func() {
		e = echo.New()
	})

	// 3 req/s
	limiter := tollbooth.NewLimiter(3, nil)

	listeningAddress := fmt.Sprintf("%s:%s", config.Details().Program.ListenAddress, config.Details().Program.ListenPort)

	if config.Details().Program.EnableLogging {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Gzip())
	e.Use(cliAgentHander)

	if config.Details().Program.EnableWeb {
		if config.Details().Program.ReleaseMode {
			// SPA frontend handler
			e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
				Root:       "dist",       // This is the path to your SPA build folder, the folder that is created from running "npm build"
				Index:      "index.html", // This is the default html page for your SPA
				Browse:     false,
				HTML5:      true,
				Filesystem: http.FS(spaFS),
			}))
		} else {
			// Development mode - cors with vite dev
			e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: []string{"http://localhost:5173"},
				AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			}))
		}
	}

	api := e.Group("/api")
	// TODO: protect https://documentation.maptiler.com/hc/en-us/articles/360020806037-Protect-your-map-key
	api.GET("/config", func(c echo.Context) error {
		var dto struct {
			MaptilerToken string
			AdminNotice   string `json:",omitempty"`
		}

		dto.MaptilerToken = config.Details().Maptiler.Token
		dto.AdminNotice = config.Details().Program.AdminNotice
		return c.JSON(http.StatusOK, dto)
	})
	api.GET("/geoip", geoip)
	api.GET("/geoip/:ip_address", geoipForAddress)
	api.PUT("/update", func(c echo.Context) error {
		var dto struct {
			Message string `json:"message"`
		}
		err := maxmind.GetInstance().Update()
		if err != nil {
			fmt.Println(err)
			dto.Message = "error updating the maxmind database"
			return c.JSON(http.StatusInternalServerError, dto)
		}
		dto.Message = "success"
		return c.JSON(http.StatusOK, dto)
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		apikey := config.Details().Program.APIKey
		return func(c echo.Context) error {
			type authorization struct {
				API_KEY string `header:"x-api-key"`
			}
			var dto struct {
				Message string `json:"message"`
			}
			request := new(authorization)
			binder := &echo.DefaultBinder{}
			binder.BindHeaders(c, request)
			if apikey == request.API_KEY {
				return next(c)
			}
			dto.Message = "unauthorized"
			return c.JSON(http.StatusUnauthorized, dto)
		}
	})

	if config.Details().Program.EnableWeb {
		// We don't serve anything else, redirect to /
		e.Any("/*", func(c echo.Context) error {
			return c.Redirect(http.StatusPermanentRedirect, "/")
		})
	}

	e.Use(tollbooth_echo.LimitHandler(limiter))

	e.Logger.Fatal(e.Start(listeningAddress))
}
