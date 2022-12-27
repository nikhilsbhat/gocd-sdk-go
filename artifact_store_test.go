package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/artifact_stores.json
	artifactStoresJSON string
	//go:embed internal/fixtures/artifact_store.json
	artifactStoreJSON string
)

func Test_client_GetArtifactStores(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to fetch the artifact stores successfully", func(t *testing.T) {
		server := mockServer([]byte(artifactStoresJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.ArtifactStoresConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			CommonConfigs: []gocd.CommonConfig{
				{
					ID:       "hub.docker",
					PluginID: "cd.go.artifact.docker.registry",
					Properties: []gocd.PluginConfiguration{
						{
							Key:   "RegistryURL",
							Value: "https://your_docker_registry_url",
						},
						{
							Key:   "Username",
							Value: "admin",
						},
						{
							Key:            "Password",
							EncryptedValue: "AES:tdfTtYtIUSAF2JXJP/3YwA==:43Kjidjuh42NHKisCAs/BQ==",
						},
					},
				},
			},
		}

		actual, err := client.GetArtifactStores()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all artifact stores present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoresJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.ArtifactStoresConfig{
			CommonConfigs: nil,
			ETAG:          "",
		}

		actual, err := client.GetArtifactStores()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all artifact stores present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoresJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.ArtifactStoresConfig{
			CommonConfigs: nil,
			ETAG:          "",
		}

		actual, err := client.GetArtifactStores()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching artifact stores from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("artifactStoreJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.ArtifactStoresConfig{
			CommonConfigs: nil,
			ETAG:          "",
		}

		actual, err := client.GetArtifactStores()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching artifact stores present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.ArtifactStoresConfig{
			CommonConfigs: nil,
			ETAG:          "",
		}

		actual, err := client.GetArtifactStores()
		assert.EqualError(t, err, "call made to get artifact stores errored with: "+
			"Get \"http://localhost:8156/go/api/admin/artifact_stores\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetArtifactStore(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to fetch an appropriate artifact store successfully", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "hub.docker",
			PluginID: "cd.go.artifact.docker.registry",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "RegistryURL",
					Value: "https://your_docker_registry_url",
				},
				{
					Key:   "Username",
					Value: "admin",
				},
				{
					Key:            "Password",
					EncryptedValue: "AES:tdfTtYtIUSAF2JXJP/3YwA==:43Kjidjuh42NHKisCAs/BQ==",
				},
			},
		}

		actual, err := client.GetArtifactStore("hub.docker")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching an appropriate artifact store due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.CommonConfig{
			Properties: nil,
		}

		actual, err := client.GetArtifactStore("hub.docker")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching an appropriate artifact store due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.CommonConfig{
			Properties: nil,
		}

		actual, err := client.GetArtifactStore("hub.docker")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching an appropriate artifact store as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("artifactStoreJSON"), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.CommonConfig{
			Properties: nil,
		}

		actual, err := client.GetArtifactStore("hub.docker")
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching an appropriate artifact store as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{
			Properties: nil,
		}

		actual, err := client.GetArtifactStore("hub.docker")
		assert.EqualError(t, err, "call made to get artifact store hub.docker errored with: "+
			"Get \"http://localhost:8156/go/api/admin/artifact_stores/hub.docker\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateArtifactStore(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to create an appropriate artifact store successfully", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "hub.docker",
			PluginID: "cd.go.artifact.docker.registry",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "RegistryURL",
					Value: "https://your_docker_registry_url",
				},
				{
					Key:   "Username",
					Value: "admin",
				},
				{
					Key:            "Password",
					EncryptedValue: "AES:tdfTtYtIUSAF2JXJP/3YwA==:43Kjidjuh42NHKisCAs/BQ==",
				},
			},
		}

		expected := storeCfg
		expected.ETAG = "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"

		actual, err := client.CreateArtifactStore(storeCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an artifact store due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{}

		actual, err := client.CreateArtifactStore(storeCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, storeCfg, actual)
	})

	t.Run("should error out while creating an artifact store due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{}

		actual, err := client.CreateArtifactStore(storeCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, storeCfg, actual)
	})

	t.Run("should error out while creating an artifact store as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("artifactStoreJSON"), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{}

		actual, err := client.CreateArtifactStore(storeCfg)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, storeCfg, actual)
	})

	t.Run("should error out while creating an artifact store as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		storeCfg := gocd.CommonConfig{ID: "docker"}

		actual, err := client.CreateArtifactStore(storeCfg)
		assert.EqualError(t, err, "call made to create artifact store docker errored with: "+
			"Post \"http://localhost:8156/go/api/admin/artifact_stores\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.CommonConfig{}, actual)
	})
}

func Test_client_UpdateArtifactStore(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to update an appropriate artifact store successfully", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "hub.docker",
			PluginID: "cd.go.artifact.docker.registry",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "RegistryURL",
					Value: "https://your_docker_registry_url",
				},
				{
					Key:   "Username",
					Value: "admin",
				},
				{
					Key:            "Password",
					EncryptedValue: "AES:tdfTtYtIUSAF2JXJP/3YwA==:43Kjidjuh42NHKisCAs/BQ==",
				},
			},
		}

		expected := storeCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateArtifactStore(storeCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an artifact store due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{}

		actual, err := client.UpdateArtifactStore(storeCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, storeCfg, actual)
	})

	t.Run("should error out while updating an artifact store due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{}

		actual, err := client.UpdateArtifactStore(storeCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, storeCfg, actual)
	})

	t.Run("should error out while updating an artifact store as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("artifactStoreJSON"), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		storeCfg := gocd.CommonConfig{}

		actual, err := client.UpdateArtifactStore(storeCfg)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, storeCfg, actual)
	})

	t.Run("should error out while updating an artifact store as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		storeCfg := gocd.CommonConfig{ID: "docker"}

		actual, err := client.UpdateArtifactStore(storeCfg)
		assert.EqualError(t, err, "call made to update artifact store docker errored with: "+
			"Put \"http://localhost:8156/go/api/admin/artifact_stores/docker\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.CommonConfig{}, actual)
	})
}

func Test_client_DeleteArtifactStore(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to delete an appropriate artifact store successfully", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeleteArtifactStore("docker")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting an artifact store due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeleteArtifactStore("docker")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting an artifact store due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoreJSON), http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeleteArtifactStore("docker")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting an artifact store as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteArtifactStore("docker")
		assert.EqualError(t, err, "call made to delete artifact store docker errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/artifact_stores/docker\": dial tcp [::1]:8156: connect: connection refused")
	})
}
