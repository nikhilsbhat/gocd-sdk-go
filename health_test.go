package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/server_healsth_messages.json
var healthMessages string

func TestConfig_GetHealthInfo(t *testing.T) {
	correctConfigHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should error out while fetching health status information from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetServerHealthMessages()
		assert.EqualError(t, err, "call made to get health info errored with "+
			"Get \"http://localhost:8156/go/api/server_health_messages\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching health status information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetServerHealthMessages()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching health status information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetServerHealthMessages()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch the server health status", func(t *testing.T) {
		server := mockServer([]byte(healthMessages), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		healthStatus, err := client.GetServerHealthMessages()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(healthStatus))
		assert.Equal(t, "WARNING", healthStatus[0].Level)
	})
}

func Test_client_GetServerHealth(t *testing.T) {
	serverHealthJSON := `{"health": "OK"}`
	t.Run("should be able to fetch server health successfully", func(t *testing.T) {
		server := mockServer([]byte(serverHealthJSON), http.StatusOK, nil, true, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := map[string]string{
			"health": "OK",
		}

		actual, err := client.GetServerHealth()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching server health from GoCD as server returned non 200 response", func(t *testing.T) {
		server := mockServer([]byte("serverHealthJSON"), http.StatusInternalServerError, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.GetServerHealth()
		assert.EqualError(t, err, "body: serverHealthJSON httpcode: 500")
		assert.Nil(t, response)
	})

	t.Run("should error out while fetching server health from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("serverHealthJSON"), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.GetServerHealth()
		assert.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Nil(t, response)
	})

	t.Run("should error out while fetching server health from GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		response, err := client.GetServerHealth()
		assert.EqualError(t, err, "call made to get server health errored with "+
			"Get \"http://localhost:8156/go/api/v1/health\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, response)
	})
}
