package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
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
		return SecretsConfig{}, fmt.Errorf("call made to get secret configs errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return SecretsConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &secretsCfg); err != nil {
		return SecretsConfig{}, ResponseReadError(err.Error())
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
		return CommonConfig{}, fmt.Errorf("call made to get secret config '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &secretCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
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
		return CommonConfig{}, fmt.Errorf("call made to create secrets config '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &secretsCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
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
		return CommonConfig{}, fmt.Errorf("call made to update secret config '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &secretsCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
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
		return fmt.Errorf("call made to delete secret config '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
