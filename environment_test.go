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

		actual, err := client.GetEnvironmentInfo()
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

		actual, err := client.GetEnvironmentInfo()
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

		actual, err := client.GetEnvironmentInfo()
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
			},
			{
				Name: "foobar2",
				Pipelines: []gocd.Pipeline{
					{
						Name: "pipeline2",
					},
				},
			},
		}

		actual, err := client.GetEnvironmentInfo()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
