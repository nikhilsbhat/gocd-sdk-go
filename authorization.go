package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

// GetAuthConfigs fetches all authorization configurations present iin GoCD.
func (conf *client) GetAuthConfigs() ([]AuthConfig, error) {
	var auth AuthConfigs
	{
	}
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(AuthConfigEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get auth configs errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return nil, err
	}

	return auth.Config.AuthConfigs, nil
}

// GetAuthConfig fetches authorization configuration for specified id.
func (conf *client) GetAuthConfig(name string) (AuthConfig, error) {
	var auth AuthConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return auth, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(filepath.Join(AuthConfigEndpoint, name))
	if err != nil {
		return auth, fmt.Errorf("call made to get auth config '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return auth, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return auth, err
	}

	auth.ETAG = resp.Header().Get("ETag")

	return auth, nil
}

// CreateAuthConfig creates an authorization configuration with the provided configurations.
func (conf *client) CreateAuthConfig(config AuthConfig) (AuthConfig, error) {
	var auth AuthConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return auth, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(AuthConfigEndpoint)
	if err != nil {
		return auth, fmt.Errorf("call made to create auth config '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return auth, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return auth, err
	}

	auth.ETAG = resp.Header().Get("ETag")

	return auth, nil
}

// UpdateAuthConfig updates some attributes of an authorization configuration.
func (conf *client) UpdateAuthConfig(config AuthConfig) (AuthConfig, error) {
	var auth AuthConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return auth, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(AuthConfigEndpoint)
	if err != nil {
		return auth, fmt.Errorf("call made to update auth config '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return auth, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return auth, err
	}

	auth.ETAG = resp.Header().Get("ETag")

	return auth, nil
}

// DeleteAuthConfig deletes the specified authorization configuration.
func (conf *client) DeleteAuthConfig(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Delete(filepath.Join(AuthConfigEndpoint, name))
	if err != nil {
		return fmt.Errorf("call made to delete auth config '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
