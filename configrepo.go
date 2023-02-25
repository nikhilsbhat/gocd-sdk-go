package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

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
		return ConfigRepo{}, &errors.APIError{Err: err, Message: "get config repo"}
	}
	if resp.StatusCode() != http.StatusOK {
		return ConfigRepo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &repoConf); err != nil {
		return ConfigRepo{}, &errors.MarshalError{Err: err}
	}

	if len(resp.Header().Get("ETag")) == 0 {
		return repoConf, &errors.NilHeaderError{Header: "ETag", Message: "updating configrepo"}
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
		return nil, &errors.APIError{Err: err, Message: "get config repos"}
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &reposConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return reposConf.ConfigRepos.ConfigRepos, nil
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
