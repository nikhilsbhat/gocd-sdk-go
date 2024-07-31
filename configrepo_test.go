package gocd_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/config_repos.json
	configReposJSON string
	//go:embed internal/fixtures/config_repo.json
	configRepoJSON string
	//go:embed internal/fixtures/config_repo_definitions.json
	configRepoDefinition string
	//go:embed internal/fixtures/config_repo_internal.json
	configRepoInternalJSON string
	eTag                   = "05548388f7ef5042cd39f7fe42e85735"
	correctHeader          = map[string]string{"Accept": gocd.HeaderVersionFour}
	configRepo             = testGetConfigRepoObj()
)

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should error out while fetching config repos information from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepos()
		assert.EqualError(t, err, "call made to get config-repos errored with: "+
			"Get \"http://localhost:8156/go/api/admin/config_repos\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repos information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigRepos()
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos\nwith BODY:backupJSON")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repos information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigRepos()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able retrieve config repos information", func(t *testing.T) {
		server := mockServer([]byte(configReposJSON), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.ConfigRepo{
			{
				ID:       "repo1",
				PluginID: "json.config.plugin",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:               "https://github.com/config-repo/gocd-json-config-example.git",
						Username:          "bob",
						EncryptedPassword: "aSdiFgRRZ6A=",
						Branch:            "master",
						AutoUpdate:        true,
					},
				},
				Configuration: []gocd.PluginConfiguration{
					{
						Key:   "username",
						Value: "admin",
					},
					{
						Key:            "password",
						EncryptedValue: "1f3rrs9uhn63hd",
					},
					{
						Key:      "url",
						Value:    "https://github.com/sample/example.git",
						IsSecure: true,
					},
				},
				Rules: []map[string]string{
					{
						"directive": "allow",
						"action":    "refer",
						"type":      "pipeline_group",
						"resource":  "*",
					},
				},
			},
		}
		actual, err := client.GetConfigRepos()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestConfig_GetConfigReposInternal(t *testing.T) {
	t.Run("should error out while fetching config repos information from server using GoCD's internal API", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigReposInternal()
		assert.EqualError(t, err, "call made to get config-repos using internal API errored with: "+
			"Get \"http://localhost:8156/go/api/internal/config_repos\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repos information using GoCD's internal API, as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigReposInternal()
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/internal/config_repos\nwith BODY:backupJSON")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching config repos information using GoCD's internal API, as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigReposInternal()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able retrieve config repos information using GoCD's internal API", func(t *testing.T) {
		server := mockServer([]byte(configRepoInternalJSON), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.ConfigRepo{
			{
				PluginID: "json.config.plugin",
				ID:       "gocd-go-sdk",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:        "https://github.com/nikhilsbhat/gocd-sdk-go.git",
						Branch:     "master",
						AutoUpdate: true,
					},
				},
				Configuration: []gocd.PluginConfiguration{},
				Rules:         []map[string]string{},
				ConfigRepoParseInfo: gocd.ConfigRepoParseInfo{
					LatestParsedModification: map[string]interface{}{
						"comment":       "Add support for GET config-repo definitions API",
						"email_address": interface{}(nil),
						"modified_time": "2023-06-27T13:46:33Z",
						"revision":      "2d1e4525a6f26cf0699c06c2ce36ab6ac512c9e6",
						"username":      "nikhilsbhat <nikhilsbhat93@gmail.com>",
					},
					GoodModification: map[string]interface{}{
						"comment":       "Add support for GET config-repo definitions API",
						"email_address": interface{}(nil), "modified_time": "2023-06-27T13:46:33Z",
						"revision": "2d1e4525a6f26cf0699c06c2ce36ab6ac512c9e6",
						"username": "nikhilsbhat <nikhilsbhat93@gmail.com>",
					},
				},
			},
			{
				PluginID: "yaml.config.plugin",
				ID:       "sample_config_repo",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:               "https://github.com/config-repo/gocd-json-config-example.git",
						Username:          "bob",
						EncryptedPassword: "AES:I/umvAruOKkDyHJFflavCQ==:4hikK7OSpJN50E4SerstZw==",
						Branch:            "master",
						AutoUpdate:        true,
					},
				},
				Configuration: []gocd.PluginConfiguration{
					{
						Key:   "url",
						Value: "https://github.com/config-repo/gocd-json-config-example.git",
					},
					{
						Key:   "username",
						Value: "admin",
					},
					{
						Key:   "password",
						Value: "admin",
					},
				},
				Rules: []map[string]string{
					{
						"action":    "refer",
						"directive": "allow",
						"resource":  "*",
						"type":      "pipeline_group",
					},
				},
				ConfigRepoParseInfo: gocd.ConfigRepoParseInfo{
					Error: "MODIFICATION CHECK FAILED FOR MATERIAL: " +
						"URL: HTTPS://GITHUB.COM/CONFIG-REPO/GOCD-JSON-CONFIG-EXAMPLE.GIT, " +
						"BRANCH: MASTER\nNO PIPELINES ARE AFFECTED BY THIS MATERIAL, " +
						"PERHAPS THIS MATERIAL IS UNUSED.\nFailed to run git clone command " +
						"STDERR: Cloning into '/Users/nikhil.bhat/idfc/gocd-setup/go-server-22.1.0/pipelines/flyweight/2b3feb60-efd7-41d3-8041-3e0d3208285e'...\n" +
						"STDERR: remote: Support for password authentication was removed on August 13, 2021.\n" +
						"STDERR: remote: Please see https://docs.github.com/en/get-started/getting-started-with-git/about-remote-repositories#cloning-with-https-urls " +
						"for information on currently recommended modes of authentication.\n" +
						"STDERR: fatal: Authentication failed for 'https://github.com/config-repo/gocd-json-config-example.git/'",
				},
			},
		}
		actual, err := client.GetConfigReposInternal()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateConfigRepoInfo(t *testing.T) {
	t.Run("server should return internal server error as malformed json passed while creating config repo", func(t *testing.T) {
		server := mockConfigRepoServer(configReposJSON, http.MethodPost, map[string]string{"Accept": gocd.HeaderVersionFour, "Content-Type": gocd.ContentJSON}, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.CreateConfigRepo(gocd.ConfigRepo{})
		assert.EqualError(t, err, "got 500 from GoCD while making POST call for "+server.URL+
			"/api/admin/config_repos\nwith BODY:json: cannot unmarshal string into Go value of type gocd.ConfigRepo")
	})

	t.Run("should error out while making client call to create config repo", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.CreateConfigRepo(gocd.ConfigRepo{})
		assert.EqualError(t, err, "call made to create config repo errored with: "+
			"Post \"http://localhost:8156/go/api/admin/config_repos\": dial tcp [::1]:8156: connect: connection refused")
	})

	t.Run("should create config repo successfully", func(t *testing.T) {
		server := mockConfigRepoServer(configRepo, http.MethodPost, map[string]string{"Accept": gocd.HeaderVersionFour, "Content-Type": gocd.ContentJSON}, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.CreateConfigRepo(*configRepo)
		assert.NoError(t, err)
	})
}

func Test_client_DeleteConfigRepo(t *testing.T) {
	repoName := "repo1"

	t.Run("should error out while deleting config repo due to server connectivity issues", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteConfigRepo(repoName)
		assert.EqualError(t, err, "call made to delete config repo 'repo1' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
	})

	t.Run("server should return 404 due to wrong header set while deleting config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionOne}, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteConfigRepo(repoName)
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should be able to delete config repo successfully", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionFour}, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteConfigRepo(repoName)
		assert.NoError(t, err)
	})
}

func Test_client_GetConfigRepo(t *testing.T) {
	repoName := "repo1"

	t.Run("should error out while fetching config repo information as server returned non 200 status code", func(t *testing.T) {
		server := mockConfigRepoServer(configRepoJSON, http.MethodPost, correctHeader, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "got 500 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:json: cannot unmarshal string into Go value of type gocd.ConfigRepo")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("should error out while fetching config repo information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("should error out while fetching config repo information from server since server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "call made to get config-repo errored with: "+
			"Get \"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("server should return 404 due to wrong header set while fetching config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodGet, map[string]string{"Accept": gocd.HeaderVersionOne}, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("server should return 404 no header set while fetching config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodGet, nil, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, gocd.ConfigRepo{}, actual)
	})

	t.Run("should error when header ETag is not set by server", func(t *testing.T) {
		server := mockConfigRepoServer(configRepo, http.MethodGet, correctHeader, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		actual, err := client.GetConfigRepo(repoName)
		assert.EqualError(t, err, "header ETag not set, this will impact while getting config-repo")
		assert.Equal(t, configRepo.ID, actual.ID)
		assert.Equal(t, *configRepo, actual)
	})

	t.Run("should be able to get config repo successfully", func(t *testing.T) {
		newConfigRepo := configRepo
		server := mockConfigRepoServer(newConfigRepo, http.MethodGet, correctHeader, true)
		client := gocd.NewClient(server.URL, auth, "info", nil)

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
		client := gocd.NewClient(server.URL, auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)
		actual, err := client.UpdateConfigRepo(*configRepo)

		assert.EqualError(t, err, "got 500 from GoCD while making PUT call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:json: cannot unmarshal string into Go value of type gocd.ConfigRepo")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while updating config repo since server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.UpdateConfigRepo(*configRepo)
		assert.EqualError(t, err, "call made to call made to update config repo errored with: "+
			"Put \"http://localhost:8156/go/api/admin/config_repos/repo1\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})

	t.Run("server should return 404 due to wrong header set while updating config repo", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		server := mockConfigRepoServer(newConfigRepo, http.MethodPut, map[string]string{"If-Match": eTag, "Accept": gocd.HeaderVersionOne}, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.UpdateConfigRepo(*newConfigRepo)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("server should return 404 no header set while updating config repo", func(t *testing.T) {
		server := mockConfigRepoServer(nil, http.MethodPut, nil, false)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.UpdateConfigRepo(*configRepo)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error since wrong ETag was specified while updating config repo", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		newCorrectHeader := correctHeader
		newCorrectHeader["If-Match"] = "eTag"
		server := mockConfigRepoServer(newConfigRepo, http.MethodPut, newCorrectHeader, true)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		actual, err := client.UpdateConfigRepo(*newConfigRepo)
		assert.EqualError(t, err, "got 406 from GoCD while making PUT call for "+server.URL+
			"/api/admin/config_repos/repo1\nwith BODY:lost update")
		assert.Equal(t, "", actual)
	})

	t.Run("should update config repo successfully", func(t *testing.T) {
		newConfigRepo := configRepo
		newConfigRepo.ETAG = eTag
		correctHeader["If-Match"] = eTag
		server := mockConfigRepoServer(newConfigRepo, http.MethodPut, correctHeader, true)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		actual, err := client.UpdateConfigRepo(*newConfigRepo)
		assert.NoError(t, err)
		assert.Equal(t, eTag, actual)
	})
}

func Test_client_ConfigRepoTriggerUpdate(t *testing.T) {
	correctConfigHeader := map[string]string{"Accept": gocd.HeaderVersionFour, "X-GoCD-Confirm": "true"}

	t.Run("Should be able to trigger update for a config repo successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"message": "OK"}`), http.StatusCreated, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := map[string]string{
			"message": "OK",
		}

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Should error out while triggering config repo update as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"message": }`), http.StatusCreated, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' looking for beginning of value")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("Should not update config repo as scheduled update is still in progress", func(t *testing.T) {
		server := mockServer([]byte(`{"message": "Update already in progress."}`), http.StatusConflict, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

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
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/config_repos/config_repo_1/trigger_update\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("Should error out while triggering update config due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ConfigRepoTriggerUpdate("config_repo_1")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/config_repos/config_repo_1/trigger_update\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("Should error out while triggering update config as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

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
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := map[string]bool{
			"in_progress": true,
		}

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Should error out while fetching config-repo status as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"in_progress": }`), http.StatusOK, correctConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' looking for beginning of value")
		assert.Equal(t, map[string]bool(nil), actual)
	})

	t.Run("Should error out while fetching config-repo status due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusConflict,
			map[string]string{"Accept": gocd.HeaderVersionThree, "X-GoCD-Confirm": "true"}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/config_repo_1/status\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, map[string]bool(nil), actual)
	})

	t.Run("Should error out while fetching config-repo status due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/config_repo_1/status\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, map[string]bool(nil), actual)
	})

	t.Run("Should error out while fetching config-repo status as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.ConfigRepoStatus("config_repo_1")
		assert.EqualError(t, err, "call made to get status of configrepo 'config_repo_1' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/config_repos/config_repo_1/status\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

func Test_client_ConfigRepoPreflightCheck(t *testing.T) {
	correctPreflightHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	preflightCheckJSON := `{"errors" : [],"valid" : true}`

	t.Run("should be able to run config-repo preflight checks successfully", func(t *testing.T) {
		server := mockServer([]byte(preflightCheckJSON), http.StatusOK,
			correctPreflightHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "debug", nil)

		pipelineFiles, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)

		pipelineMap := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipelineMap, "yaml.config.plugin", "sample")
		assert.NoError(t, err)
		assert.Equal(t, true, actual)
	})

	t.Run("should error out while running config-repo preflight checks in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(preflightCheckJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		pipelineFiles, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)

		pipelineMap := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipelineMap, "yaml.config.plugin", "sample")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/config_repo_ops/preflight?pluginId=yaml.config.plugin&repoId=sample\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, false, actual)
	})

	t.Run("should error out while running config-repo preflight checks in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(preflightCheckJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		pipelineFiles, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)

		pipelineMap := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipelineMap, "yaml.config.plugin", "sample")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/config_repo_ops/preflight?pluginId=yaml.config.plugin&repoId=sample\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, false, actual)
	})

	t.Run("should error out while running config-repo preflight checks in GoCD as preflight checks failed", func(t *testing.T) {
		preflightCheckJSONNew := `{"errors" : ["invalid merge configurations, duplicate key TEST","invalid merge configurations, duplicate key ENV"],"valid" : false}`
		server := mockServer([]byte(preflightCheckJSONNew), http.StatusOK, correctPreflightHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		pipelineFiles, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)

		pipelineMap := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipelineMap, "yaml.config.plugin", "sample")
		assert.EqualError(t, err, "invalid merge configurations, duplicate key TEST\ninvalid merge configurations, duplicate key ENV")
		assert.Equal(t, false, actual)
	})

	t.Run("should error out while running config-repo preflight checks in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("preflightCheckJSON"), http.StatusOK, correctPreflightHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		pipelineFiles, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)

		pipelineMap := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipelineMap, "yaml.config.plugin", "sample")
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, false, actual)
	})

	t.Run("should error out while running config-repo preflight checks in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		pipelineFiles, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)

		pipelineMap := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipelineMap, "yaml.config.plugin", "sample")
		assert.EqualError(t, err, "call made to preflight check of confirepo 'sample' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/config_repo_ops/preflight?pluginId=yaml.config.plugin&repoId=sample\": "+
			"dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, false, actual)
	})

	t.Run("should be able to run config-repo preflight checks successfully", func(t *testing.T) {
		goCDAuth := gocd.Auth{
			UserName: "admin",
			Password: "admin",
		}
		client := gocd.NewClient("http://localhost:8153/go", goCDAuth, "info", nil)

		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)

		pipelineFiles, err := client.GetPipelineFiles(filepath.Join(homeDir, "opensource/gocd-git-path-sample"), nil, "*.gocd.yaml")
		assert.NoError(t, err)

		pipeliness := client.SetPipelineFiles(pipelineFiles)

		actual, err := client.ConfigRepoPreflightCheck(pipeliness, "yaml.config.plugin", "sample-repo")
		assert.NoError(t, err)
		assert.Equal(t, true, actual)
	})
}

