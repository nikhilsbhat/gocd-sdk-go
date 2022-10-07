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

var (
	//go:embed internal/fixtures/config_repos.json
	configReposJSON string
	eTag            = "05548388f7ef5042cd39f7fe42e85735"
	correctHeader   = map[string]string{"Accept": gocd.HeaderVersionFour}
	configRepo      = testGetConfigRepoObj()
)

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should error out while fetching config repos information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepos()
		assert.EqualError(t, err, "call made to get config repos errored with "+
			"Get \"http://localhost:8156/go/api/admin/config_repos\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repos information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepos()
		assert.EqualError(t, err, gocd.APIErrorWithBody("backupJSON", http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repos information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepos()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able retrieve config repos information", func(t *testing.T) {
		server := mockServer([]byte(configReposJSON), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepos()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(actual))
	})
}

func Test_client_CreateConfigRepoInfo(t *testing.T) {
	t.Run("server should return internal server error as malformed json passed while creating config repo", func(t *testing.T) {
		server := mockConfigRepoServer(configReposJSON, http.MethodPost, map[string]string{"Accept": gocd.HeaderVersionFour, "Content-Type": gocd.ContentJSON}, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.CreateConfigRepo(gocd.ConfigRepo{})
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

		err := client.CreateConfigRepo(gocd.ConfigRepo{})
		assert.EqualError(t, err, "call made to create config repo errored with: "+
			"Post \"http://localhost:8156/go/api/admin/config_repos\": dial tcp [::1]:8156: connect: connection refused")
	})

	t.Run("should create config repo successfully", func(t *testing.T) {
		server := mockConfigRepoServer(configRepo, http.MethodPost, map[string]string{"Accept": gocd.HeaderVersionFour, "Content-Type": gocd.ContentJSON}, false)
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
		assert.EqualError(t, err, "call made to create config repo errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
	})

	t.Run("server should return 404 due to wrong header set while deleting config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionOne}, false)
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
		server := mockConfigRepoServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionFour}, false)
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

func Test_client_GetConfigRepo(t *testing.T) {
	repoName := "repo1"
	t.Run("should error out while fetching config repo information as server returned non 200 status code", func(t *testing.T) {
		server := mockConfigRepoServer(configReposJSON, http.MethodPost, correctHeader, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "body: json: cannot unmarshal string into Go value of type gocd.ConfigRepo httpcode: 500")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("should error out while fetching config repo information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("should error out while fetching config repo information from server since server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "call made to get config repo errored with Get "+
			"\"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("server should return 404 due to wrong header set while fetching config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodGet, map[string]string{"Accept": gocd.HeaderVersionOne}, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("server should return 404 no header set while fetching config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodGet, nil, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("should error when header ETag is not set by server", func(t *testing.T) {
		server := mockConfigRepoServer(configRepo, http.MethodGet, correctHeader, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)
		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "header ETag not set, this will impact while updating configrepo")
		assert.Equal(t, configRepo.ID, actual.ID)
		assert.Equal(t, *configRepo, actual)
	})

	t.Run("should be able to get config repo successfully", func(t *testing.T) {
		newConfigRepo := configRepo
		server := mockConfigRepoServer(newConfigRepo, http.MethodGet, correctHeader, true)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		newConfigRepo.ETAG = eTag
		actual, err := client.GetConfigRepo(repoName)
		assert.Nil(t, err)
		assert.Equal(t, newConfigRepo.ID, actual.ID)
		assert.Equal(t, *newConfigRepo, actual)
	})
}

func Test_client_UpdateConfigRepo(t *testing.T) {
	t.Run("should error out while updating config repo information as server returned non 200 status code", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		newCorrectHeader := correctHeader
		newCorrectHeader["If-Match"] = eTag
		server := mockConfigRepoServer(configReposJSON, http.MethodPut, newCorrectHeader, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)
		actual, err := client.UpdateConfigRepo(*configRepo, eTag)

		assert.EqualError(t, err, "body: json: cannot unmarshal string into Go value of type gocd.ConfigRepo httpcode: 500")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while updating config repo since server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.UpdateConfigRepo(*configRepo, eTag)
		assert.EqualError(t, err, "put call made to update config repo errored with: Put "+
			"\"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})

	t.Run("server should return 404 due to wrong header set while updating config repo", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		server := mockConfigRepoServer(newConfigRepo, http.MethodPut, map[string]string{"If-Match": eTag, "Accept": gocd.HeaderVersionOne}, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.UpdateConfigRepo(*newConfigRepo, eTag)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, "", actual)
	})

	t.Run("server should return 404 no header set while updating config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodPut, nil, false)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.UpdateConfigRepo(*configRepo, eTag)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, "", actual)
	})

	t.Run("should error since wrong ETag was specified while updating config repo", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		newCorrectHeader := correctHeader
		newCorrectHeader["If-Match"] = "eTag"
		server := mockConfigRepoServer(newConfigRepo, http.MethodPut, newCorrectHeader, true)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)
		actual, err := client.UpdateConfigRepo(*newConfigRepo, eTag)
		assert.EqualError(t, err, "body: lost update httpcode: 406")
		assert.Equal(t, "", actual)
	})

	t.Run("should update config repo successfully", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		correctHeader["If-Match"] = eTag
		server := mockConfigRepoServer(newConfigRepo, http.MethodPut, correctHeader, true)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)
		actual, err := client.UpdateConfigRepo(*newConfigRepo, eTag)
		assert.NoError(t, err)
		assert.Equal(t, eTag, actual)
	})
}

