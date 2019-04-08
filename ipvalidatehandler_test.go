package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	test1 = IPCountries{"172.58.7.17", []string{"United States", "Germany", "Russia"}}
	test2 = IPCountries{"172.58.7.17", []string{"Japan", "Germany", "Russia"}}
)

func IPRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/ipvalidate", ipValidateHandler).Methods("POST")
	return router
}

func TestIPValidateHandler(t *testing.T) {
	jsonTest1, _ := json.Marshal(test1)
	request, _ := http.NewRequest("POST", "/v1/ipvalidate", bytes.NewBuffer(jsonTest1))
	response := httptest.NewRecorder()
	IPRouter().ServeHTTP(response, request)
	assert.Equal(t, `{"valid_status": "true"}`, string(response.Body.Bytes()), "")
}

func TestIPValidateHandlerTestFalseStatus(t *testing.T) {
	jsonTest2, _ := json.Marshal(test2)
	request, _ := http.NewRequest("POST", "/v1/ipvalidate", bytes.NewBuffer(jsonTest2))
	response := httptest.NewRecorder()
	IPRouter().ServeHTTP(response, request)
	assert.Equal(t, `{"valid_status": "false"}`, string(response.Body.Bytes()), "")
}
