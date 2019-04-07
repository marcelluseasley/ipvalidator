package main

import (
	"net/http"
	"fmt"	
	"log"
	"encoding/json"
)

// IPCountries - holds IP address and list of approved countries.
type IPCountries struct {
	IP string `json:"ip"`
	Countries []string `json:"approved_countries"`
}



func ipValidateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ipcountries IPCountries
	err := decoder.Decode(&ipcountries)
	if err != nil {
		log.Printf("Could not process incoming JSON: %v", err)
		fmt.Fprint(w, "invalid JSON request")
	}
	fmt.Printf(string(ipcountries.Countries[0]))
}

