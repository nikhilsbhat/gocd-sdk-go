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
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepo()
		assert.EqualError(t, err, "call made to get config repo errored with "+
			"Get \"http://localhost:8156/go/api/admin/config_repos\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repo information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repo information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able retrieve config repo information", func(t *testing.T) {
		server := mockServer([]byte(configReposJSON), http.StatusOK, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(actual))
	})
}

func Test_client_CreateConfigRepoInfo(t *testing.T) {
	t.Run("server should return internal server error as malformed json passed", func(t *testing.T) {
		server := configRepoServer(configReposJSON, http.MethodPost, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		configRepo := &gocd.ConfigRepo{}
		err := client.CreateConfigRepo(*configRepo)
		assert.EqualError(t, err, "body: json: cannot unmarshal string into Go value of type gocd.ConfigRepo httpcode: 500")
	})

	t.Run("should error out while making client call to create config repo", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		configRepo := &gocd.ConfigRepo{}
		err := client.CreateConfigRepo(*configRepo)
		assert.EqualError(t, err, "post call made to create config repo errored with: "+
			"Post \"http://localhost:8156/go/api/admin/config_repos\": dial tcp [::1]:8156: connect: connection refused")
	})

	t.Run("should create config repo successfully", func(t *testing.T) {
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

		server := configRepoServer(configRepo, http.MethodPost, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.CreateConfigRepo(*configRepo)
		assert.NoError(t, err)
	})
}

func Test_client_DeleteConfigRepo(t *testing.T) {
	repoName := "repo1"
	t.Run("should error out while deleting config repo due to server connectivity issues", func(t *testing.T) {
		//server := configRepoServer(nil, http.MethodDelete)
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteConfigRepo(repoName)
		assert.EqualError(t, err, "post call made to create config repo errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
	})

	t.Run("server should return 404 due to wrong header set while deleting config repo", func(t *testing.T) {
		server := configRepoServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionOne})
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.DeleteConfigRepo(repoName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should be able to delete config repo successfully", func(t *testing.T) {
		server := configRepoServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionFour})
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.DeleteConfigRepo(repoName)
		assert.NoError(t, err)
	})
}

func configRepoServer(request interface{}, method string, header map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		for key, value := range header {
			if r.Header.Get(key) != value {
				writer.WriteHeader(http.StatusNotFound)
				if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
					log.Fatalln(err)
				}
				return
			}
		}

		if method == http.MethodDelete {
			writer.WriteHeader(http.StatusOK)
			if _, err := writer.Write([]byte(`{"message": "The config repo 'repo-1' was deleted successfully."}`)); err != nil {
				log.Fatalln(err)
			}
			return
		}

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
	}))
}
