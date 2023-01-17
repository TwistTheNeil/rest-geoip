package maxmind

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/viper"
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
	IP string
}

// Open a maxmind database
func (m *DB) Open() error {
	var err error
	dbLocation := viper.GetString("MAXMIND_DB_LOCATION") + viper.GetString("MAXMIND_DB")
	fmt.Printf("Opening db %s\n", dbLocation)

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
