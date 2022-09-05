package gocd_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/config_repos.json
var configReposJSON string

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should error out while fetching config repo information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepoInfo()
		assert.EqualError(t, err, "call made to get config repo errored with "+
			"Get \"http://localhost:8153/go/api/admin/config_repos\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repo information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepoInfo()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repo information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepoInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able retrieve config repo information", func(t *testing.T) {
		server := mockServer([]byte(configReposJSON), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepoInfo()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(actual))
	})
}

func Test_client_CreateConfigRepoInfo(t *testing.T) {
	t.Run("server should return internal server error as malformed json passed", func(t *testing.T) {
		server := configRepoServer(configReposJSON)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		configRepo := &gocd.ConfigRepo{}
		err := client.CreateConfigRepoInfo(*configRepo)
		assert.EqualError(t, err, "body: json: cannot unmarshal string into Go value of type gocd.ConfigRepo httpcode: 500")
	})

	t.Run("should error out while making client call to create config repo", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		configRepo := &gocd.ConfigRepo{}
		err := client.CreateConfigRepoInfo(*configRepo)
		assert.EqualError(t, err, "post call made to create config repo errored with: "+
			"Post \"http://localhost:8153/go/api/admin/config_repos\": dial tcp 127.0.0.1:8153: connect: connection refused")
	})

	t.Run("should error out while creating config repo", func(t *testing.T) {
		configRepo := &gocd.ConfigRepo{
			PluginID:      "json.config.plugin",
			ID:            "repo1",
			Configuration: nil,
			Rules:         nil,
		}
		configRepo.Material.Type = "git"
		configRepo.Material.Attributes.URL = "https://github.com/config-repo/gocd-json-config-example.git"
		configRepo.Material.Attributes.AutoUpdate = false
		configRepo.Material.Attributes.Branch = "master"
		configRepo.Rules = []map[string]interface{}{
			{
				"directive": "allow",
				"action":    "refer",
				"type":      "pipeline_group",
				"resource":  "*",
			},
		}
		configRepo.Configuration = []map[string]interface{}{
			{
				"key":   "username",
				"value": "admin",
			},
			{
				"key":             "password",
				"encrypted_value": "1f3rrs9uhn63hd",
			},
			{
				"key":       "url",
				"value":     "https://github.com/sample/example.git",
				"is_secure": true,
			},
		}

		server := configRepoServer(configRepo)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.CreateConfigRepoInfo(*configRepo)
		assert.NoError(t, err)
	})
}

func configRepoServer(request interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, t *http.Request) {
		var configRepo gocd.ConfigRepo
		requestByte, err := json.Marshal(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err = writer.Write([]byte(fmt.Sprintf("%s %s", string(requestByte), err.Error()))); err != nil {
				log.Fatalln(err)
			}

			return
		}

		if err = json.Unmarshal(requestByte, &configRepo); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err = writer.Write([]byte(err.Error())); err != nil {
				log.Fatalln(err)
			}
		}

		writer.WriteHeader(http.StatusOK)
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err = writer.Write([]byte("OK")); err != nil {
			log.Fatalln(err)
		}
	}))
}
