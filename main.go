package main

import (
	"errors"
	"fmt"
	"os"
	"rest-geoip/internal/config"
	"rest-geoip/internal/errortypes"
	"rest-geoip/internal/maxmind"
	"rest-geoip/internal/router"
	"rest-geoip/internal/signals"
)

func main() {
	config.Init()
	signals.Trap()
	err := maxmind.GetInstance().Open()

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
