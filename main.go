package main

import (
	"os"
	"rest-geoip/routes"
)

func main() {
	listenAddress := "localhost"
	listenPort := "8080"

	if os.Getenv("LISTEN_ADDRESS") != "" {
		listenAddress = os.Getenv("LISTEN_ADDRESS")
	}
	if os.Getenv("LISTEN_PORT") != "" {
		listenPort = os.Getenv("LISTEN_PORT")
	}

	router := routes.SetupRouter()
	if err := router.Run(listenAddress + ":" + listenPort); err != nil {
		panic(err)
	}
}
