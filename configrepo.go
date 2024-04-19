package gocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

type PipelineFiles struct {
	Name string
	Path string
}

// GetConfigRepo fetches information of a specific config-repo from GoCD server.
func (conf *client) GetConfigRepo(repo string) (ConfigRepo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ConfigRepo{}, err
	}

	var repoConf ConfigRepo

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionFour,
		}).
		Get(filepath.Join(ConfigReposEndpoint, repo))
	if err != nil {
		return ConfigRepo{}, &errors.APIError{Err: err, Message: "get config-repo"}
	}

	if resp.StatusCode() != http.StatusOK {
		return ConfigRepo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &repoConf); err != nil {
		return ConfigRepo{}, &errors.MarshalError{Err: err}
	}

	if len(resp.Header().Get("ETag")) == 0 {
		return repoConf, &errors.NilHeaderError{Header: "ETag", Message: "getting config-repo"}
	}

	repoConf.ETAG = resp.Header().Get("ETag")

	return repoConf, nil
}

// GetConfigRepos fetches information of all config-repos from GoCD server.
func (conf *client) GetConfigRepos() ([]ConfigRepo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var reposConf ConfigRepoConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionFour,
		}).
		Get(ConfigReposEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get config-repos"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &reposConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return reposConf.ConfigRepos.ConfigRepos, nil
}

// GetConfigReposInternal fetches information about all config repos from the GoCD server using GoCD's internal API.
// Use GetConfigRepos for fetching all config-repos information; use this only if you know why it is being used.
func (conf *client) GetConfigReposInternal() ([]ConfigRepo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var reposConf ConfigRepoConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionFour,
		}).
		Get(ConfigReposInternalEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get config-repos using internal API"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &reposConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return reposConf.ConfigRepos.ConfigRepos, nil
}

// GetConfigRepoDefinitions fetches information of a specific config-repo from GoCD server.
func (conf *client) GetConfigRepoDefinitions(repo string) (ConfigRepo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ConfigRepo{}, err
	}

	var repoConf ConfigRepo

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionFour,
		}).
		Get(filepath.Join(ConfigReposEndpoint, repo, "definitions"))
	if err != nil {
		return ConfigRepo{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get config-repo definitions for '%s'", repo)}
	}

	if resp.StatusCode() != http.StatusOK {
		return ConfigRepo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &repoConf); err != nil {
		return ConfigRepo{}, &errors.MarshalError{Err: err}
	}

	return repoConf, nil
}

// CreateConfigRepo fetches information of all config-repos in GoCD server.
func (conf *client) CreateConfigRepo(repoObj ConfigRepo) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionFour,
			"Content-Type": ContentJSON,
		}).
		SetBody(repoObj).
		Post(ConfigReposEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "create config repo"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// UpdateConfigRepo updates the config repo configurations with the latest configurations provided.
func (conf *client) UpdateConfigRepo(repo ConfigRepo) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionFour,
			"Content-Type": ContentJSON,
			"If-Match":     repo.ETAG,
		}).
		SetBody(repo).
		Put(filepath.Join(ConfigReposEndpoint, repo.ID))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: "call made to update config repo"}
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return resp.Header().Get("ETag"), nil
}

// DeleteConfigRepo deletes a specific config repo.
func (conf *client) DeleteConfigRepo(repo string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionFour,
			"Content-Type": ContentJSON,
		}).
		Delete(filepath.Join(ConfigReposEndpoint, repo))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete config repo '%s'", repo)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// ConfigRepoTriggerUpdate triggers config repo update for a specific config-repo.
func (conf *client) ConfigRepoTriggerUpdate(name string) (map[string]string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionFour,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(ConfigReposEndpoint, name, "trigger_update"))
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("trigger update configrepo '%s'", name)}
	}

	if (resp.StatusCode() != http.StatusCreated) && (resp.StatusCode() != http.StatusConflict) {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var response map[string]string
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return response, &errors.MarshalError{Err: err}
	}

	return response, nil
}

