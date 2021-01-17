package main

import (
	"rest-geoip/routes"

	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("MAXMIND_DB_LOCATION", "/opt/")
	viper.SetDefault("MAXMIND_DB", "GeoLite2-City.mmdb")
	viper.SetDefault("LOGGING", "false")
	viper.SetDefault("WEB", true)
	viper.SetDefault("LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("LISTEN_PORT", "8080")
	viper.AutomaticEnv()

	server := routes.SetupRouter()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