func Test_client_SetPipelineFiles(t *testing.T) {
	t.Run("should be able to return the map equivalent of []gocd.PipelineFiles", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		pipelines := []gocd.PipelineFiles{
			{
				Name: "pipeline1",
				Path: "absolute/path/to/pipeline/1",
			},
			{
				Name: "pipeline2",
				Path: "absolute/path/to/pipeline/2",
			},
			{
				Name: "pipeline3",
				Path: "absolute/path/to/pipeline/3",
			},
			{
				Name: "pipeline4",
				Path: "absolute/path/to/pipeline/4",
			},
			{
				Name: "pipeline5",
				Path: "absolute/path/to/pipeline/5",
			},
		}

		expected := map[string]string{
			"pipeline1": "absolute/path/to/pipeline/1",
			"pipeline2": "absolute/path/to/pipeline/2",
			"pipeline3": "absolute/path/to/pipeline/3",
			"pipeline4": "absolute/path/to/pipeline/4",
			"pipeline5": "absolute/path/to/pipeline/5",
		}

		actual := client.SetPipelineFiles(pipelines)

		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetPipelineFiles(t *testing.T) {
	t.Run("should be able to identify the path as directory and fetch the pipelines", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "debug", nil)

		expected := []gocd.PipelineFiles{
			{
				Name: "mail_server_config.json",
				Path: "/Users/nikhil.bhat/my-opensource/gocd-sdk-go/internal/fixtures/mail_server_config.json",
			},
			{
				Name: "role_config.json",
				Path: "/Users/nikhil.bhat/my-opensource/gocd-sdk-go/internal/fixtures/role_config.json",
			},
			{
				Name: "roles_config.json",
				Path: "/Users/nikhil.bhat/my-opensource/gocd-sdk-go/internal/fixtures/roles_config.json",
			},
			{
				Name: "secrets_config.json",
				Path: "/Users/nikhil.bhat/my-opensource/gocd-sdk-go/internal/fixtures/secrets_config.json",
			},
		}

		actual, err := client.GetPipelineFiles("internal/fixtures", nil, "*_config.json")
		assert.NoError(t, err)
		assert.Equal(t, len(expected), len(actual))
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while parsing directory to fetch the pipelines since pattern is missing", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "debug", nil)

		actual, err := client.GetPipelineFiles("internal/fixture", nil, "*_config.json")
		assert.EqualError(t, err, "lstat internal/fixture: no such file or directory")
		assert.Nil(t, actual)
	})

	t.Run("should error out while parsing directory due to wrong path", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "debug", nil)

		actual, err := client.GetPipelineFiles("internal/fixtures", nil)
		assert.EqualError(t, err, "pipeline files pattern not passed (ex: *.gocd.yaml)")
		assert.Nil(t, actual)
	})

	t.Run("should be able to identify the path as file", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "debug", nil)

		expected := []gocd.PipelineFiles{
			{
				Name: "mail_server_config.json",
				Path: "/Users/nikhil.bhat/my-opensource/gocd-sdk-go/internal/fixtures/mail_server_config.json",
			},
		}

		actual, err := client.GetPipelineFiles("", []string{"internal/fixtures/mail_server_config.json"})
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while identifying the pipeline files due to wrong path", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "debug", nil)

		actual, err := client.GetPipelineFiles("", []string{"internal/fixture/mail_server_config.json"})
		assert.EqualError(t, err, "stat internal/fixture/mail_server_config.json: no such file or directory")
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

