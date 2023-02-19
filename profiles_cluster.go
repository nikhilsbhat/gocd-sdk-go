package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

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
		return ProfilesConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &profilesCfg); err != nil {
		return ProfilesConfig{}, &errors.MarshalError{Err: err}
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
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &profilesCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
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
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("create cluster profile '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &profileCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
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
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("update cluster profile '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &storeCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete cluster profile '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
