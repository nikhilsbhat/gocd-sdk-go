package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
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
		return ArtifactStoresConfig{}, fmt.Errorf("call made to get artifact stores errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return ArtifactStoresConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &storeCfg); err != nil {
		return ArtifactStoresConfig{}, ResponseReadError(err.Error())
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
		return CommonConfig{}, fmt.Errorf("call made to get artifact store %s errored with: %w", name, err)
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
		return CommonConfig{}, fmt.Errorf("call made to create artifact store %s errored with: %w", config.ID, err)
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
		return CommonConfig{}, fmt.Errorf("call made to update artifact store %s errored with: %w", config.ID, err)
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
		return fmt.Errorf("call made to delete artifact store %s errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
