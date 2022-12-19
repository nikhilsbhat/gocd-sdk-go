package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

func (conf *client) GetClusterProfiles() (ProfilesConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ProfilesConfig{}, err
	}

	var profilesCfg ProfilesConfigs
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(ClusterProfileEndpoint)
	if err != nil {
		return ProfilesConfig{}, fmt.Errorf("call made to get cluster profiles errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return ProfilesConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &profilesCfg); err != nil {
		return ProfilesConfig{}, ResponseReadError(err.Error())
	}

	profilesCfg.ProfilesConfigs.ETAG = resp.Header().Get("ETag")

	return profilesCfg.ProfilesConfigs, nil
}

func (conf *client) GetClusterProfile(name string) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var profilesCfg CommonConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(ClusterProfileEndpoint, name))
	if err != nil {
		return CommonConfig{}, fmt.Errorf("call made to get cluster profile '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &profilesCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
	}

	profilesCfg.ETAG = resp.Header().Get("ETag")

	return profilesCfg, nil
}

func (conf *client) CreateClusterProfile(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var profileCfg CommonConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(ClusterProfileEndpoint)
	if err != nil {
		return CommonConfig{}, fmt.Errorf("call made to create cluster profile '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &profileCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
	}

	profileCfg.ETAG = resp.Header().Get("ETag")

	return profileCfg, nil
}

func (conf *client) UpdateClusterProfile(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var storeCfg CommonConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(filepath.Join(ClusterProfileEndpoint, config.ID))
	if err != nil {
		return CommonConfig{}, fmt.Errorf("call made to update cluster profile '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &storeCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
	}

	storeCfg.ETAG = resp.Header().Get("ETag")

	return storeCfg, nil
}

func (conf *client) DeleteClusterProfile(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Delete(filepath.Join(ClusterProfileEndpoint, name))
	if err != nil {
		return fmt.Errorf("call made to delete cluster profile '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
