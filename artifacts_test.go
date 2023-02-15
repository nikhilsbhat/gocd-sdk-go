package gocd_test

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imdario/mergo"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/artifacts_info.json
var artifactInfoJSON string

func Test_client_GetArtifactConfig(t *testing.T) {
	correctArtifactsHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to get the artifact information from GoCD successfully", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte(artifactInfoJSON), http.StatusOK, correctArtifactsHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ArtifactInfo{
			ArtifactsDir: "foo",
			PurgeSettings: struct {
				PurgeStartDiskSpace float64 `json:"purge_start_disk_space,omitempty" yaml:"purge_start_disk_space,omitempty"`
				PurgeUptoDiskSpace  float64 `json:"purge_upto_disk_space,omitempty" yaml:"purge_upto_disk_space,omitempty"`
			}{PurgeStartDiskSpace: 10, PurgeUptoDiskSpace: 20},
			ETAG: "17f5a9edf150884e5fc4315b4a7814cd",
		}

		actual, err := client.GetArtifactConfig()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, expected.ETAG, actual.ETAG)
	})

	t.Run("should error out while getting the artifact information from GoCD successfully due to wrong headers", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte(artifactInfoJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetArtifactConfig()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})

	t.Run("should error out while getting the artifact information from GoCD successfully due to no headers set", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte(artifactInfoJSON), http.StatusOK, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetArtifactConfig()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})

	t.Run("should error out while getting the artifact information as server returned malformed data", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte("artifactInfoJSON"), http.StatusOK, correctArtifactsHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetArtifactConfig()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})

	t.Run("should error out while making client call to fetch the artifact information as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetArtifactConfig()
		assert.EqualError(t, err, "call made to get artifacts info errored with "+
			"Get \"http://localhost:8156/go/api/admin/config/server/artifact_config\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})
}

func Test_client_UpdateArtifactConfig(t *testing.T) {
	etag := "17f5a9edf150884e5fc4315b4a7814cd"
	correctArtifactsHeader := map[string]string{
		"Accept":       gocd.HeaderVersionOne,
		"Content-Type": gocd.ContentJSON,
		"If-Match":     etag,
	}
	t.Run("should be able to update the artifact config successfully", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte(artifactInfoJSON), http.StatusOK, correctArtifactsHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		artifactConfig := gocd.ArtifactInfo{
			ArtifactsDir: "foo",
			PurgeSettings: struct {
				PurgeStartDiskSpace float64 `json:"purge_start_disk_space,omitempty" yaml:"purge_start_disk_space,omitempty"`
				PurgeUptoDiskSpace  float64 `json:"purge_upto_disk_space,omitempty" yaml:"purge_upto_disk_space,omitempty"`
			}{
				PurgeStartDiskSpace: 20.0,
			},
			ETAG: "17f5a9edf150884e5fc4315b4a7814cd",
		}

		actual, err := client.UpdateArtifactConfig(artifactConfig)
		assert.NoError(t, err)
		assert.Equal(t, float64(20), actual.PurgeSettings.PurgeStartDiskSpace)
	})

	t.Run("should error out while updating the artifact information in GoCD due to wrong headers", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte(artifactInfoJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.UpdateArtifactConfig(gocd.ArtifactInfo{})
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})

	t.Run("should error out while updating the artifact information in GoCD successfully due to no headers set", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte(artifactInfoJSON), http.StatusOK, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.UpdateArtifactConfig(gocd.ArtifactInfo{})
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})

	t.Run("should error out while updating the artifact information as server returned malformed data", func(t *testing.T) {
		server := mockArtifactInfoServer([]byte("artifactInfoJSON"), http.StatusInternalServerError, correctArtifactsHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		artifactConfig := gocd.ArtifactInfo{
			ArtifactsDir: "foo",
			PurgeSettings: struct {
				PurgeStartDiskSpace float64 `json:"purge_start_disk_space,omitempty" yaml:"purge_start_disk_space,omitempty"`
				PurgeUptoDiskSpace  float64 `json:"purge_upto_disk_space,omitempty" yaml:"purge_upto_disk_space,omitempty"`
			}{
				PurgeStartDiskSpace: 20.0,
			},
			ETAG: "17f5a9edf150884e5fc4315b4a7814cd",
		}

		actual, err := client.UpdateArtifactConfig(artifactConfig)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})

	t.Run("should error out while updating the artifact information in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.UpdateArtifactConfig(gocd.ArtifactInfo{})
		assert.EqualError(t, err, "call made to update artifacts info errored with Post "+
			"\"http://localhost:8156/go/api/admin/config/server/artifact_config\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.ArtifactInfo{}, actual)
	})
}

func mockArtifactInfoServer(body []byte, statusCode int, header map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
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

		if req.Method == http.MethodPost { //nolint:nestif
			if statusCode == http.StatusInternalServerError {
				writer.WriteHeader(http.StatusOK)
				_, err := writer.Write(body)
				if err != nil {
					log.Fatalln(err)
				}

				return
			}

			var artifactCFG gocd.ArtifactInfo
			if err := json.Unmarshal(body, &artifactCFG); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				_, err = writer.Write([]byte(err.Error()))
				if err != nil {
					log.Fatalln(err)
				}
			}

			var newCFG gocd.ArtifactInfo
			if err := json.NewDecoder(req.Body).Decode(&newCFG); err != nil {
				log.Fatalln(err)
			}

			if err := mergo.Merge(&artifactCFG, &newCFG, mergo.WithOverride); err != nil {
				log.Fatalf("merging config errored with: %s", err.Error())
			}

			out, err := json.Marshal(artifactCFG)
			if err != nil {
				log.Fatalln(err)
			}

			writer.WriteHeader(statusCode)
			_, err = writer.Write(out)
			if err != nil {
				log.Fatalln(err)
			}

			return
		}

		writer.Header().Set("Etag", "17f5a9edf150884e5fc4315b4a7814cd")
		writer.WriteHeader(statusCode)
		_, err := writer.Write(body)
		if err != nil {
			log.Fatalln(err)
		}
	}))
}
