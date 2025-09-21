package maxmind

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"rest-geoip/internal/config"
	"rest-geoip/internal/errortypes"
	"rest-geoip/internal/fs"
	"sync"

	"github.com/oschwald/maxminddb-golang"
)

// DB struct
type DB struct {
	db *maxminddb.Reader
}

var instance *DB
var once sync.Once

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
	Subdivisions []struct {
		IsoCode   string `maxminddb:"iso_code"`
		GeoNameID uint   `maxminddb:"geoname_id"`
	} `maxminddb:"subdivisions"`
	IP string
}

// Open a maxmind database
func (m *DB) Open() error {
	dbLocation := config.Details().Maxmind.DBLocation + config.Details().Maxmind.DBFileName
	fmt.Printf("Opening db %s\n", dbLocation)

	_, err := os.Stat(dbLocation)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			e := errortypes.NewErrorDatabaseNotFound(err, dbLocation)
			return fmt.Errorf("maxmind.Open: db not found: %w", e)
		}
	}
	m.db, err = maxminddb.Open(dbLocation)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Close a maxmind database
func (m *DB) Close() error {
	return m.db.Close()
}

func (m *DB) Update() error {
	if config.Details().Maxmind.LicenseKey == "" {
		return fmt.Errorf("Error: Can't update database when no license key is set (MAXMIND_LICENSE env var needs to be set)")
	}
	if m == nil {
		err := m.Close()
		if err != nil {
			fmt.Println("Failed to close maxmind database")
			return err
		}
	}
	if err := DownloadAndUpdate(); err != nil {
		fmt.Println("Failed to update maxmind database")
		return err
	}
	if err := m.Open(); err != nil {
		fmt.Println("Failed to open maxmind database")
		return err
	}

	return nil
}

// Lookup results from a maxmind db lookup
func (m *DB) Lookup(ip net.IP) (Record, error) {
	var record Record

	err := m.db.Lookup(ip, &record)
	if err != nil {
		return record, err
	}

	record.IP = ip.String()
	return record, nil
}

// GetInstance of a maxmindReader
func GetInstance() *DB {
	once.Do(func() {
		instance = &DB{}
	})
	return instance
}

// DownloadAndUpdate the maxmind database
func DownloadAndUpdate() error {
	// TODO: check that db is closed
	dbURL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + config.Details().Maxmind.LicenseKey + "&suffix=tar.gz"
	md5URL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + config.Details().Maxmind.LicenseKey + "&suffix=tar.gz.md5"
	dbDest := config.Details().Maxmind.DBLocation + "/Geolite.tar.gz"
	md5Dest := config.Details().Maxmind.DBLocation + "/Geolite.tar.gz.md5"

	// Make channels to pass errors in WaitGroup
	downloadErrors := make(chan error)
	wgDone := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)

	go fs.Download(dbURL, dbDest, &wg, downloadErrors)
	go fs.Download(md5URL, md5Dest, &wg, downloadErrors)

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

	if err := fs.VerifyMD5HashFromFile(dbDest, md5Dest); err != nil {
		return err
	}

	// Prepare a reader for extracting the tar.gz
	r, err := os.Open(dbDest) // #nosec G304
	if err != nil {
		return err
	}

	if err := fs.ExtractTarGz(r, config.Details().Maxmind.DBLocation); err != nil {
		return err
	}

	// Move mmdb to MAXMIND_DB_LOCATION
	geoCityDBPath, _, err := fs.FindFile(config.Details().Maxmind.DBLocation, "mmdb$")
	if err != nil {
		return err
	}

	if err = fs.MoveFile(geoCityDBPath, config.Details().Maxmind.DBLocation+"/"+config.Details().Maxmind.DBFileName); err != nil {
		return err
	}

	// Remove all temporary downloaded files
	matches, err := filepath.Glob(config.Details().Maxmind.DBLocation + "GeoLite2-City_*")
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
