package db

import (
	"github.com/oschwald/maxminddb-golang"
	"log"
	"net"
)

type ICountryProvider interface {
	GetCountry(ipAddress string) string
	GetLatLng(ipAddress string) (float64, float64)
	GetTimeZone(ipAddress string) string
}

type dbProvider struct {
	db *maxminddb.Reader
}

var CountryProvider ICountryProvider

func init() {
	db, err := maxminddb.Open("./db/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	CountryProvider = &dbProvider{db: db}

}

func (provider *dbProvider) GetCountry(ipAddress string) string {

	ip := net.ParseIP(ipAddress)

	var record3 struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err := provider.db.Lookup(ip, &record3)
	if err != nil {
		log.Fatal(err)
	}

	return record3.Country.ISOCode
}

func (provider *dbProvider) GetLatLng(ipAddress string) (float64, float64) {

	ip := net.ParseIP(ipAddress)

	var record3 struct {
		Location struct {
			Latitude  float64 `maxminddb:"latitude"`
			Longitude float64 `maxminddb:"longitude"`
		} `maxminddb:"location"`
	}

	err := provider.db.Lookup(ip, &record3)
	if err != nil {
		log.Fatal(err)
	}

	return record3.Location.Latitude, record3.Location.Longitude

}

func (provider *dbProvider) GetTimeZone(ipAddress string) string {

	ip := net.ParseIP(ipAddress)

	var record3 struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`

		Location struct {
			TimeZone string `maxminddb:"time_zone"`
		} `maxminddb:"location"`
	}

	err := provider.db.Lookup(ip, &record3)
	if err != nil {
		log.Fatal(err)
	}

	return record3.Location.TimeZone

}