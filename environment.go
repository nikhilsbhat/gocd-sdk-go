package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetEnvironmentInfo fetches information of backup configured in GoCD server.
func (conf *client) GetEnvironmentInfo() ([]Environment, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionThree,
	})

	var envConf EnvironmentInfo
	resp, err := newClient.httpClient.R().Get(GoCdEnvironmentEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get environment errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, ApiWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &envConf); err != nil {
		return nil, ResponseReadError(err.Error())
	}

	return envConf.Environments.Environments, nil
}
