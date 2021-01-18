package signals

import (
	"fmt"
	"os"
	"os/signal"
	"rest-geoip/maxmind"
	"syscall"
)

// Trap traps signals for these purposes:
// SIGUSR1: Update database
// SIGUSR2: Close database
func Trap() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for {
			switch <-signals {
			case syscall.SIGUSR1:
				fmt.Println("SIGUSR1 called. Updating maxmind db")
				maxmind.GetInstance().Close()
				maxmind.DownloadAndUpdate()
				maxmind.GetInstance().Open()
			case syscall.SIGUSR2:
				fmt.Println("SIGUSR2 called. Closing maxmind db")
				maxmind.GetInstance().Close()
			}
		}
	}()
}
