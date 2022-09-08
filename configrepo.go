package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

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
		return ConfigRepo{}, fmt.Errorf("call made to get config repo errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return ConfigRepo{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &repoConf); err != nil {
		return ConfigRepo{}, ResponseReadError(err.Error())
	}

	if len(resp.Header().Get("ETag")) == 0 {
		return repoConf, fmt.Errorf("header ETag not set, this will impact while updating configrepo") //nolint:goerr113
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
		return nil, fmt.Errorf("call made to get config repos errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &reposConf); err != nil {
		return nil, ResponseReadError(err.Error())
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
		return fmt.Errorf("post call made to create config repo errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// UpdateConfigRepo updates the config repo configurations with the latest configurations provided.
func (conf *client) UpdateConfigRepo(repoObj ConfigRepo, etag string) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionFour,
			"Content-Type": ContentJSON,
			"If-Match":     etag,
		}).
		SetBody(repoObj).
		Put(filepath.Join(ConfigReposEndpoint, repoObj.ID))
	if err != nil {
		return "", fmt.Errorf("put call made to update config repo errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", APIErrorWithBody(resp.String(), resp.StatusCode())
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
		return fmt.Errorf("post call made to create config repo errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
