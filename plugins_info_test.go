package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/plugins_info.json
	pluginsInfoJSON string
	//go:embed internal/fixtures/plugin_info.json
	pluginInfoJSON string
)

func Test_client_GetPluginsInfo(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}

	t.Run("should be able to fetch all plugin info successfully", func(t *testing.T) {
		server := mockServer([]byte(pluginsInfoJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PluginsInfo{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Plugins: []*gocd.Plugin{
				{
					ID:                 "json.config.plugin",
					PluginFileLocation: "/Users/varshavs/gocd/server/plugins/bundled/gocd-json-config-plugin.jar",
					BundledPlugin:      true,
					About: map[string]interface{}{
						"name":                     "JSON Configuration Plugin",
						"version":                  "0.2",
						"target_go_version":        "16.1.0",
						"description":              "Configuration plugin that supports GoCD configuration in JSON",
						"target_operating_systems": []interface{}{},
						"vendor": map[string]interface{}{
							"name": "Tomasz Setkowski",
							"url":  "https://github.com/tomzo/gocd-json-config-plugin",
						},
					},
					Status: struct {
						State string `json:"state,omitempty" yaml:"state,omitempty"`
					}(struct{ State string }{State: "active"}),
					Extensions: []gocd.PluginAttributes{
						{
							Type: "configrepo",
							PluginSettings: &gocd.PluginSettingAttribute{
								Configurations: []*gocd.PluginConfiguration{
									{
										Key: "pipeline_pattern",
										Metadata: map[string]interface{}{
											"secure":   false,
											"required": false,
										},
									},
									{
										Key: "environment_pattern",
										Metadata: map[string]interface{}{
											"secure":   false,
											"required": false,
										},
									},
								},
							},
						},
					},
				},
			},
		}

		actual, err := client.GetPluginsInfo()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginsInfoJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PluginsInfo{}

		actual, err := client.GetPluginsInfo()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/plugin_info?include_bad=true\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pluginsInfoJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PluginsInfo{}

		actual, err := client.GetPluginsInfo()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/plugin_info?include_bad=true\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pluginsInfoJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PluginsInfo{}

		actual, err := client.GetPluginsInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.PluginsInfo{}

		actual, err := client.GetPluginsInfo()
		assert.EqualError(t, err, "call made to get all plugins info errored with: "+
			"Get \"http://localhost:8156/go/api/admin/plugin_info?include_bad=true\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetPluginInfo(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}

	t.Run("should be able to fetch all plugin info successfully", func(t *testing.T) {
		server := mockServer([]byte(pluginInfoJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Plugin{
			ETAG:               "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:                 "my_plugin",
			PluginFileLocation: "/path/to/server/plugins/external/my_plugin.jar",
			BundledPlugin:      false,
			About: map[string]interface{}{
				"name":                     "My Plugin",
				"version":                  "0.2",
				"target_go_version":        "16.1.0",
				"description":              "Short desc",
				"target_operating_systems": []interface{}{},
				"vendor": map[string]interface{}{
					"name": "GoCD contributors",
					"url":  "https://github.com/tomzo/gocd-json-config-plugin",
				},
			},
			Status: struct {
				State string `json:"state,omitempty" yaml:"state,omitempty"`
			}(struct{ State string }{State: "active"}),
			Extensions: []gocd.PluginAttributes{
				{
					Type: "configrepo",
					PluginSettings: &gocd.PluginSettingAttribute{
						Configurations: []*gocd.PluginConfiguration{
							{
								Key: "pipeline_pattern",
								Metadata: map[string]interface{}{
									"secure":   false,
									"required": false,
								},
							},
							{
								Key: "environment_pattern",
								Metadata: map[string]interface{}{
									"secure":   false,
									"required": false,
								},
							},
						},
					},
				},
			},
		}

		actual, err := client.GetPluginInfo("")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pluginInfoJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Plugin{}

		actual, err := client.GetPluginInfo("")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/plugin_info\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pluginInfoJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Plugin{}

		actual, err := client.GetPluginInfo("")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/plugin_info\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pluginInfoJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Plugin{}

		actual, err := client.GetPluginInfo("")
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all plugin info present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.Plugin{}

		actual, err := client.GetPluginInfo("my_plugin")
		assert.EqualError(t, err, "call made to get plugin info 'my_plugin' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/plugin_info/my_plugin\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
