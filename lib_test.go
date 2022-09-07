package gocd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
)

func mockServer(body []byte, statusCode int, header map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range header {
			if r.Header.Get(key) != value {
				w.WriteHeader(http.StatusNotFound)
				if _, err := w.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}

		w.WriteHeader(statusCode)
		_, err := w.Write(body)
		if err != nil {
			log.Fatalln(err)
		}
	}))
}
