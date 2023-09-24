package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
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
		return nil, &errors.APIError{Err: err, Message: "get all packages"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &packagesCfg); err != nil {
		return nil, &errors.MarshalError{Err: err}
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
		return Package{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get package '%s'", repoID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return Package{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &packageCfg); err != nil {
		return Package{}, &errors.MarshalError{Err: err}
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
		return Package{}, &errors.APIError{Err: err, Message: fmt.Sprintf("create package '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return Package{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &packageCfg); err != nil {
		return Package{}, &errors.MarshalError{Err: err}
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
		return Package{}, &errors.APIError{Err: err, Message: fmt.Sprintf("update package '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return Package{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &packageCfg); err != nil {
		return Package{}, &errors.MarshalError{Err: err}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete package '%s'", repoID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
