package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/permissions.json
var permissionsJSON string

func Test_client_GetPermissions(t *testing.T) {
	correctPermissionHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should be able to fetch permissions information successfully", func(t *testing.T) {
		server := mockServer([]byte(permissionsJSON), http.StatusOK,
			correctPermissionHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Permission{
			Environment: gocd.EntityPermissions{
				View:       []string{"DEV", "UAT", "QA", "PROD"},
				Administer: []string{"DEV"},
			},
			ConfigRepo: gocd.EntityPermissions{
				View:       []string{"dev-pipelines-repo", "prod-pipelines-repo"},
				Administer: []string{"dev-pipelines-repo"},
			},
			ClusterProfile: gocd.EntityPermissions{
				View:       []string{"dev-cluster", "prod-cluster"},
				Administer: []string{"dev-cluster"},
			},
			ElasticAgentProfile: gocd.EntityPermissions{
				View:       []string{"build-agent", "deploy-agent"},
				Administer: []string{"build-agent"},
			},
		}

		actual, err := client.GetPermissions(nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching permissions info from GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(permissionsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Permission{}

		actual, err := client.GetPermissions(nil)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/auth/permissions\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching permissions info from GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(permissionsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Permission{}

		actual, err := client.GetPermissions(nil)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/auth/permissions\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching permissions info from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("permissionsJSON"), http.StatusOK, correctPermissionHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Permission{}

		actual, err := client.GetPermissions(nil)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching permissions info from GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.Permission{}

		actual, err := client.GetPermissions(nil)
		assert.EqualError(t, err, "call made to get permissions errored with: "+
			"Get \"http://localhost:8156/go/api/auth/permissions\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
