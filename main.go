package main

import (
	"rest-geoip/routes"
)

func main() {
	server := routes.SetupRouter()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
