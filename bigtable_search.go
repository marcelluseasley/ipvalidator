package main

import (
	"context"
	"log"
	"net"
	"strings"

	"cloud.google.com/go/bigtable"
	"google.golang.org/api/option"
)

// project - contains the project in Google Cloud
// instance - the BigTable instance
// tableIPCountry - table keyed by IP network (from geolite2 database) to country geocode ID
// tableGeonameCountry - table keyed by Geoname ID (from geolite2 database) to country name
const (
	project             = "industrious-eye-236701"
	instance            = "ip-country-bt"
	tableIPCountry      = "IP-Country"
	tableGeonameCountry = "Geoname-Country"
)

// getFirstOctet - gets the first octet from the IP address
// This will narrow down having to search the entire IP-Country
// to just the IP addresses that start with the first octet.
func getFirstOctet(ip string) string {
	return strings.Split(ip, ".")[0]
}

// getGeoIDFromIP - takes an IP address in form "aaa.bbb.ccc.ddd"
// takes the first octect "aaa" and uses it to filter down the number
// of rows needed to search in the BigTable table
// For each of the rows that begin with the octet of our IP address,
// we parse the CIDR with net.ParseCIDR, and check to see if the IP address
// is contained within that CIDR range. If so, then we return the geocode ID
func getGeoIDFromIP(ip string) string {
	var geoID string
	ctx := context.Background()
	rowRange := bigtable.PrefixRange(getFirstOctet(ip))

	client, err := bigtable.NewClient(ctx, project, instance, option.WithCredentialsFile("ipgeo-readrights.json"))
	if err != nil {
		log.Fatalf("could not create data operations client: %v", err)
	}

	tblIPCountry := client.Open(tableIPCountry)

	var netIP net.IP
	netIP.UnmarshalText([]byte(ip))
	tblIPCountry.ReadRows(ctx, rowRange, func(row bigtable.Row) bool {
		for _, j := range row["ipGeoID"] {
			_, ipv4Net, _ := net.ParseCIDR(j.Row)
			if ipv4Net.Contains(netIP) {
				geoID = string(j.Value)
				return false
			}
		}

		return true
	}, bigtable.RowFilter(bigtable.ColumnFilter("geoname_id")))
	return geoID
}

// lookupCountryFromGeoID - takes the geocode ID and uses it as the
// key in the BigTable table tableGeonameCountry, which returns the
// associated country as a string
func lookupCountryFromGeoID(geoID string) string {
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, project, instance, option.WithCredentialsFile("ipgeo-readrights.json"))
	if err != nil {
		log.Fatalf("could not create data operations client: %v", err)
	}

	tblGeonameCountry := client.Open(tableGeonameCountry)
	row, err := tblGeonameCountry.ReadRow(ctx, geoID)
	if err != nil {
		log.Fatalf("could not read row from table: %v", err)
	}
	country := row["geoIDCountry"]
	return string(country[0].Value)
}
