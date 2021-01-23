package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"rest-geoip/lib/maxmind"
	"sync"

	"github.com/spf13/viper"
)

// DownloadAndUpdate the maxmind database
func DownloadAndUpdate() error {
	if err := maxmind.GetInstance().Close(); err != nil {
		fmt.Println(err)
		fmt.Println("Failed to close maxmind db")
		return err
	}

	dbURL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + viper.GetString("MAXMIND_LICENSE") + "&suffix=tar.gz"
	md5URL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + viper.GetString("MAXMIND_LICENSE") + "&suffix=tar.gz.md5"
	dbDest := viper.GetString("MAXMIND_DB_LOCATION") + "/Geolite.tar.gz"
	md5Dest := viper.GetString("MAXMIND_DB_LOCATION") + "/Geolite.tar.gz.md5"

	// Make channels to pass errors in WaitGroup
	downloadErrors := make(chan error)
	wgDone := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)

	go Download(dbURL, dbDest, &wg, downloadErrors)
	go Download(md5URL, md5Dest, &wg, downloadErrors)

	// wait until WaitGroup is done
	// Sends a signal we need to catch in the select
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received
	select {
	case <-wgDone:
		break
	case err := <-downloadErrors:
		// close(downloadErrors)
		return err
	}

	if err := VerifyMD5HashFromFile(dbDest, md5Dest); err != nil {
		return err
	}

	// Prepare a reader for extracting the tar.gz
	r, err := os.Open(dbDest) // #nosec G304
	if err != nil {
		return err
	}

	if err := ExtractTarGz(r, viper.GetString("MAXMIND_DB_LOCATION")); err != nil {
		return err
	}

	// Move mmdb to MAXMIND_DB_LOCATION
	geoCityDBPath, _, err := FindFile(viper.GetString("MAXMIND_DB_LOCATION"), "mmdb$")
	if err != nil {
		return err
	}

	if err = MoveFile(geoCityDBPath, viper.GetString("MAXMIND_DB_LOCATION")+"/"+viper.GetString("MAXMIND_DB")); err != nil {
		return err
	}

	// Remove all temporary downloaded files
	matches, err := filepath.Glob(viper.GetString("MAXMIND_DB_LOCATION") + "GeoLite2-City_*")
	if err != nil {
		return err
	}
	matches = append(matches, dbDest)
	matches = append(matches, md5Dest)
	for _, v := range matches {
		if err := os.RemoveAll(v); err != nil {
			return err
		}
	}

	if err := maxmind.GetInstance().Open(); err != nil {
		fmt.Println(err)
		fmt.Println("Failed to open maxmind db")
		return err
	}

	return nil
}
