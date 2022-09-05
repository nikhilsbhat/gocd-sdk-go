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
	t.Run("should error out while fetching health status information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetHealthMessages()
		assert.EqualError(t, err, "call made to get health info errored with Get \"http://localhost:8153/go/api/server_health_messages\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching health status information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetHealthMessages()
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching health status information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetHealthMessages()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch the server health status", func(t *testing.T) {
		server := mockServer([]byte(healthMessages), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		healthStatus, err := client.GetHealthMessages()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(healthStatus))
		assert.Equal(t, "WARNING", healthStatus[0].Level)
	})
}
