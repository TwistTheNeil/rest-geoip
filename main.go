package main

import (
	"fmt"
	"os"
	"rest-geoip/internal/maxmind"
	"rest-geoip/internal/router"

	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("MAXMIND_DB_LOCATION", "/opt/")
	viper.SetDefault("MAXMIND_DB", "GeoLite2-City.mmdb")
	viper.SetDefault("LOGGING", false)
	viper.SetDefault("WEB", true)
	viper.SetDefault("LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("LISTEN_PORT", "1323")
	viper.AutomaticEnv()

	// signals.Trap()
	err := maxmind.
		GetInstance().
		Open()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: Maxmind database not opened during initialization. Please check that it exists.")
		os.Exit(1)
	}

	router.InitRouter()
}
