package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed internal/fixtures/version.json
var versionInfo string

func Test_config_GetVersionInfo(t *testing.T) {
	correctVersionHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should error out while fetching version information from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetVersionInfo()
		require.EqualError(t, err, "call made to get version information errored with: "+
			"Get \"http://localhost:8156/go/api/version\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.VersionInfo{}, actual)
	})

	t.Run("should error out while fetching version information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctVersionHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetVersionInfo()
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+"/api/version\nwith BODY:backupJSON")
		assert.Equal(t, gocd.VersionInfo{}, actual)
	})

	t.Run("should error out while fetching version information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctVersionHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetVersionInfo()
		require.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.VersionInfo{}, actual)
	})

	t.Run("should be able to fetch the version info", func(t *testing.T) {
		server := mockServer([]byte(versionInfo), http.StatusOK, correctVersionHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.VersionInfo{
			Version:     "16.6.0",
			FullVersion: "16.6.0 (3348-a7a5717cbd60c30006314fb8dd529796c93adaf0)",
			GitSHA:      "a7a5717cbd60c30006314fb8dd529796c93adaf0",
		}

		actual, err := client.GetVersionInfo()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