func Test_client_ConfigRepoTriggerUpdate(t *testing.T) {
	correctConfigHeader := map[string]string{"Accept": gocd.HeaderVersionFour, "X-GoCD-Confirm": "true"}

	t.Run("Should be able to trigger update for a config repo successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"message": "OK"}`), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		expected := map[string]string{
			"message": "OK",
		}

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Should error out while triggering config repo update as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"message": }`), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' looking for beginning of value")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("Should not update config repo as scheduled update is still in progress", func(t *testing.T) {
		server := mockServer([]byte(`{"message": "Update already in progress."}`), http.StatusConflict, correctConfigHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		expected := map[string]string{
			"message": "Update already in progress.",
		}

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Should error out while triggering update config due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusConflict,
			map[string]string{"Accept": gocd.HeaderVersionThree, "X-GoCD-Confirm": "true"}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("Should error out while triggering update config due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("Should error out while triggering update config as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "call made to trigger update configrepo 'config_repo_1' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/config_repos/config_repo_1/trigger_update\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

func Test_client_ConfigRepoStatus(t *testing.T) {
	correctConfigHeader := map[string]string{"Accept": gocd.HeaderVersionFour}

	t.Run("Should be able to fetch the status of the config-repo successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"in_progress": true}`), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		expected := map[string]bool{
			"in_progress": true,
		}

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Should error out while fetching config-repo status as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"in_progress": }`), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' looking for beginning of value")
		assert.Equal(t, map[string]bool(nil), actual)
	})

	t.Run("Should error out while fetching config-repo status due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusConflict,
			map[string]string{"Accept": gocd.HeaderVersionThree, "X-GoCD-Confirm": "true"}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, map[string]bool(nil), actual)
	})

	t.Run("Should error out while fetching config-repo status due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Equal(t, map[string]bool(nil), actual)
	})

	t.Run("Should error out while fetching config-repo status as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "call made to get status of configrepo 'config_repo_1' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/config_repos/config_repo_1/trigger_update\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

func mockConfigRepoServer(request interface{}, method string, header map[string]string, etag bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if header == nil {
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

			return
		}

		for key, value := range header {
			if len(req.Header.Get("If-Match")) != 0 {
				if header["If-Match"] != request.(*gocd.ConfigRepo).ETAG {
					writer.WriteHeader(http.StatusNotAcceptable)
					if _, err := writer.Write([]byte(`lost update`)); err != nil {
						log.Fatalln(err)
					}

					return
				}
			}
			if req.Header.Get(key) != value {
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

		if etag {
			writer.Header().Set("ETag", eTag)
		}

		writer.WriteHeader(http.StatusOK)

		if method == http.MethodGet {
			if _, err = writer.Write(requestByte); err != nil {
				log.Fatalln(err)
			}
		}
	}))
}

func testGetConfigRepoObj() *gocd.ConfigRepo {
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

	return configRepo
}
