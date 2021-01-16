package main

import (
	"rest-geoip/routes"

	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	server := routes.SetupRouter()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
