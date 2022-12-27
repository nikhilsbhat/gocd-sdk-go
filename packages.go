package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

func (conf *client) GetPackages() ([]Package, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var packagesCfg Packages
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(PackagesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get all packages errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &packagesCfg); err != nil {
		return nil, ResponseReadError(err.Error())
	}

	return packagesCfg.Packages.Packages, nil
}

func (conf *client) GetPackage(repoID string) (Package, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Package{}, err
	}

	var packageCfg Package
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(filepath.Join(PackagesEndpoint, repoID))
	if err != nil {
		return Package{}, fmt.Errorf("call made to get package '%s' errored with: %w", repoID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return Package{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &packageCfg); err != nil {
		return Package{}, ResponseReadError(err.Error())
	}

	packageCfg.ETAG = resp.Header().Get("ETag")

	return packageCfg, nil
}

func (conf *client) CreatePackage(config Package) (Package, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Package{}, err
	}

	var packageCfg Package
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(PackagesEndpoint)
	if err != nil {
		return Package{}, fmt.Errorf("call made to create package '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return Package{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &packageCfg); err != nil {
		return Package{}, ResponseReadError(err.Error())
	}

	packageCfg.ETAG = resp.Header().Get("ETag")

	return packageCfg, nil
}

func (conf *client) UpdatePackage(config Package) (Package, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Package{}, err
	}

	var packageCfg Package
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(filepath.Join(PackagesEndpoint, config.ID))
	if err != nil {
		return Package{}, fmt.Errorf("call made to update package '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return Package{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &packageCfg); err != nil {
		return Package{}, ResponseReadError(err.Error())
	}

	packageCfg.ETAG = resp.Header().Get("ETag")

	return packageCfg, nil
}

func (conf *client) DeletePackage(repoID string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Delete(filepath.Join(PackagesEndpoint, repoID))
	if err != nil {
		return fmt.Errorf("call made to delete package '%s' errored with: %w", repoID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
