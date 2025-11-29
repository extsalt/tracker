package geo

import (
	"log"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var (
	db   *geoip2.Reader
	once sync.Once
)

func Init(dbPath string) error {
	var err error
	once.Do(func() {
		db, err = geoip2.Open(dbPath)
		if err != nil {
			log.Printf("Failed to open GeoIP database: %v", err)
		}
	})
	return err
}

func Lookup(ipStr string) (country, state, city string, err error) {
	if db == nil {
		return "", "", "", nil
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", "", nil
	}

	record, err := db.City(ip)
	if err != nil {
		return "", "", "", err
	}

	country = record.Country.IsoCode
	if len(record.Subdivisions) > 0 {
		state = record.Subdivisions[0].IsoCode
	}
	city = record.City.Names["en"]

	return country, state, city, nil
}
