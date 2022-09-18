package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/environments.json
var environmentJSON string

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
		assert.EqualError(t, err, "call made to get environment errored with "+
			"Get \"http://localhost:8156/go/api/admin/environments\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all config repos present as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctEnvHeader, false)
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
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctEnvHeader, false)
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

	t.Run("should be able to fetch all config repos present in GoCD server", func(t *testing.T) {
		server := mockServer([]byte(environmentJSON), http.StatusOK, correctEnvHeader, false)
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
				EnvVars: []struct {
					Name           string `json:"name,omitempty"`
					Value          string `json:"value,omitempty"`
					EncryptedValue string `json:"encrypted_value,omitempty"`
					Secure         bool   `json:"secure,omitempty"`
				}{
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
				EnvVars: []struct {
					Name           string `json:"name,omitempty"`
					Value          string `json:"value,omitempty"`
					EncryptedValue string `json:"encrypted_value,omitempty"`
					Secure         bool   `json:"secure,omitempty"`
				}{
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
		server := mockServer([]byte(encryptionJSON), http.StatusOK, correctEnvHeader, false)
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
			EnvVars: []struct {
				Name           string `json:"name,omitempty"`
				Value          string `json:"value,omitempty"`
				EncryptedValue string `json:"encrypted_value,omitempty"`
				Secure         bool   `json:"secure,omitempty"`
			}{
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
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}, false)
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
		server := mockServer(nil, http.StatusOK, nil, false)
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
