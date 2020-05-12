package geoip

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type MaxmindRecord struct {
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
}

func GetIpAddressInfo(ipAddress string) (MaxmindRecord, error) {
	var record MaxmindRecord

	db, err := maxminddb.Open("db.mmdb")
	if err != nil {
		return record, err
	}
	defer db.Close()

	ip := net.ParseIP(ipAddress)

	err = db.Lookup(ip, &record)
	if err != nil {
		return record, err
	}

	return record, nil
}
