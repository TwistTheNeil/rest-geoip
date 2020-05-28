package main

import (
	"rest-geoip/routes"

	"github.com/markbates/pkger"
)

func main() {
	// Explicitly include /templates
	pkger.Include("/templates")

	router := routes.SetupRouter()
	if err := router.Run(); err != nil {
		panic(err)
	}
}
