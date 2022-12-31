package gocd_test

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/environments.json
	environmentsJSON string
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
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetEnvironments()
		assert.EqualError(t, err, "call made to get environments errored with "+
			"Get \"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all config repos present as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctEnvHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetEnvironments()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all config repos present as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetEnvironments()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch all environment information present in GoCD server", func(t *testing.T) {
		server := mockServer([]byte(environmentsJSON), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

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
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateEnvironments(t *testing.T) {
	correctEnvHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to create the environment successfully", func(t *testing.T) {
		server := mockServer([]byte(encryptionJSON), http.StatusOK, correctEnvHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

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
		assert.NoError(t, err)
	})

	t.Run("should error out while creating environment due to wrong headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		environment := gocd.Environment{}

		err := client.CreateEnvironment(environment)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while creating environment due to missing", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		environment := gocd.Environment{}

		err := client.CreateEnvironment(environment)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while creating environment due to missing", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		environment := gocd.Environment{}
		err := client.CreateEnvironment(environment)
		assert.EqualError(t, err, "call made to create environment errored with Post"+
			" \"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_DeleteEnvironment(t *testing.T) {
	t.Run("should be able to delete the environment successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeleteEnvironment("env1")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting the environment as wrong headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeleteEnvironment("env1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting the environment as no headers set", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeleteEnvironment("env1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting the environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteEnvironment("env1")
		assert.EqualError(t, err, "call made to delete environment env1 errored with "+
			"Delete \"http://localhost:8156/go/api/admin/environments/env1\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_UpdateEnvironment(t *testing.T) {
	correctUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON, "If-Match": "26b227605daf6f2d7768c8edaf61b861"}
	t.Run("should be able to update the environment successfully", func(t *testing.T) {
		server := mockServer([]byte(environmentUpdateJSON), http.StatusOK, correctUpdateHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		environment := gocd.Environment{
			Name: "new_environment",
			Pipelines: []gocd.Pipeline{{
				Name: "up42",
			}},
			ETAG: "26b227605daf6f2d7768c8edaf61b861",
		}

		var expected gocd.Environment
		unMarshallErr := json.Unmarshal([]byte(environmentUpdateJSON), &expected)
		assert.NoError(t, unMarshallErr)

		actual, err := client.UpdateEnvironment(environment)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating the environment due to wrong headers set", func(t *testing.T) {
		server := mockServer([]byte(environmentUpdateJSON), http.StatusOK,
			map[string]string{
				"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON,
				"If-Match": "26b227605daf6f2d7768c8edaf61b861",
			}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		environment := gocd.Environment{}
		actual, err := client.UpdateEnvironment(environment)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while updating the environment due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(environmentUpdateJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		environment := gocd.Environment{}
		actual, err := client.UpdateEnvironment(environment)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while updating the environment as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("environmentUpdateJSON"), http.StatusOK, correctUpdateHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		environment := gocd.Environment{
			ETAG: "26b227605daf6f2d7768c8edaf61b861",
		}
		actual, err := client.UpdateEnvironment(environment)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while updating the environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		environment := gocd.Environment{}

		actual, err := client.UpdateEnvironment(environment)
		assert.EqualError(t, err, "call made to update environment errored with Patch "+
			"\"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Environment{}, actual)
	})
}

func Test_client_PatchEnvironment(t *testing.T) {
	correctPatchHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to patch the environment successfully", func(t *testing.T) {
		server := mockServer([]byte(environmentPatchJSON), http.StatusOK, correctPatchHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		patch := gocd.PatchEnvironment{
			Name: "new_environment",
			Pipelines: struct {
				Add    []string `json:"add,omitempty"`
				Remove []string `json:"remove,omitempty"`
			}{
				Add: []string{"up42"},
			},
			EnvVars: struct {
				Add []struct {
					Name  string `json:"name,omitempty"`
					Value string `json:"value,omitempty"`
				} `json:"add,omitempty"`
				Remove []string `json:"remove,omitempty"`
			}{
				Add: []struct {
					Name  string `json:"name,omitempty"`
					Value string `json:"value,omitempty"`
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
		assert.NoError(t, unMarshallErr)

		actual, err := client.PatchEnvironment(patch)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while patching GoCD environment as required headers are missing", func(t *testing.T) {
		server := mockServer([]byte(environmentPatchJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)
		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while patching GoCD environment as wrong headers passed", func(t *testing.T) {
		server := mockServer([]byte(environmentPatchJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)
		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while patching GoCD environment as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("environmentPatchJSON"), http.StatusOK, correctPatchHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)
		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while patching GoCD environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		patch := gocd.PatchEnvironment{}

		actual, err := client.PatchEnvironment(patch)
		assert.EqualError(t, err, "call made to patch environment errored with Patch "+
			"\"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Environment{}, actual)
	})
}

func Test_client_GetEnvironment(t *testing.T) {
	correctGetHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	envName := "my_environment"

	t.Run("should be able to fetch environment config from GoCD server successfully", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, correctGetHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

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
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching GoCD environment as wrong headers set", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetEnvironment(envName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while fetching GoCD environment as headers are missing", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetEnvironment(envName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while fetching GoCD environment as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("environmentJSON"), http.StatusOK, correctGetHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetEnvironment(envName)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'e' looking for beginning of value")
		assert.Equal(t, gocd.Environment{}, actual)
	})

	t.Run("should error out while fetching GoCD environment as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetEnvironment(envName)
		assert.EqualError(t, err, "call made to get environment errored with Get "+
			"\"http://localhost:8156/go/api/admin/environments/my_environment\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.Environment{}, actual)
	})
}
