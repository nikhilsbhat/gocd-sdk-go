package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/cluster_profiles.json
	clusterProfilesJSON string
	//go:embed internal/fixtures/cluster_profile.json
	clusterProfileJSON string
)

func Test_client_GetClusterProfiles(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to fetch the cluster profiles successfully", func(t *testing.T) {
		server := mockServer([]byte(clusterProfilesJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			CommonConfigs: []gocd.CommonConfig{
				{
					ID:       "prod-cluster",
					PluginID: "cd.go.contrib.elastic-agent.docker",
					Properties: []gocd.PluginConfiguration{
						{
							Key:   "GoServerUrl",
							Value: "https://ci.example.com/go",
						},
					},
				},
			},
		}

		actual, err := client.GetClusterProfiles()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all cluster profiles present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfilesJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetClusterProfiles()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all cluster profiles present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfilesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetClusterProfiles()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all cluster profiles from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("clusterProfilesJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetClusterProfiles()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all cluster profiles present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetClusterProfiles()
		assert.EqualError(t, err, "call made to get cluster profiles errored with: "+
			"Get \"http://localhost:8156/go/api/admin/elastic/cluster_profiles\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetClusterProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	profileName := "prod-cluster"
	t.Run("should be able to fetch a specific cluster profile successfully", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "prod-cluster",
			PluginID: "cd.go.contrib.elastic-agent.docker",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "GoServerUrl",
					Value: "https://ci.example.com/go",
				},
			},
		}

		actual, err := client.GetClusterProfile(profileName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific cluster profile present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetClusterProfile(profileName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific cluster profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetClusterProfile(profileName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific cluster profile from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("clusterProfileJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetClusterProfile(profileName)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific cluster profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{}

		actual, err := client.GetClusterProfile(profileName)
		assert.EqualError(t, err, "call made to get cluster profile 'prod-cluster' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/elastic/cluster_profiles/prod-cluster\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateClusterProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to create a specific cluster profile successfully", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "prod-cluster",
			PluginID: "cd.go.contrib.elastic-agent.docker",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "GoServerUrl",
					Value: "https://ci.example.com/go",
				},
			},
		}

		expected := profileCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.CreateClusterProfile(profileCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific cluster profile present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateClusterProfile(profileCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific cluster profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.CreateClusterProfile(profileCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific cluster profile from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("clusterProfileJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.CreateClusterProfile(profileCfg)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific cluster profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		profileCfg := gocd.CommonConfig{ID: "prod-cluster"}
		expected := gocd.CommonConfig{}

		actual, err := client.CreateClusterProfile(profileCfg)
		assert.EqualError(t, err, "call made to create cluster profile 'prod-cluster' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/elastic/cluster_profiles\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateClusterProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to update a specific cluster profile successfully", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{
			ETAG:     "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:       "prod-cluster",
			PluginID: "cd.go.contrib.elastic-agent.docker",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "GoServerUrl",
					Value: "https://ci.example.com/go",
				},
			},
		}

		expected := profileCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateClusterProfile(profileCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific cluster profile present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateClusterProfile(profileCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific cluster profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateClusterProfile(profileCfg)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific cluster profile from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("clusterProfileJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateClusterProfile(profileCfg)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific cluster profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		profileCfg := gocd.CommonConfig{ID: "prod-cluster"}
		expected := gocd.CommonConfig{}

		actual, err := client.UpdateClusterProfile(profileCfg)
		assert.EqualError(t, err, "call made to update cluster profile 'prod-cluster' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/elastic/cluster_profiles/prod-cluster\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeleteClusterProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to delete an appropriate cluster profile successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteClusterProfile("prod-cluster")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting a cluster profile due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteClusterProfile("prod-cluster")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting a cluster profile due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteClusterProfile("prod-cluster")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting a cluster profile as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteClusterProfile("prod-cluster")
		assert.EqualError(t, err, "call made to delete cluster profile 'prod-cluster' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/elastic/cluster_profiles/prod-cluster\": dial tcp [::1]:8156: connect: connection refused")
	})
}
