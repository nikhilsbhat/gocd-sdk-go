package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// GetAuthConfigs fetches all authorization configurations present iin GoCD.
func (conf *client) GetAuthConfigs() ([]CommonConfig, error) {
	var auth AuthConfigs

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
		return nil, &errors.APIError{Err: err, Message: "get auth configs"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return nil, err
	}

	return auth.Config.AuthConfigs, nil
}

// GetAuthConfig fetches authorization configuration for specified id.
func (conf *client) GetAuthConfig(name string) (CommonConfig, error) {
	var auth CommonConfig

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
		return auth, &errors.APIError{Err: err, Message: fmt.Sprintf("get auth config '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return auth, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return auth, err
	}

	auth.ETAG = resp.Header().Get("ETag")

	return auth, nil
}

// CreateAuthConfig creates an authorization configuration with the provided configurations.
func (conf *client) CreateAuthConfig(config CommonConfig) (CommonConfig, error) {
	var auth CommonConfig

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
		return auth, &errors.APIError{Err: err, Message: fmt.Sprintf("create auth config '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return auth, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &auth); err != nil {
		return auth, err
	}

	auth.ETAG = resp.Header().Get("ETag")

	return auth, nil
}

// UpdateAuthConfig updates some attributes of an authorization configuration.
func (conf *client) UpdateAuthConfig(config CommonConfig) (CommonConfig, error) {
	var auth CommonConfig

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
		return auth, &errors.APIError{Err: err, Message: fmt.Sprintf("update auth config '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return auth, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete auth config '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
