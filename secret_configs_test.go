package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/fixtures/secrets_configs.json
	secretConfigsJSON string
	//go:embed internal/fixtures/secrets_config.json
	secretConfigJSON string
)

func Test_client_GetSecretConfigs(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to fetch the secret configs successfully", func(t *testing.T) {
		server := mockServer([]byte(secretConfigsJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SecretsConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			CommonConfigs: []gocd.CommonConfig{
				{
					ID:       "demo",
					PluginID: "cd.go.secrets.file-based-plugin",
					Properties: []gocd.PluginConfiguration{
						{
							Key:   "SecretsFilePath",
							Value: "path/to/secret/file.db",
						},
					},
					Rules: []map[string]string{
						{
							"directive": "allow",
							"action":    "refer",
							"type":      "pipeline_group",
							"resource":  "first",
						},
					},
				},
			},
		}

		actual, err := client.GetSecretConfigs()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all secret configs present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SecretsConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetSecretConfigs()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/secret_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all secret configs present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SecretsConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetSecretConfigs()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/secret_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all secret configs from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("secretConfigsJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SecretsConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetSecretConfigs()
		require.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all secret configs present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.SecretsConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetSecretConfigs()
		require.EqualError(t, err, "call made to get secret configs errored with: "+
			"Get \"http://localhost:8156/go/api/admin/secret_configs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetSecretConfig(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to fetch a specific secret config successfully", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "demo",
			PluginID: "cd.go.secrets.file-based-plugin",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "SecretsFilePath",
					Value: "path/to/secret/file.db",
				},
			},
			Rules: []map[string]string{
				{
					"directive": "allow",
					"action":    "refer",
					"type":      "pipeline_group",
					"resource":  "first",
				},
			},
		}

		actual, err := client.GetSecretConfig("demo")
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific secret config present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetSecretConfig("demo")
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/secret_configs/demo\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific secret config present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetSecretConfig("demo")
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/secret_configs/demo\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific secret config from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("secretConfigJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetSecretConfig("demo")
		require.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific secret config present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{}

		actual, err := client.GetSecretConfig("demo")
		require.EqualError(t, err, "call made to get secret config 'demo' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/secret_configs/demo\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeleteSecretConfig(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	secretCfgID := "demo"

	t.Run("should be able to delete an appropriate secret config successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteSecretConfig(secretCfgID)
		require.NoError(t, err)
	})

	t.Run("should error out while deleting an secret config due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteSecretConfig(secretCfgID)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/secret_configs/demo\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting an secret config due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteSecretConfig(secretCfgID)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/secret_configs/demo\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting an secret config as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteSecretConfig(secretCfgID)
		require.EqualError(t, err, "call made to delete secret config 'demo' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/secret_configs/demo\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_UpdateSecretConfig(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	secretCfgID := "demo"

	t.Run("should be able to update an specific secret config successfully", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       secretCfgID,
			PluginID: "cd.go.secrets.file-based-plugin",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "SecretsFilePath",
					Value: "path/to/secret/file.db",
				},
			},
			Rules: []map[string]string{
				{
					"directive": "allow",
					"action":    "refer",
					"type":      "pipeline_group",
					"resource":  "first",
				},
			},
		}

		expected := profileCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateSecretConfig(profileCfg)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific secret config present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateSecretConfig(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/secret_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific secret config present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateSecretConfig(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/secret_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific secret config from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("secretConfigJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		secretCfg := gocd.CommonConfig{}
		expected := secretCfg

		actual, err := client.UpdateSecretConfig(secretCfg)
		require.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific secret config present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		profileCfg := gocd.CommonConfig{ID: secretCfgID}
		expected := gocd.CommonConfig{}

		actual, err := client.UpdateSecretConfig(profileCfg)
		require.EqualError(t, err, "call made to update secret config 'demo' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/secret_configs/demo\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateSecretConfig(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	secretCfgID := "demo"

	t.Run("should be able to create an specific secret config successfully", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		secretCfg := gocd.CommonConfig{
			ID:       secretCfgID,
			PluginID: "cd.go.secrets.file-based-plugin",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "SecretsFilePath",
					Value: "path/to/secret/file.db",
				},
			},
			Rules: []map[string]string{
				{
					"directive": "allow",
					"action":    "refer",
					"type":      "pipeline_group",
					"resource":  "first",
				},
			},
		}

		expected := secretCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.CreateSecretConfig(secretCfg)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific secret config present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		secretCfg := gocd.CommonConfig{}
		expected := secretCfg

		actual, err := client.CreateSecretConfig(secretCfg)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/secret_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific secret config present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(secretConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.CreateSecretConfig(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/secret_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific secret config from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("secretConfigJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		secretCfg := gocd.CommonConfig{}
		expected := secretCfg

		actual, err := client.CreateSecretConfig(secretCfg)
		require.EqualError(t, err, "reading response body errored with: invalid character 's' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific secret config present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		profileCfg := gocd.CommonConfig{ID: secretCfgID}
		expected := gocd.CommonConfig{}

		actual, err := client.CreateSecretConfig(profileCfg)
		require.EqualError(t, err, "call made to create secrets config 'demo' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/secret_configs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
