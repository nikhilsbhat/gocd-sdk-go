package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
)

func (conf *client) GetRoles() (RolesConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return RolesConfig{}, err
	}

	var rolesCfg RolesConfigs
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(RolesEndpoint)
	if err != nil {
		return RolesConfig{}, fmt.Errorf("call made to get all roles errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return RolesConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &rolesCfg); err != nil {
		return RolesConfig{}, ResponseReadError(err.Error())
	}

	rolesCfg.RolesConfigs.ETAG = resp.Header().Get("ETag")

	return rolesCfg.RolesConfigs, nil
}

func (conf *client) GetRole(name string) (Role, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Role{}, err
	}

	var roleCfg Role
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(filepath.Join(RolesEndpoint, name))
	if err != nil {
		return roleCfg, fmt.Errorf("call made to get role '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return roleCfg, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return roleCfg, ResponseReadError(err.Error())
	}

	roleCfg.ETAG = resp.Header().Get("ETag")

	return roleCfg, nil
}

func (conf *client) GetRolesByType(roleType string) (RolesConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return RolesConfig{}, err
	}

	var roleCfg RolesConfigs
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		SetQueryParam("type", strings.ToLower(roleType)).
		Get(RolesEndpoint)
	if err != nil {
		return RolesConfig{}, fmt.Errorf("call made to get role by type '%s' errored with: %w", roleType, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return RolesConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return RolesConfig{}, ResponseReadError(err.Error())
	}

	roleCfg.RolesConfigs.ETAG = resp.Header().Get("ETag")

	return roleCfg.RolesConfigs, nil
}

func (conf *client) CreateRole(config Role) (Role, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Role{}, err
	}

	var roleCfg Role
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(RolesEndpoint)
	if err != nil {
		return roleCfg, fmt.Errorf("call made to create role '%s' errored with: %w", config.Name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return roleCfg, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return roleCfg, ResponseReadError(err.Error())
	}

	roleCfg.ETAG = resp.Header().Get("ETag")

	return roleCfg, nil
}

func (conf *client) UpdateRole(config Role) (Role, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Role{}, err
	}

	var roleCfg Role
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(RolesEndpoint)
	if err != nil {
		return roleCfg, fmt.Errorf("call made to update role '%s' errored with: %w", config.Name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return roleCfg, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return roleCfg, ResponseReadError(err.Error())
	}

	roleCfg.ETAG = resp.Header().Get("ETag")

	return roleCfg, nil
}

func (conf *client) DeleteRole(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Delete(filepath.Join(RolesEndpoint, name))
	if err != nil {
		return fmt.Errorf("call made to delete role '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
