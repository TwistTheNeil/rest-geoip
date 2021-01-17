package maxmind

import (
	"net"
	"os"
	"path/filepath"
	"rest-geoip/customerrors"
	"rest-geoip/utils"
	"sync"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/viper"
)

// Record captures the data resulting from a query to the maxmind database
type Record struct {
	Country struct {
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union"`
		ISOCode           string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
	IP string
}

// Info returns results from a maxmind db lookup
func Info(ip net.IP) (Record, error) {
	var record Record

	db, err := maxminddb.Open(viper.GetString("MAXMIND_DB_LOCATION") + viper.GetString("MAXMIND_DB"))
	if err != nil {
		return record, customerrors.ErrMMDBNotFound
	}
	defer db.Close()

	err = db.Lookup(ip, &record)
	if err != nil {
		return record, err
	}

	record.IP = ip.String()
	return record, nil
}

// DownloadAndUpdate the maxmind database
func DownloadAndUpdate() error {
	dbURL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + viper.GetString("MAXMIND_LICENSE") + "&suffix=tar.gz"
	md5URL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + viper.GetString("MAXMIND_LICENSE") + "&suffix=tar.gz.md5"
	dbDest := viper.GetString("MAXMIND_DB_LOCATION") + "/Geolite.tar.gz"
	md5Dest := viper.GetString("MAXMIND_DB_LOCATION") + "/Geolite.tar.gz.md5"

	// Make channels to pass errors in WaitGroup
	downloadErrors := make(chan error)
	wgDone := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)

	go utils.Download(dbURL, dbDest, &wg, downloadErrors)
	go utils.Download(md5URL, md5Dest, &wg, downloadErrors)

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

	if err := utils.VerifyMD5HashFromFile(dbDest, md5Dest); err != nil {
		return err
	}

	// Prepare a reader for extracting the tar.gz
	r, err := os.Open(dbDest) // #nosec G304
	if err != nil {
		return err
	}

	if err := utils.ExtractTarGz(r, viper.GetString("MAXMIND_DB_LOCATION")); err != nil {
		return err
	}

	// Move mmdb to MAXMIND_DB_LOCATION
	geoCityDBPath, _, err := utils.FindFile(viper.GetString("MAXMIND_DB_LOCATION"), "mmdb$")
	if err != nil {
		return err
	}

	if err = utils.MoveFile(geoCityDBPath, viper.GetString("MAXMIND_DB_LOCATION")+"/"+viper.GetString("MAXMIND_DB")); err != nil {
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

	return nil
}
