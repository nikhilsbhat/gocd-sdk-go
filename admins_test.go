package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/system_admins.json
var systemAdmins string

func Test_client_GetAdminsInfo(t *testing.T) {
	t.Run("should error out while fetching system admins present from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetAdminsInfo()
		assert.EqualError(t, err, "call made to get system admin errored with: Get \"http://localhost:8153/go/api/admin/security/system_admins\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while fetching system admins present as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetAdminsInfo()
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while fetching system admins present as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetAdminsInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should be able to fetch admins present in GoCD server", func(t *testing.T) {
		server := mockServer([]byte(systemAdmins), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		expected := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
		}

		actual, err := client.GetAdminsInfo()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