func Test_client_GetConfigRepoDefinitions(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionFour}
	configRepoName := "config-repo-group"

	t.Run("should be able to fetch the definitions defined in config repo successfully", func(t *testing.T) {
		server := mockServer([]byte(configRepoDefinition), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ConfigRepo{
			Environments: []gocd.Environment{
				{
					Name: "dev",
				},
			},
			Groups: []gocd.PipelineGroup{
				{
					Name: configRepoName,
					Pipelines: []gocd.Pipeline{
						{
							Name: "pipeline1",
						},
						{
							Name: "pipeline2",
						},
					},
				},
			},
		}

		actual, err := client.GetConfigRepoDefinitions(configRepoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching definitions defined in config repo present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoresJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ConfigRepo{}

		actual, err := client.GetConfigRepoDefinitions(configRepoName)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/config-repo-group/definitions\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching definitions defined in config repo present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(artifactStoresJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ConfigRepo{}

		actual, err := client.GetConfigRepoDefinitions(configRepoName)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config_repos/config-repo-group/definitions\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching definitions defined in config repo from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("artifactStoreJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ConfigRepo{}

		actual, err := client.GetConfigRepoDefinitions(configRepoName)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching definitions defined in config repo present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.ConfigRepo{}

		actual, err := client.GetConfigRepoDefinitions(configRepoName)
		assert.EqualError(t, err, "call made to get config-repo definitions for 'config-repo-group' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/config_repos/config-repo-group/definitions\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func testGetConfigRepoObj() *gocd.ConfigRepo {
	configRepo := &gocd.ConfigRepo{
		PluginID:      "json.config.plugin",
		ID:            "repo1",
		Configuration: nil,
		Rules:         nil,
		Material:      gocd.Material{},
	}

	configRepo.Material.Type = "git"
	configRepo.Material.Attributes = gocd.Attribute{}
	configRepo.Material.Attributes.URL = "https://github.com/config-repo/gocd-json-config-example.git"
	configRepo.Material.Attributes.AutoUpdate = false
	configRepo.Material.Attributes.Branch = "master"
	configRepo.Rules = []map[string]string{
		{
			"directive": "allow",
			"action":    "refer",
			"type":      "pipeline_group",
			"resource":  "*",
		},
	}
	configRepo.Configuration = []gocd.PluginConfiguration{
		{
			Key:   "username",
			Value: "admin",
		},
		{
			Key:            "password",
			EncryptedValue: "1f3rrs9uhn63hd",
		},
		{
			Key:      "url",
			Value:    "https://github.com/sample/example.git",
			IsSecure: true,
		},
	}

	return configRepo
}
