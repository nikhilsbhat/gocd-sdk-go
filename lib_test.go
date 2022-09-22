package gocd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
)

func mockServer(body []byte, statusCode int, header map[string]string, nilHeader bool, additionalHeaders map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if !nilHeader {
			if header == nil {
				writer.WriteHeader(http.StatusNotFound)
				if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
					log.Fatalln(err)
				}

				return
			}
		}

		for key, value := range header {
			if req.Header.Get(key) != value {
				writer.WriteHeader(http.StatusNotFound)
				if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
					log.Fatalln(err)
				}

				return
			}
		}

		for key, value := range additionalHeaders {
			writer.Header().Set(key, value)
		}

		writer.WriteHeader(statusCode)
		_, err := writer.Write(body)
		if err != nil {
			log.Fatalln(err)
		}
	}))
}
