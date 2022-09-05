package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/config_repos.json
var configRepos string

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should error out while fetching config repo information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepoInfo()
		assert.EqualError(t, err, "call made to get config repo errored with Get \"http://localhost:8153/go/api/admin/config_repos\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repo information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepoInfo()
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repo information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepoInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able retrieve config repo information", func(t *testing.T) {
		server := mockServer([]byte(configRepos), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepoInfo()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(actual))
	})
}
