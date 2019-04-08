package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"archive/zip"
	"encoding/csv"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/bigtable"
)

const (
	tableIPCountry  = "IP-Country"
	columnFamily1   = "ipGeoID"
	columnCIDR      = "network"
	columnGeonameID = "geoname_id"

	tableGeonameCountry = "Geoname-Country"
	columnFamily2       = "geoIDCountry"
	columnCountry       = "country"

	geolite2URL                        = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip"
	geoliteCountryIPBlocksFile         = "GeoLite2-Country-Blocks-IPv4.csv"
	geoliteCountryLocationsEnglishFile = "GeoLite2-Country-Locations-en.csv"
)

var dir string

func downloadGeoLiteZIPFile(url string) (*bytes.Reader, error) {

	buf := &bytes.Buffer{}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil

}

func extractZipFiles() {
	urlReader, err := downloadGeoLiteZIPFile(geolite2URL)
	if err != nil {
		log.Fatalf("Unable to download GeoLite zip file: %s", err)
	}

	zReader, err := zip.NewReader(urlReader, int64(urlReader.Len()))
	if err != nil {
		log.Fatalf("Unable to process GeoLite zip file: %s", err)
	}
	for _, zFile := range zReader.File {
		writeFile(zFile)
	}

}

func writeFile(zFile *zip.File) {
	fileDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current directory for writing: %s", err)
	}
	f, err := zFile.Open()
	if err != nil {
		log.Printf("Unable to read file: %s\n", zFile.Name)
		return
	}
	defer f.Close()

	path := strings.Replace(filepath.Join(fileDirectory, zFile.Name), `/`, string(filepath.Separator), -1)
	dir, _ = filepath.Split(path)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		log.Printf("unable to create directory: %s", dir)
		return
	}

	ff, err := os.Create(path)
	if err != nil {
		log.Printf("Unable to create file: %s", path)
		return
	}
	defer ff.Close()

	_, err = io.Copy(ff, f)
	if err != nil {
		log.Printf("Cannot write file %s: %s", path, err)
		return
	}
	return

}

func getCSVRecords(fileDir string, csvFileName string) ([][]string, error) {
	f, err := os.Open(fileDir + csvFileName)
	if err != nil {
		log.Printf("Unable to read from CSV file: %s", err)
		return nil, err
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Printf("Unable to read records from CSV file: %s", err)
		return nil, err
	}

	return records, nil
}

func createBigTableTablesFromCSVRecords(records [][]string, tableName string, columnFamily string, column1 string) {

	project := "industrious-eye-236701"
	instance := "ip-country-bt"

	ctx := context.Background()

	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		log.Fatalf("Could not create admin client: %v", err)
	}

	tables, err := adminClient.Tables(ctx)
	if err != nil {
		log.Fatalf("Could not fetch table list: %v", err)
	}

	if !sliceContains(tables, tableName) {
		log.Printf("Creating table %s", tableName)
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			log.Fatalf("Could not create table %s: %v", tableName, err)
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, tableName)
	if err != nil {
		log.Fatalf("Could not read info for table %s: %v", tableName, err)
	}

	if !sliceContains(tblInfo.Families, columnFamily) {
		if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamily); err != nil {
			log.Fatalf("Could not create column family %s: %v", columnFamily, err)
		}
	}

	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}

	switch tableName {
	case tableIPCountry:
		log.Println("inside case tableIPCountry")
		tbl := client.Open(tableName)

		muts := make([]*bigtable.Mutation, 100000)
		rowKeys := make([]string, 100000)
		x := 0
		for i, record := range records {
			if i == 0 {
				continue
			}
			muts[x] = bigtable.NewMutation()
			muts[x].Set(columnFamily, column1, bigtable.Now(), []byte(record[1]))

			rowKeys[x] = fmt.Sprintf("%s", record[0])
			if i%100000 == 0 {
				rowErrs, err := tbl.ApplyBulk(ctx, rowKeys, muts)
				if err != nil {
					log.Fatalf("Could not apply bulk row mutation: %v", err)
				}
				if rowErrs != nil {
					for _, rowErr := range rowErrs {
						log.Printf("Error writing row: %v", rowErr)
					}
					log.Fatalf("Could not write some rows")
				}
				x = -1
			}
			x++

		}
		rowErrs, err := tbl.ApplyBulk(ctx, rowKeys, muts)
		if err != nil {
			log.Fatalf("Could not apply bulk row mutation: %v", err)
		}
		if rowErrs != nil {
			for _, rowErr := range rowErrs {
				log.Printf("Error writing row: %v", rowErr)
			}
			log.Fatalf("Could not write some rows")
		}

	case tableGeonameCountry:
		log.Println("inside case tableGeonameCountry")
		tbl := client.Open(tableName)


		for i, record := range records {
			if i == 0 {
				continue
			}
			mut := bigtable.NewMutation()
			mut.Set(columnFamily, column1, bigtable.Now(), []byte(record[5]))

			rowKey := fmt.Sprintf("%s", record[0])

			err := tbl.Apply(ctx, rowKey, mut)
			if err != nil {
				log.Fatalf("Could not apply bulk row mutation: %v", err)
			}

		}

	}

	if err = adminClient.Close(); err != nil {
		log.Fatalf("Could not close admin client: %v", err)
	}

}

func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func main() {
	extractZipFiles()

	tableIPRecords, err := getCSVRecords(dir, geoliteCountryIPBlocksFile)
	if err != nil {
		log.Fatalf("Error retreiving records from CSV file: %s", err)
	}

	createBigTableTablesFromCSVRecords(tableIPRecords, tableIPCountry, columnFamily1, columnGeonameID)

	tableGeonameCountryRecords, err := getCSVRecords(dir, geoliteCountryLocationsEnglishFile)
	if err != nil {
		log.Fatalf("Error retreiving records from CSV file: %s", err)
	}

	createBigTableTablesFromCSVRecords(tableGeonameCountryRecords, tableGeonameCountry, columnFamily2, columnCountry)

}
