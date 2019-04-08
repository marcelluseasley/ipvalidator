package main

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func HealthRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/healthcheck", healthcheckHandler).Methods("GET")
	return router
}
func TestHealthcheckHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "/v1/healthcheck", nil)
	response := httptest.NewRecorder()
	HealthRouter().ServeHTTP(response, request)
	assert.Equal(t, "alive", string(response.Body.Bytes()), "alive is expected")
}
