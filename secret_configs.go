package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetSecretConfigs() (SecretsConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return SecretsConfig{}, err
	}

	var secretsCfg SecretsConfigs

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(SecretsConfigEndpoint)
	if err != nil {
		return SecretsConfig{}, &errors.APIError{Err: err, Message: "get secret configs"}
	}

	if resp.StatusCode() != http.StatusOK {
		return SecretsConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &secretsCfg); err != nil {
		return SecretsConfig{}, &errors.MarshalError{Err: err}
	}

	secretsCfg.SecretsConfigs.ETAG = resp.Header().Get("ETag")

	return secretsCfg.SecretsConfigs, nil
}

func (conf *client) GetSecretConfig(name string) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var secretCfg CommonConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(filepath.Join(SecretsConfigEndpoint, name))
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get secret config '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &secretCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
	}

	secretCfg.ETAG = resp.Header().Get("ETag")

	return secretCfg, nil
}

func (conf *client) CreateSecretConfig(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var secretsCfg CommonConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(SecretsConfigEndpoint)
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("create secrets config '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &secretsCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
	}

	secretsCfg.ETAG = resp.Header().Get("ETag")

	return secretsCfg, nil
}

func (conf *client) UpdateSecretConfig(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var secretsCfg CommonConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(filepath.Join(SecretsConfigEndpoint, config.ID))
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("update secret config '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &secretsCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
	}

	secretsCfg.ETAG = resp.Header().Get("ETag")

	return secretsCfg, nil
}

func (conf *client) DeleteSecretConfig(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Delete(filepath.Join(SecretsConfigEndpoint, name))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete secret config '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
