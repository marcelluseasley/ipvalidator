package main

import (
	"net/http"
	"log"
	
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/ipvalidate", ipValidateHandler).Methods("POST")
	router.HandleFunc("/v1/healthcheck", healthcheckHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}

