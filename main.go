package main

import (
	"rest-geoip/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/geoip", routes.GeoIPInfo)
	r.GET("/ip", routes.IPAddress)

	r.Run()
}
