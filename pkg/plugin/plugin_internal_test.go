package plugin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_SetURL(t *testing.T) {
	originalYAMLAPIURL := yamlPluginAPIURL
	originalJSONAPIURL := jsonPluginAPIURL
	originalGroovyAPIURL := groovyPluginAPIURL
	t.Cleanup(func() {
		yamlPluginAPIURL = originalYAMLAPIURL
		jsonPluginAPIURL = originalJSONAPIURL
		groovyPluginAPIURL = originalGroovyAPIURL
	})

	tests := []struct {
		name            string
		pipelineType    string
		version         string
		latestRelease   string
		apiURL          *string
		expectedVersion string
		expectedURL     string
	}{
		{
			name:            "yaml with configured version",
			pipelineType:    "yaml",
			version:         "0.13.0",
			expectedVersion: "0.13.0",
			expectedURL:     "https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.13.0/yaml-config-plugin-0.13.0.jar",
		},
		{
			name:            "json with configured version",
			pipelineType:    "json",
			version:         "0.6.0",
			expectedVersion: "0.6.0",
			expectedURL:     "https://github.com/tomzo/gocd-json-config-plugin/releases/download/0.6.0/json-config-plugin-0.6.0.jar",
		},
		{
			name:            "groovy with configured version",
			pipelineType:    "groovy",
			version:         "2.1.3-512",
			expectedVersion: "2.1.3-512",
			expectedURL:     "https://github.com/gocd-contrib/gocd-groovy-dsl-config-plugin/releases/download/v2.1.3-512/gocd-groovy-dsl-config-plugin-2.1.3-512.jar",
		},
		{
			name:            "yaml with latest version",
			pipelineType:    "yaml",
			latestRelease:   "0.14.0",
			apiURL:          &yamlPluginAPIURL,
			expectedVersion: "0.14.0",
			expectedURL:     "https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.14.0/yaml-config-plugin-0.14.0.jar",
		},
		{
			name:            "json with latest version",
			pipelineType:    "json",
			latestRelease:   "0.7.0",
			apiURL:          &jsonPluginAPIURL,
			expectedVersion: "0.7.0",
			expectedURL:     "https://github.com/tomzo/gocd-json-config-plugin/releases/download/0.7.0/json-config-plugin-0.7.0.jar",
		},
		{
			name:            "groovy with latest version",
			pipelineType:    "groovy",
			latestRelease:   "v2.2.0-100",
			apiURL:          &groovyPluginAPIURL,
			expectedVersion: "2.2.0-100",
			expectedURL:     "https://github.com/gocd-contrib/gocd-groovy-dsl-config-plugin/releases/download/v2.2.0-100/gocd-groovy-dsl-config-plugin-2.2.0-100.jar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.apiURL != nil {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, err := fmt.Fprintf(w, `[{"name":%q}]`, tt.latestRelease)
					require.NoError(t, err)
				}))
				t.Cleanup(server.Close)
				*tt.apiURL = server.URL
			}

			cfg := NewPluginConfig(tt.version, "", "", "debug").(*Config)
			cfg.PipelineType = tt.pipelineType

			err := cfg.setURL()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedVersion, cfg.Version)
			assert.Equal(t, tt.expectedURL, cfg.URL)
		})
	}
}

func TestConfig_SetURLSkipsExistingURL(t *testing.T) {
	cfg := NewPluginConfig("", "", "https://example.com/plugin.jar", "debug").(*Config)
	cfg.PipelineType = "yaml"

	err := cfg.setURL()
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/plugin.jar", cfg.URL)
	assert.Empty(t, cfg.Version)
}

func TestConfig_Accessors(t *testing.T) {
	cfg := NewPluginConfig("", "", "", "debug")

	err := cfg.SetType([]string{"pipeline.gocd.yaml"})
	require.NoError(t, err)
	assert.Equal(t, "yaml", cfg.GetType())

	cfg.SetVersion("0.13.0")
	assert.Equal(t, "0.13.0", cfg.GetVersion())
}

func TestConfig_GetLatestRelease(t *testing.T) {
	t.Run("returns latest tag", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, `[{"name":"v1.2.3"}]`)
			require.NoError(t, err)
		}))
		t.Cleanup(server.Close)

		cfg := NewPluginConfig("", "", "", "debug").(*Config)

		version, err := cfg.GetLatestRelease(server.URL)
		require.NoError(t, err)
		assert.Equal(t, "v1.2.3", version)
	})

	t.Run("returns non ok error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "nope", http.StatusTeapot)
		}))
		t.Cleanup(server.Close)

		cfg := NewPluginConfig("", "", "", "debug").(*Config)

		version, err := cfg.GetLatestRelease(server.URL)
		require.ErrorContains(t, err, "got 418 from GoCD while making GET call")
		assert.Empty(t, version)
	})

	t.Run("returns marshal error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, `{`)
			require.NoError(t, err)
		}))
		t.Cleanup(server.Close)

		cfg := NewPluginConfig("", "", "", "debug").(*Config)

		version, err := cfg.GetLatestRelease(server.URL)
		require.ErrorContains(t, err, "unexpected end of JSON input")
		assert.Empty(t, version)
	})
}

func TestConfig_DownloadUnitPaths(t *testing.T) {
	t.Run("returns configured local path", func(t *testing.T) {
		cfg := NewPluginConfig("", "/tmp/plugin.jar", "", "debug")

		pluginPath, err := cfg.Download()
		require.NoError(t, err)
		assert.Equal(t, "/tmp/plugin.jar", pluginPath)
	})

	t.Run("uses existing downloaded plugin", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		pluginDir := filepath.Join(home, ".gocd", "plugins")
		require.NoError(t, os.MkdirAll(pluginDir, 0o755))
		expectedPath := filepath.Join(pluginDir, "plugin.jar")
		require.NoError(t, os.WriteFile(expectedPath, []byte("plugin"), 0o644))

		cfg := NewPluginConfig("", "", "https://example.com/releases/plugin.jar", "debug").(*Config)

		pluginPath, err := cfg.Download()
		require.NoError(t, err)
		assert.Equal(t, expectedPath, pluginPath)
		assert.Equal(t, expectedPath, cfg.Path)
	})

	t.Run("downloads plugin", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := fmt.Fprint(w, "plugin")
			require.NoError(t, err)
		}))
		t.Cleanup(server.Close)

		cfg := NewPluginConfig("", "", server.URL+"/plugin.jar", "debug").(*Config)
		expectedPath := filepath.Join(home, ".gocd", "plugins", "plugin.jar")

		pluginPath, err := cfg.Download()
		require.NoError(t, err)
		assert.Equal(t, expectedPath, pluginPath)
		assert.Equal(t, expectedPath, cfg.Path)

		content, err := os.ReadFile(expectedPath)
		require.NoError(t, err)
		assert.Equal(t, "plugin", string(content))
	})

	t.Run("returns download non ok error", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "missing", http.StatusNotFound)
		}))
		t.Cleanup(server.Close)

		cfg := NewPluginConfig("", "", server.URL+"/plugin.jar", "debug")

		pluginPath, err := cfg.Download()
		require.EqualError(t, err, "downloading plugin returned non OK response code '404' with BODY: ''")
		assert.Empty(t, pluginPath)
	})
}
