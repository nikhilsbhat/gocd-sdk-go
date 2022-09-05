package gocd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
)

func mockServer(body []byte, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, t *http.Request) {
		w.WriteHeader(statusCode)
		_, err := w.Write(body)
		if err != nil {
			log.Fatalln(err)
		}
	}))
}
