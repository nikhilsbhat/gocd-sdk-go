package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
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
		return nil, &errors.APIError{Err: err, Message: "get environments"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &envConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return envConf.Environments.Environments, nil
}

// GetEnvironment fetches information of a specific environment from GoCD.
func (conf *client) GetEnvironment(name string) (Environment, error) {
	var env Environment

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return env, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(filepath.Join(EnvironmentEndpoint, name))
	if err != nil {
		return env, &errors.APIError{Err: err, Message: fmt.Sprintf("get environment '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return env, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &env); err != nil {
		return env, &errors.MarshalError{Err: err}
	}

	env.ETAG = resp.Header().Get("ETag")

	return env, nil
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("create environment '%s'", environment.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// PatchEnvironment Update some attributes of an environment.
func (conf *client) PatchEnvironment(environment any) (Environment, error) {
	envPatch := environment.(PatchEnvironment)

	var env Environment

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return env, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(envPatch).
		Patch(filepath.Join(EnvironmentEndpoint, envPatch.Name))
	if err != nil {
		return env, &errors.APIError{Err: err, Message: fmt.Sprintf("patch environment '%s'", envPatch.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return env, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err := json.Unmarshal(resp.Body(), &env); err != nil {
		return env, &errors.MarshalError{Err: err}
	}

	return env, nil
}

// UpdateEnvironment will update the environment configurations of a already created GoCD environment.
func (conf *client) UpdateEnvironment(environment Environment) (Environment, error) {
	var env Environment

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return env, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
			"If-Match":     environment.ETAG,
		}).
		SetBody(environment).
		Put(filepath.Join(EnvironmentEndpoint, environment.Name))
	if err != nil {
		return env, &errors.APIError{Err: err, Message: fmt.Sprintf("update environment '%s'", environment.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return env, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &env); err != nil {
		return env, &errors.MarshalError{Err: err}
	}

	env.ETAG = resp.Header().Get("ETag")

	return env, nil
}

// DeleteEnvironment deletes the specified GoCD environment.
func (conf *client) DeleteEnvironment(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).Delete(filepath.Join(EnvironmentEndpoint, name))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete environment '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) GetEnvironmentsMerged(names []string) ([]Environment, error) {
	var env EnvironmentInfo

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).Get(EnvironmentInternalEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("get environment mapping of '%s'", strings.Join(names, ", "))}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &env); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	environments := make([]Environment, 0)

	for _, environment := range env.Environments.Environments {
		for _, name := range names {
			if environment.Name == name {
				environments = append(environments, environment)
			}
		}
	}

	if len(environments) != 0 {
		return environments, nil
	}

	return nil, &errors.GoCDSDKError{Message: fmt.Sprintf("no environments found with names '%s' to get mappings", strings.Join(names, ", "))}
}
