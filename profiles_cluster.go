package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetClusterProfiles() (ProfilesConfig, error) {
	newClient := conf.clone()

	var profilesCfg ProfilesConfigs

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(ClusterProfileEndpoint)
	if err != nil {
		return ProfilesConfig{}, &errors.APIError{Err: err, Message: "get cluster profiles"}
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
	newClient := conf.clone()

	var profilesCfg CommonConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(ClusterProfileEndpoint, name))
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get cluster profile '%s'", name)}
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
	newClient := conf.clone()

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
	newClient := conf.clone()

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
	newClient := conf.clone()

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
