package gocd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
)

func mockServer(body []byte, statusCode int, header map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
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

		writer.WriteHeader(statusCode)
		_, err := writer.Write(body)
		if err != nil {
			log.Fatalln(err)
		}
	}))
}
