package gocd

import (
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
	resp, err := newClient.httpClient.R().SetResult(&reposConf).Get(GoCdConfigReposEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get config repo errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}

	return reposConf.ConfigRepos.ConfigRepos, nil
}
