package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetEnvironments fetches information of backup configured in GoCD server.
func (conf *client) GetEnvironments() ([]Environment, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var envConf EnvironmentInfo
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(EnvironmentEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get environment errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &envConf); err != nil {
		return nil, ResponseReadError(err.Error())
	}

	return envConf.Environments.Environments, nil
}

// CreateEnvironment creates GoCD environment with the specified configurations.
func (conf *client) CreateEnvironment(environment Environment) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(environment).
		Post(EnvironmentEndpoint)
	if err != nil {
		return fmt.Errorf("call made to create environment errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
