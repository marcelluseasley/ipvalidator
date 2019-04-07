package main

import (
	"net/http"
	"fmt"	
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "alive")
}
