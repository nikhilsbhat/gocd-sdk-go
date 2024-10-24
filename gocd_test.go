package gocd_test

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestGetGoCDMethodNames(t *testing.T) {
	t.Run("should list all method names", func(t *testing.T) {
		response := gocd.GetGoCDMethodNames()
		assert.Equal(t, 144, len(response))
		assert.Equal(t, "AgentKillTask", response[0])
		assert.Equal(t, "UpdatePipelineGroup", response[137])
	})
}

func TestNewClient(t *testing.T) {
	t.Run("should be able to use token based auth", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)

		auth = gocd.Auth{
			BearerToken: "dlskjmxelkjmxwlerkmwelkfmwek",
		}

		client := gocd.NewClient(server.URL, auth, "info", nil)

		assert.NotNil(t, client)
	})

	t.Run("should be able to use token with CA based auth", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)

		caContent, err := os.ReadFile("internal/fixtures/ca.sample.pem")
		require.NoError(t, err)

		auth = gocd.Auth{
			BearerToken: "dlskjmxelkjmxwlerkmwelkfmwek",
		}

		client := gocd.NewClient(server.URL, auth, "info", caContent)

		assert.NotNil(t, client)
	})
}
