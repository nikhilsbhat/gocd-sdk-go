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
	//go:embed internal/fixtures/elastic_agent_profiles.json
	elasticAgentProfilesJSON string
	//go:embed internal/fixtures/elastic_agent_profile.json
	elasticAgentProfileJSON string
	//go:embed internal/fixtures/elastic_agent_profiles_usage.json
	elasticAgentProfileUsageJSON string
)

func Test_client_GetElasticAgentProfiles(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to fetch the elastic agent profiles successfully", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfilesJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			CommonConfigs: []gocd.CommonConfig{
				{
					ID:               "unit-tests",
					ClusterProfileID: "prod-cluster",
					Properties: []gocd.PluginConfiguration{
						{
							Key:   "Image",
							Value: "alpine:latest",
						},
						{
							Key:   "Command",
							Value: "",
						},
						{
							Key:   "Environment",
							Value: "",
						},
						{
							Key:   "MaxMemory",
							Value: "200M",
						},
						{
							Key:   "ReservedMemory",
							Value: "150M",
						},
					},
				},
			},
		}

		actual, err := client.GetElasticAgentProfiles()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all elastic agent profiles present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfilesJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetElasticAgentProfiles()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/elastic/profiles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all elastic agent profiles present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfilesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetElasticAgentProfiles()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/elastic/profiles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all elastic agent profiles from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("elasticAgentProfilesJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetElasticAgentProfiles()
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all elastic agent profiles present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.ProfilesConfig{
			CommonConfigs: nil,
		}

		actual, err := client.GetElasticAgentProfiles()
		require.EqualError(t, err, "call made to get elastic agent profiles errored with: "+
			"Get \"http://localhost:8156/go/api/elastic/profiles\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetElasticAgentProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	profileName := "prod-cluster"

	t.Run("should be able to fetch a specific elastic agent profile successfully", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{
			ETAG:             "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:               "unit-tests",
			ClusterProfileID: "prod-cluster",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "Image",
					Value: "gocdcontrib/gocd-dev-build",
				},
				{
					Key:   "Environment",
					Value: "JAVA_HOME=/opt/java\nMAKE_OPTS=-j8",
				},
			},
		}

		actual, err := client.GetElasticAgentProfile(profileName)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfilesJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetElasticAgentProfile(profileName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/elastic/profiles/prod-cluster\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfilesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetElasticAgentProfile(profileName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/elastic/profiles/prod-cluster\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("elasticAgentProfilesJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetElasticAgentProfile(profileName)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{}

		actual, err := client.GetElasticAgentProfile(profileName)
		require.EqualError(t, err, "call made to get elastic agent profile 'prod-cluster' errored with: "+
			"Get \"http://localhost:8156/go/api/elastic/profiles/prod-cluster\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeleteElasticAgentProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	profileName := "prod-cluster"

	t.Run("should be able to delete an appropriate elastic agent profile successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteElasticAgentProfile(profileName)
		require.NoError(t, err)
	})

	t.Run("should error out while deleting an elastic agent profile due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteElasticAgentProfile(profileName)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/elastic/profiles/prod-cluster\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting an elastic agent profile due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteElasticAgentProfile(profileName)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/elastic/profiles/prod-cluster\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting an elastic agent profile as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteElasticAgentProfile(profileName)
		require.EqualError(t, err, "call made to delete elastic agent profile 'prod-cluster' errored with: "+
			"Delete \"http://localhost:8156/go/api/elastic/profiles/prod-cluster\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_UpdateElasticAgentProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to update an specific elastic agent profile successfully", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{
			ETAG:             "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:               "unit-tests",
			ClusterProfileID: "prod-cluster",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "Image",
					Value: "gocdcontrib/gocd-dev-build",
				},
				{
					Key:   "Environment",
					Value: "JAVA_HOME=/opt/java\nMAKE_OPTS=-j8",
				},
			},
		}

		expected := profileCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateElasticAgentProfile(profileCfg)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific elastic agent profile present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/elastic/profiles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific elastic agent profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/elastic/profiles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific elastic agent profile from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("elasticAgentProfileJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.UpdateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating an specific elastic agent profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		profileCfg := gocd.CommonConfig{ID: "unit-tests"}
		expected := gocd.CommonConfig{}

		actual, err := client.UpdateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "call made to update elastic agent profile 'unit-tests' errored with: "+
			"Put \"http://localhost:8156/go/api/elastic/profiles/unit-tests\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateElasticAgentProfile(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to create an specific elastic agent profile successfully", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		elasticAgentProfileCfg := gocd.CommonConfig{
			ETAG:             "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:               "unit-tests",
			ClusterProfileID: "prod-cluster",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "Image",
					Value: "gocdcontrib/gocd-dev-build",
				},
				{
					Key:   "Environment",
					Value: "JAVA_HOME=/opt/java\nMAKE_OPTS=-j8",
				},
			},
		}

		expected := elasticAgentProfileCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.CreateElasticAgentProfile(elasticAgentProfileCfg)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific elastic agent profile present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.CreateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/elastic/profiles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific elastic agent profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.CreateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/elastic/profiles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific elastic agent profile from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("elasticAgentProfileJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		profileCfg := gocd.CommonConfig{}
		expected := profileCfg

		actual, err := client.CreateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating an specific elastic agent profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		profileCfg := gocd.CommonConfig{ID: "unit-tests"}
		expected := gocd.CommonConfig{}

		actual, err := client.CreateElasticAgentProfile(profileCfg)
		require.EqualError(t, err, "call made to create elastic agent profile 'unit-tests' errored with: "+
			"Post \"http://localhost:8156/go/api/elastic/profiles\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetElasticAgentProfileUsage(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	profileName := "prod-cluster"

	t.Run("should be able to fetch usage of an elastic agent profile successfully", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileUsageJSON), http.StatusOK,
			correctArtifactHeader, true, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.ElasticProfileUsage{
			{
				PipelineName:         "pipeline-1",
				StageName:            "integration",
				JobName:              "build",
				PipelineConfigOrigin: "config_repo",
			},
			{
				PipelineName:         "pipeline-2",
				StageName:            "package",
				JobName:              "default",
				PipelineConfigOrigin: "config_repo",
			},
		}

		actual, err := client.GetElasticAgentProfileUsage(profileName)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should fail while fetching usage of an elastic agent profile successfully due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfileUsageJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetElasticAgentProfileUsage(profileName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/internal/elastic/profiles/prod-cluster/usages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(elasticAgentProfilesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetElasticAgentProfileUsage(profileName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/internal/elastic/profiles/prod-cluster/usages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile present in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("elasticAgentProfilesJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetElasticAgentProfileUsage(profileName)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching a specific elastic agent profile present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetElasticAgentProfileUsage(profileName)
		require.EqualError(t, err, "call made to get elastic agent profile usage 'prod-cluster' errored with: "+
			"Get \"http://localhost:8156/go/api/internal/elastic/profiles/prod-cluster/usages\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}
