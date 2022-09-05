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
		"Accept": GoCdHeaderVersionFour,
	})

	var reposConf ConfigRepoConfig
	resp, err := newClient.httpClient.R().Get(GoCdConfigReposEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get config repo errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, ApiWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &reposConf); err != nil {
		return nil, ResponseReadError(err.Error())
	}

	return reposConf.ConfigRepos.ConfigRepos, nil
}
