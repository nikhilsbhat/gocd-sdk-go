package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/maintenance.json
var maintenanceJson string

func Test_client_EnableMaintenanceMode(t *testing.T) {
	t.Run("should enable the maintenance mode successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusNoContent, map[string]string{"Accept": gocd.HeaderVersionOne, "X-GoCD-Confirm": "true"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.EnableMaintenanceMode()
		assert.NoError(t, err)
	})

	t.Run("should error out while enabling maintenance mode as no valid headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusNoContent, map[string]string{"Accept": gocd.HeaderVersionTwo, "X-GoCD-Confirm": "true"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.EnableMaintenanceMode()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while making client call to enable maintenance mode", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.EnableMaintenanceMode()
		assert.EqualError(t, err, "call made to enable maintenance mode errored with Post "+
			"\"http://localhost:8156/go/api/admin/maintenance_mode/enable\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_DisableMaintenanceMode(t *testing.T) {
	t.Run("should disable the maintenance mode successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusNoContent, map[string]string{"Accept": gocd.HeaderVersionOne, "X-GoCD-Confirm": "true"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DisableMaintenanceMode()
		assert.NoError(t, err)
	})

	t.Run("should error out while disabling maintenance mode as no valid headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusNoContent, map[string]string{"Accept": gocd.HeaderVersionTwo, "X-GoCD-Confirm": "true"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DisableMaintenanceMode()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while making client call to disable maintenance mode", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DisableMaintenanceMode()
		assert.EqualError(t, err, "call made to enable maintenance mode errored with Post "+
			"\"http://localhost:8156/go/api/admin/maintenance_mode/disable\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_GetMaintenanceModeInfo(t *testing.T) {
	t.Run("should fetch the maintenance mode information successfully", func(t *testing.T) {
		server := mockServer([]byte(maintenanceJson), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := &gocd.Maintenance{}
		expected.MaintenanceInfo.Enabled = false
		expected.MaintenanceInfo.Metadata.UpdatedBy = "admin"
		expected.MaintenanceInfo.Metadata.UpdatedOn = "2019-01-02T04:18:28Z"

		actual, err := client.GetMaintenanceModeInfo()
		assert.NoError(t, err)
		assert.Equal(t, actual, actual)
	})

	t.Run("should error out with 404 while fetching maintenance mode information due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte("maintenanceJson"), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetMaintenanceModeInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'm' looking for beginning of value")
		assert.Equal(t, gocd.Maintenance{}, actual)
	})

	t.Run("should error out with 404 while fetching maintenance mode information due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(maintenanceJson), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetMaintenanceModeInfo()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Maintenance{}, actual)
	})

	t.Run("should error out while fetching maintenance mode information as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetMaintenanceModeInfo()
		assert.EqualError(t, err, "call made to enable maintenance mode errored with "+
			"Get \"http://localhost:8156/go/api/admin/maintenance_mode/info\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Maintenance{}, actual)
	})
}
