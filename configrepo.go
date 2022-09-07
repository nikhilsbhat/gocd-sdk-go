package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

// GetConfigRepo fetches information of all config-repos in GoCD server.
func (conf *client) GetConfigRepo() ([]ConfigRepo, error) {
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
		return nil, fmt.Errorf("call made to get config repo errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIWithCodeError(resp.StatusCode())
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
			"Content-Type": contentJSON,
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

// DeleteConfigRepo deletes a specific config repo.
func (conf *client) DeleteConfigRepo(repo string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionFour,
			"Content-Type": contentJSON,
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
