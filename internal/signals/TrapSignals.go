package signals

import (
	"fmt"
	"os"
	"os/signal"
	"rest-geoip/internal/maxmind"
	"syscall"
)

// Trap traps signals for these purposes:
// SIGUSR1: Update database
func Trap() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for {
			switch <-signals {
			case syscall.SIGUSR1:
				fmt.Println("SIGUSR1 called. Updating maxmind db")
				if err := maxmind.GetInstance().Close(); err != nil {
					fmt.Println("Failed to close maxmind database")
				}
				if err := maxmind.DownloadAndUpdate(); err != nil {
					fmt.Println("Failed to update maxmind database")
					continue
				}
				if err := maxmind.GetInstance().Open(); err != nil {
					fmt.Println("Failed to open maxmind database")
					continue
				}
			}
		}
	}()
}
