package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/plugin_settings.json
	pluginSettingJSON string
	pluginName        = "github.oauth.login"
)

func Test_client_GetPluginSettings(t *testing.T) {
	t.Run("should be able to fetch the plugin setting successfully", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.PluginSettings{
			ID: pluginName,
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "consumer_key",
					Value: "consumerkey",
				},
				{
					Key:   "consumer_secret",
					Value: "consumersecret",
				},
				{
					Key:   "server_base_url",
					Value: "https://ci.example.com",
				},
				{
					Key:            "password",
					EncryptedValue: "aSdiFgRRZ6A=",
				},
			},
			ETAG: "05548388f7ef5042cd39f7fe42e85735",
		}

		actual, err := client.GetPluginSettings(pluginName)
		assert.NoError(t, err)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.ETAG, actual.ETAG)
		assert.ElementsMatch(t, expected.Configuration, actual.Configuration)
	})

	t.Run("should error out while fetching the plugin setting due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.PluginSettings{}

		actual, err := client.GetPluginSettings(pluginName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching the plugin setting due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			nil, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.PluginSettings{}

		actual, err := client.GetPluginSettings(pluginName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching the plugin setting as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pluginSettingJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.PluginSettings{}

		actual, err := client.GetPluginSettings(pluginName)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching the plugin setting as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.PluginSettings{}

		actual, err := client.GetPluginSettings(pluginName)
		assert.EqualError(t, err, "call made to get 'github.oauth.login' plugin setting errored with: "+
			"Get \"http://localhost:8156/go/api/admin/plugin_settings/github.oauth.login\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreatePluginSettings(t *testing.T) {
	pluginHeaders := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to create the plugin setting successfully", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK, pluginHeaders, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{
			ID: pluginName,
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "consumer_key",
					Value: "consumerkey",
				},
				{
					Key:   "consumer_secret",
					Value: "consumersecret",
				},
				{
					Key:   "server_base_url",
					Value: "https://ci.example.com",
				},
				{
					Key:            "password",
					EncryptedValue: "aSdiFgRRZ6A=",
				},
			},
		}

		expected := pluginSettings

		actual, err := client.CreatePluginSettings(pluginSettings)
		assert.NoError(t, err)
		assert.Equal(t, expected.ID, actual.ID)
		assert.ElementsMatch(t, expected.Configuration, actual.Configuration)
	})

	t.Run("should error out while creating the plugin setting due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.CreatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating the plugin setting due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			nil, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.CreatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating the plugin setting as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pluginSettingJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.CreatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating the plugin setting as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.CreatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "call made to create plugin setting of 'github.oauth.login' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/plugin_settings\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdatePluginSettings(t *testing.T) {
	pluginHeaders := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON, "If-Match": "05548388f7ef5042cd39f7fe42e85735"}

	t.Run("should be able to update the plugin setting successfully", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK, pluginHeaders, false, map[string]string{"ETag": "e89135b38ddbcd9380c83eb524647bdd"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{
			ID: pluginName,
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "consumer_key",
					Value: "consumerkey",
				},
				{
					Key:   "consumer_secret",
					Value: "consumersecret",
				},
				{
					Key:   "server_base_url",
					Value: "https://ci.example.com",
				},
				{
					Key:            "password",
					EncryptedValue: "aSdiFgRRZ6A=",
				},
			},
		}

		expected := pluginSettings
		pluginSettings.ETAG = "05548388f7ef5042cd39f7fe42e85735"

		actual, err := client.UpdatePluginSettings(pluginSettings)
		assert.NoError(t, err)
		assert.Equal(t, expected.ID, actual.ID)
		assert.ElementsMatch(t, expected.Configuration, actual.Configuration)
	})

	t.Run("should error out while updating the plugin setting due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.UpdatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating the plugin setting due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginSettingJSON), http.StatusOK,
			nil, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.UpdatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating the plugin setting as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pluginSettingJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, map[string]string{"ETag": "05548388f7ef5042cd39f7fe42e85735"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.UpdatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating the plugin setting as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		pluginSettings := gocd.PluginSettings{ID: pluginName}
		expected := gocd.PluginSettings{}

		actual, err := client.UpdatePluginSettings(pluginSettings)
		assert.EqualError(t, err, "call made to update plugin setting of 'github.oauth.login' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/plugin_settings/github.oauth.login\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
