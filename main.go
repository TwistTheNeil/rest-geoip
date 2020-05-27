package main

import (
	"rest-geoip/routes"

	"github.com/markbates/pkger"
)

func main() {
	// Explicitly include /templates
	pkger.Include("/templates")
	routes.SetupAndServe()
}
