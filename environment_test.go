package gocd_test

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/fixtures/environments.json
	environmentsJSON string
	//go:embed internal/fixtures/environment_merged.json
	environmentsMergedJSON string
	//go:embed internal/fixtures/environment.json
	environmentJSON string
	//go:embed internal/fixtures/environment_update.json
	environmentUpdateJSON string
	//go:embed internal/fixtures/environment_patch.json
	environmentPatchJSON string
)

func Test_client_GetEnvironmentInfo(t *testing.T) {
	correctEnvHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should error out while fetching all config repos present from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetEnvironments()
		require.EqualError(t, err, "call made to get environments errored with: "+
			"Get \"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all config repos present as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironments()
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/admin/environments\nwith BODY:backupJSON")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all config repos present as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironments()
		require.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch all environment information present in GoCD server", func(t *testing.T) {
		server := mockServer([]byte(environmentsJSON), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Environment{
			{
				Name: "foobar1",
				Pipelines: []gocd.Pipeline{
					{
						Name: "pipeline1",
					},
				},
				EnvVars: []gocd.EnvVars{
					{
						Name:   "username",
						Value:  "admin",
						Secure: false,
					},
					{
						Name:           "password",
						EncryptedValue: "LSd1TI0eLa+DjytHjj0qjA==",
						Secure:         true,
					},
				},
			},
			{
				Name: "foobar2",
				Pipelines: []gocd.Pipeline{
					{
						Name: "pipeline2",
					},
				},
				EnvVars: []gocd.EnvVars{
					{
						Name:   "username",
						Value:  "admin",
						Secure: false,
					},
					{
						Name:           "password",
						EncryptedValue: "LSd1TI0eLa+DjytHjj0qjA==",
						Secure:         true,
					},
				},
			},
		}

		actual, err := client.GetEnvironments()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateEnvironments(t *testing.T) {
	correctEnvHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to create the environment successfully", func(t *testing.T) {
		server := mockServer([]byte(encryptionJSON), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{
			Name: "environment_1",
			Pipelines: []gocd.Pipeline{
				{
					Name: "pipeline1",
				},
			},
			EnvVars: []gocd.EnvVars{
				{
					Name:   "env1",
					Value:  "env_value_1",
					Secure: false,
				},
				{
					Name:           "env2",
					EncryptedValue: "ksd64675xd-023-0-0293r0",
					Secure:         true,
				},
			},
		}

		err := client.CreateEnvironment(environment)
		require.NoError(t, err)
	})

	t.Run("should error out while creating environment due to wrong headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{}

		err := client.CreateEnvironment(environment)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/environments\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating environment due to missing", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{}

		err := client.CreateEnvironment(environment)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/environments\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating environment due to missing", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		environment := gocd.Environment{Name: "test"}
		err := client.CreateEnvironment(environment)
		require.EqualError(t, err, "call made to create environment 'test' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_DeleteEnvironment(t *testing.T) {
	t.Run("should be able to delete the environment successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteEnvironment("env1")
		require.NoError(t, err)
	})

	t.Run("should error out while deleting the environment as wrong headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteEnvironment("env1")
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/environments/env1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting the environment as no headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteEnvironment("env1")
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/environments/env1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting the environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteEnvironment("env1")
		require.EqualError(t, err, "call made to delete environment 'env1' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/environments/env1\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_UpdateEnvironment(t *testing.T) {
	correctUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON, "If-Match": "26b227605daf6f2d7768c8edaf61b861"}

	t.Run("should be able to update the environment successfully", func(t *testing.T) {
		server := mockServer([]byte(environmentUpdateJSON), http.StatusOK, correctUpdateHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{
			Name: "new_environment",
			Pipelines: []gocd.Pipeline{{
				Name: "up42",
			}},
			ETAG: "26b227605daf6f2d7768c8edaf61b861",
		}

		var expected gocd.Environment
		unMarshallErr := json.Unmarshal([]byte(environmentUpdateJSON), &expected)
		require.NoError(t, unMarshallErr)

		actual, err := client.UpdateEnvironment(environment)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating the environment due to wrong headers set", func(t *testing.T) {
		server := mockServer([]byte(environmentUpdateJSON), http.StatusOK,
			map[string]string{
				"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON,
				"If-Match": "26b227605daf6f2d7768c8edaf61b861",
			}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{}
		actual, err := client.UpdateEnvironment(environment)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/environments\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while updating the environment due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(environmentUpdateJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{}
		actual, err := client.UpdateEnvironment(environment)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/environments\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while updating the environment as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("environmentUpdateJSON"), http.StatusOK, correctUpdateHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		environment := gocd.Environment{
			ETAG: "26b227605daf6f2d7768c8edaf61b861",
		}
		actual, err := client.UpdateEnvironment(environment)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while updating the environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		environment := gocd.Environment{Name: "test"}

		actual, err := client.UpdateEnvironment(environment)
		require.EqualError(t, err, "call made to update environment 'test' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/environments/test\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Environment{}, actual)
	})
}

func Test_client_PatchEnvironment(t *testing.T) {
	correctPatchHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to patch the environment successfully", func(t *testing.T) {
		server := mockServer([]byte(environmentPatchJSON), http.StatusOK, correctPatchHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		patch := gocd.PatchEnvironment{
			Name: "new_environment",
			Pipelines: struct {
				Add    []string `json:"add,omitempty" yaml:"add,omitempty"`
				Remove []string `json:"remove,omitempty" yaml:"remove,omitempty"`
			}{
				Add: []string{"up42"},
			},
			EnvVars: struct {
				Add []struct {
					Name  string `json:"name,omitempty" yaml:"name,omitempty"`
					Value string `json:"value,omitempty" yaml:"value,omitempty"`
				} `json:"add,omitempty" yaml:"add,omitempty"`
				Remove []string `json:"remove,omitempty" yaml:"remove,omitempty"`
			}{
				Add: []struct {
					Name  string `json:"name,omitempty" yaml:"name,omitempty"`
					Value string `json:"value,omitempty" yaml:"value,omitempty"`
				}{
					{
						Name:  "GO_SERVER_URL",
						Value: "https://ci.example.com/go",
					},
				},
			},
		}

		var expected gocd.Environment
		unMarshallErr := json.Unmarshal([]byte(environmentPatchJSON), &expected)
		require.NoError(t, unMarshallErr)

		actual, err := client.PatchEnvironment(patch)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while patching GoCD environment as required headers are missing", func(t *testing.T) {
		server := mockServer([]byte(environmentPatchJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		require.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/admin/environments\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while patching GoCD environment as wrong headers passed", func(t *testing.T) {
		server := mockServer([]byte(environmentPatchJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		require.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/admin/environments\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while patching GoCD environment as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("environmentPatchJSON"), http.StatusOK, correctPatchHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while patching GoCD environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		patch := gocd.PatchEnvironment{Name: "test"}

		actual, err := client.PatchEnvironment(patch)
		require.EqualError(t, err, "call made to patch environment 'test' errored with: "+
			"Patch \"http://localhost:8156/go/api/admin/environments/test\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Environment{}, actual)
	})
}

func Test_client_GetEnvironment(t *testing.T) {
	correctGetHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	envName := "my_environment"

	t.Run("should be able to fetch environment config from GoCD server successfully", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, correctGetHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Environment{
			Name:      envName,
			Pipelines: []gocd.Pipeline{{Name: "up42"}},
			EnvVars: []gocd.EnvVars{
				{
					Name:   "username",
					Secure: false,
					Value:  "admin",
				}, {
					Name:           "password",
					Secure:         true,
					EncryptedValue: "LSd1TI0eLa+DjytHjj0qjA==",
				},
			},
		}

		actual, err := client.GetEnvironment(envName)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching GoCD environment as wrong headers set", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironment(envName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/environments/my_environment\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while fetching GoCD environment as headers are missing", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironment(envName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/environments/my_environment\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while fetching GoCD environment as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("environmentJSON"), http.StatusOK, correctGetHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironment(envName)
		require.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while fetching GoCD environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetEnvironment(envName)
		require.EqualError(t, err, "call made to get environment 'my_environment' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/environments/my_environment\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Environment{}, actual)
	})
}

func Test_client_GetEnvironmentMappings(t *testing.T) {
	correctEnvHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should error out while fetching all selected environment mappings from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetEnvironmentsMerged([]string{"example_environment"})
		require.EqualError(t, err, "call made to get environment mapping of 'example_environment' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/internal/environments/merged\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all selected environment mappings returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironmentsMerged([]string{"example_environment"})
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/admin/internal/environments/merged\nwith BODY:backupJSON")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all selected environment mappings present as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironmentsMerged([]string{"example_environment"})
		require.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all selected environment mappings present in GoCD server since environment was missing", func(t *testing.T) {
		server := mockServer([]byte(environmentsMergedJSON), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetEnvironmentsMerged([]string{"sample"})
		require.EqualError(t, err, "no environments found with names 'sample' to get mappings")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch all selected environment mappings present in GoCD server", func(t *testing.T) {
		server := mockServer([]byte(environmentsMergedJSON), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Environment{
			{
				Name: "sample_environment",
				Pipelines: []gocd.Pipeline{
					{Name: "gocd-prometheus-exporter"},
					{Name: "helm-images"},
					{Name: "helm-drift"},
					{Name: "gocd-cli"},
					{Name: "gocd-golang-sdk"},
				},
				EnvVars: []gocd.EnvVars{
					{Name: "TEST_ENV13", Value: "value_env13", EncryptedValue: "", Secure: false},
					{Name: "TEST_ENV12", Value: "value_env18", EncryptedValue: "", Secure: false},
					{Name: "TEST_ENV11", Value: "value_env11", EncryptedValue: "", Secure: false},
				},
				Origins: []gocd.EnvironmentOrigin{
					{Type: "gocd", ID: ""},
					{Type: "config_repo", ID: "sample"},
				},
				ETAG: "",
			},
			{
				Name: "example_environment",
				Pipelines: []gocd.Pipeline{
					{Name: "gocd-prometheus-exporter"},
				},
				EnvVars: []gocd.EnvVars{
					{Name: "TEST_ENV13", Value: "value_env13", EncryptedValue: "", Secure: false},
				},
				Origins: []gocd.EnvironmentOrigin{
					{Type: "gocd", ID: ""},
				},
				ETAG: "",
			},
		}

		actual, err := client.GetEnvironmentsMerged([]string{"sample_environment", "example_environment"})
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
