package main

import (
	"rest-geoip/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", routes.GetIpAddressInfo)

	r.Run()
}
