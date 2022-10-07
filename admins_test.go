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
	correctAdminHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	t.Run("should error out while fetching system admins present from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetSystemAdmins()
		assert.EqualError(t, err, "call made to get system admin errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/system_admins\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while fetching system admins present as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctAdminHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetSystemAdmins()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while fetching system admins present as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctAdminHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetSystemAdmins()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should get 404 from server as header messed up", func(t *testing.T) {
		server := mockServer([]byte(systemAdmins), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetSystemAdmins()
		assert.EqualError(t, err, "goCd server returned code 404 with message")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should be able to fetch admins present in GoCD server", func(t *testing.T) {
		server := mockServer([]byte(systemAdmins), http.StatusOK, correctAdminHeader, false, nil)
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

		actual, err := client.GetSystemAdmins()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateSystemAdmins(t *testing.T) {
	correctAdminHeader := map[string]string{
		"Accept":       gocd.HeaderVersionTwo,
		"Content-Type": gocd.ContentJSON,
		"If-Match":     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
	}
	t.Run("should be able to update the system admins successfully", func(t *testing.T) {
		server := mockServer([]byte(systemAdmins), http.StatusOK, correctAdminHeader,
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		users := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
			ETAG:  "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		expected := users
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateSystemAdmins(users)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating system admins due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(systemAdmins), http.StatusOK, map[string]string{
			"Accept":       gocd.HeaderVersionThree,
			"Content-Type": gocd.ContentJSON,
			"If-Match":     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		},
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		users := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
			ETAG:  "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.UpdateSystemAdmins(users)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while updating system admins due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(systemAdmins), http.StatusOK, nil,
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		users := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
			ETAG:  "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.UpdateSystemAdmins(users)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while updating system admins as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("systemAdmins"), http.StatusOK, correctAdminHeader,
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		users := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
			ETAG:  "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.UpdateSystemAdmins(users)
		assert.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while updating system admins as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("systemAdmins"), http.StatusOK, correctAdminHeader,
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		users := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
			ETAG:  "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.UpdateSystemAdmins(users)
		assert.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})

	t.Run("should error out while updating system admins as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		users := gocd.SystemAdmins{
			Roles: []string{"manager"},
			Users: []string{"john", "maria"},
			ETAG:  "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.UpdateSystemAdmins(users)
		assert.EqualError(t, err, "call made to update system admin errored with: "+
			"Put \"http://localhost:8156/go/api/admin/security/system_admins\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.SystemAdmins{}, actual)
	})
}
