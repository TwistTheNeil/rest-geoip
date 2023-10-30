package main

import (
	"errors"
	"fmt"
	"os"
	"rest-geoip/internal/errortypes"
	"rest-geoip/internal/maxmind"
	"rest-geoip/internal/random"
	"rest-geoip/internal/router"
	"rest-geoip/internal/signals"

	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("MAXMIND_DB_LOCATION", "/opt/")
	viper.SetDefault("MAXMIND_DB", "GeoLite2-City.mmdb")
	viper.SetDefault("LOGGING", false)
	viper.SetDefault("WEB", true)
	viper.SetDefault("LISTEN_ADDRESS", "0.0.0.0")
	viper.SetDefault("LISTEN_PORT", "1323")
	viper.SetDefault("RELEASE_MODE", "true")
	viper.SetDefault("MAPTILER_TOKEN", "token")
	viper.AutomaticEnv()
	if viper.GetString("API_KEY") == "" {
		generatedKey, err := random.GenerateKey(512)
		if err != nil {
			fmt.Println("Error: no api key set and we weren't able to generate a key. Exiting")
			os.Exit(1)
		}

		viper.SetDefault("API_KEY", generatedKey)
		fmt.Printf("Generated API key: %s\n", viper.GetString("API_KEY"))
	}

	signals.Trap()
	err := maxmind.
		GetInstance().
		Open()

	if err != nil {
		var ErrDatabaseNotFound *errortypes.ErrorDatabaseNotFound
		if errors.As(errors.Unwrap(err), &ErrDatabaseNotFound) {
			fmt.Println(err)
			fmt.Println("Attempting to fetch database")
			fetchErr := maxmind.GetInstance().Update()
			if fetchErr != nil {
				fmt.Println(fetchErr)
				os.Exit(1)
			}
		} else {

			fmt.Println("Error: Maxmind database not opened during initialization. Please check that it exists and is readable.")
			os.Exit(1)
		}
	}

	router.InitRouter()
}
