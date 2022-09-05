package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetConfigRepoInfo fetches information of all config-repos in GoCD server.
func (conf *client) GetConfigRepoInfo() ([]ConfigRepo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": HeaderVersionFour,
	})

	var reposConf ConfigRepoConfig
	resp, err := newClient.httpClient.R().Get(ConfigReposEndpoint)
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

// CreateConfigRepoInfo fetches information of all config-repos in GoCD server.
func (conf *client) CreateConfigRepoInfo(repoObj ConfigRepo) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept":       HeaderVersionFour,
		"Content-Type": contentJSON,
	})

	resp, err := newClient.httpClient.R().SetBody(repoObj).Post(ConfigReposEndpoint)
	if err != nil {
		return fmt.Errorf("post call made to create config repo errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
