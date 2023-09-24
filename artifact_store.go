package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetArtifactStores() (ArtifactStoresConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ArtifactStoresConfig{}, err
	}

	var storeCfg ArtifactStoresConfigs

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(ArtifactStoreEndpoint)
	if err != nil {
		return ArtifactStoresConfig{}, &errors.APIError{Err: err, Message: "get artifact stores"}
	}

	if resp.StatusCode() != http.StatusOK {
		return ArtifactStoresConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &storeCfg); err != nil {
		return ArtifactStoresConfig{}, &errors.MarshalError{Err: err}
	}

	storeCfg.ArtifactStoresConfigs.ETAG = resp.Header().Get("ETag")

	return storeCfg.ArtifactStoresConfigs, nil
}

func (conf *client) GetArtifactStore(name string) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var storeCfg CommonConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(ArtifactStoreEndpoint, name))
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get artifact store %s", name)}
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

func (conf *client) CreateArtifactStore(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var storeCfg CommonConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(ArtifactStoreEndpoint)
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("create artifact store %s", config.ID)}
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

func (conf *client) UpdateArtifactStore(config CommonConfig) (CommonConfig, error) {
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
		Put(filepath.Join(ArtifactStoreEndpoint, config.ID))
	if err != nil {
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("update artifact store %s", config.ID)}
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

func (conf *client) DeleteArtifactStore(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Delete(filepath.Join(ArtifactStoreEndpoint, name))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete artifact store %s", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