// ConfigRepoStatus fetches the latest available status of the specified config repo.
func (conf *client) ConfigRepoStatus(repo string) (map[string]bool, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionFour,
		}).
		Get(filepath.Join(ConfigReposEndpoint, repo, "status"))
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("get status of configrepo '%s'", repo)}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var response map[string]bool

	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return response, &errors.MarshalError{Err: err}
	}

	return response, nil
}

// ConfigRepoPreflightCheck runs the pre-flight checks on the config-repo with the provided pipeline files.
// Checks posted definition file(s) for syntax and merge errors without updating the current GoCD configuration.
func (conf *client) ConfigRepoPreflightCheck(pipelines map[string]string, pluginID string, repoID string) (bool, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return false, err
	}

	request := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		SetQueryParam("pluginId", pluginID).
		SetQueryParam("repoId", repoID)

	for name, path := range pipelines {
		pipelineBytes, err := os.ReadFile(path)
		if err != nil {
			return false, err
		}

		request.SetFileReader("files[]", name, bytes.NewReader(pipelineBytes))
	}

	resp, err := request.Post(PreflightCheckEndpoint)
	if err != nil {
		return false, &errors.APIError{Err: err, Message: fmt.Sprintf("preflight check of confirepo '%s'", repoID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return false, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var response map[string]interface{}

	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return false, &errors.MarshalError{Err: err}
	}

	if value, ok := response["errors"]; ok {
		errorSlice := GetSLice(value)
		if len(errorSlice) != 0 {
			return false, &errors.GoCDSDKError{Message: strings.Join(errorSlice, "\n")}
		}
	}

	return response["valid"].(bool), nil
}

// SetPipelineFiles transforms an array of pipeline files([]PipelineFiles) to map which could be utilised by ConfigRepoPreflightCheck.
func (conf *client) SetPipelineFiles(pipelines []PipelineFiles) map[string]string {
	fileMap := make(map[string]string)
	for _, pipeline := range pipelines {
		fileMap[pipeline.Name] = pipeline.Path
	}

	return fileMap
}

// GetPipelineFiles reads the pipeline file or recursively read the directory to get all the pipelines matching the pattern and transforms to []PipelineFiles
// So that SetPipelineFiles can convert them to format that ConfigRepoPreflightCheck understands.
func (conf *client) GetPipelineFiles(pathAndPattern ...string) ([]PipelineFiles, error) {
	path := pathAndPattern[0]
	patterns := pathAndPattern[1:]

	var pipelineFiles []PipelineFiles

	fileInfo, err := os.Stat(pathAndPattern[0])
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		conf.logger.Debugf("pipeline files path '%s' is a file finding absolute path of the same", path)

		absFilePath, err := filepath.Abs(path)
		if err != nil {
			return pipelineFiles, err
		}

		_, fileName := filepath.Split(absFilePath)
		pipelineFiles = append(pipelineFiles, PipelineFiles{
			Name: fileName,
			Path: absFilePath,
		})

		return pipelineFiles, nil
	}

	if len(pathAndPattern) <= 1 {
		return nil, &errors.GoCDSDKError{Message: "pipeline files pattern not passed (ex: *.gocd.yaml)"}
	}

	conf.logger.Debugf("pipeline files path '%s' is a directory, finding all the files matching the pattern '%s'", path, strings.Join(patterns, ","))

	if err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, pattern := range patterns {
			match, err := filepath.Match(pattern, info.Name())
			if err != nil {
				conf.logger.Errorf("matching GoCD pipeline file errored with '%s'", err)
			}

			if match {
				conf.logger.Debugf("identified pipeline '%s' under path '%s'", info.Name(), filepath.Dir(path))

				absPath, err := filepath.Abs(path)
				if err != nil {
					conf.logger.Errorf("finding absolute path of pipeline '%s' errored with '%s'", info.Name(), err)
				} else {
					path = absPath
				}

				pipelineFiles = append(pipelineFiles, PipelineFiles{
					Name: info.Name(),
					Path: path,
				})
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return pipelineFiles, nil
}
