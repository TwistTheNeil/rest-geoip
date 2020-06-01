package main

import (
	"rest-geoip/routes"
)

func main() {
	router := routes.SetupRouter()
	if err := router.Run(); err != nil {
		panic(err)
	}
}
