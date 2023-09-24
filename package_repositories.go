package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetPackageRepositories() ([]PackageRepository, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var packageRepositoriesCfg PackageRepositories

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(PackageRepositoriesEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get package repositories"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &packageRepositoriesCfg); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return packageRepositoriesCfg.Repositories.PackageRepositories, nil
}

func (conf *client) GetPackageRepository(repoID string) (PackageRepository, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PackageRepository{}, err
	}

	var repositoryCfg PackageRepository

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(PackageRepositoriesEndpoint, repoID))
	if err != nil {
		return PackageRepository{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get package repository '%s'", repoID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return PackageRepository{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &repositoryCfg); err != nil {
		return PackageRepository{}, &errors.MarshalError{Err: err}
	}

	repositoryCfg.ETAG = resp.Header().Get("ETag")

	return repositoryCfg, nil
}

func (conf *client) CreatePackageRepository(config PackageRepository) (PackageRepository, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PackageRepository{}, err
	}

	var repositoryCfg PackageRepository

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(PackageRepositoriesEndpoint)
	if err != nil {
		return PackageRepository{}, &errors.APIError{Err: err, Message: fmt.Sprintf("create package repository '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return PackageRepository{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &repositoryCfg); err != nil {
		return PackageRepository{}, &errors.MarshalError{Err: err}
	}

	repositoryCfg.ETAG = resp.Header().Get("ETag")

	return repositoryCfg, nil
}

func (conf *client) UpdatePackageRepository(config PackageRepository) (PackageRepository, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PackageRepository{}, err
	}

	var repositoryCfg PackageRepository

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(filepath.Join(PackageRepositoriesEndpoint, config.ID))
	if err != nil {
		return PackageRepository{}, &errors.APIError{Err: err, Message: fmt.Sprintf("update package repository '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return PackageRepository{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &repositoryCfg); err != nil {
		return PackageRepository{}, &errors.MarshalError{Err: err}
	}

	repositoryCfg.ETAG = resp.Header().Get("ETag")

	return repositoryCfg, nil
}

func (conf *client) DeletePackageRepository(repoID string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Delete(filepath.Join(PackageRepositoriesEndpoint, repoID))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete package repository '%s'", repoID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
