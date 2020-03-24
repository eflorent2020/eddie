package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, bodyData string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(bodyData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStatusUnauthorized(t *testing.T) {
	// Grab our router
	router := SetupRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/api/v1/rest/users/", "")
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func TestLogin(t *testing.T) {
	db, err = bootstrap()
	defer db.Close()
	// Grab our router
	router := SetupRouter()
	getToken(t, router)
}

func getToken(t *testing.T, router http.Handler) string {
	w := performRequest(router, "POST", "/login", `{"username": "admin@eddie.onthewifi.com","password": "aze123"}`)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	token, exists := response["token"]
	assert.Nil(t, err)
	assert.True(t, exists)
	return token.(string)

}

// see https://github.com/gothinkster/golang-gin-realworld-example-app/blob/master/common/unit_test.go
