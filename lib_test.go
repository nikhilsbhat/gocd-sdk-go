package gocd_test

import (
	"net/http"
	"net/http/httptest"
)

func mockServer(body []byte, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, t *http.Request) {
		w.WriteHeader(statusCode)
		w.Write(body)
	}))
}
