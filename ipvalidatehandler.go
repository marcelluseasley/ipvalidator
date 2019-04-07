package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// IPCountries - holds IP address and list of approved countries.
type IPCountries struct {
	IP        string   `json:"ip"`
	Countries []string `json:"approved_countries"`
}

func ipValidateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received connection from: %v", r.RemoteAddr)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	decoder := json.NewDecoder(r.Body)
	var ipcountries IPCountries
	err := decoder.Decode(&ipcountries)
	if err != nil {
		log.Printf("Could not process incoming JSON: %v", err)
		fmt.Fprint(w, "invalid JSON request")
	}
	

	geocode := getGeoIDFromIP(ipcountries.IP)
	country := lookupCountryFromGeoID(geocode)

	if contains(ipcountries.Countries, country){
		fmt.Fprintf(w, `{"valid_status": "true"	}`)
	} else {
		fmt.Fprintf(w, `{"valid_status": "false"}`)
	}


}

func contains(forest []string, tree string) bool {
	for _, t := range forest {
		if tree == t {
			return true
		}
	}
	return false
}