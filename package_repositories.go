package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
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
		return nil, fmt.Errorf("call made to get package repositories errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &packageRepositoriesCfg); err != nil {
		return nil, ResponseReadError(err.Error())
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
		return PackageRepository{}, fmt.Errorf("call made to get package repository '%s' errored with: %w", repoID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return PackageRepository{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &repositoryCfg); err != nil {
		return PackageRepository{}, ResponseReadError(err.Error())
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
		return PackageRepository{}, fmt.Errorf("call made to create package repository '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return PackageRepository{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &repositoryCfg); err != nil {
		return PackageRepository{}, ResponseReadError(err.Error())
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
		return PackageRepository{}, fmt.Errorf("call made to update package repository '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return PackageRepository{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &repositoryCfg); err != nil {
		return PackageRepository{}, ResponseReadError(err.Error())
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
		return fmt.Errorf("call made to delete package repository '%s' errored with: %w", repoID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
