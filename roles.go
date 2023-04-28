package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

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
		return RolesConfig{}, &errors.APIError{Err: err, Message: "get all roles"}
	}

	if resp.StatusCode() != http.StatusOK {
		return RolesConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &rolesCfg); err != nil {
		return RolesConfig{}, &errors.MarshalError{Err: err}
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
		return roleCfg, &errors.APIError{Err: err, Message: fmt.Sprintf("get role '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return roleCfg, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return roleCfg, &errors.MarshalError{Err: err}
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
		return RolesConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get role by type '%s'", roleType)}
	}

	if resp.StatusCode() != http.StatusOK {
		return RolesConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return RolesConfig{}, &errors.MarshalError{Err: err}
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
		return roleCfg, &errors.APIError{Err: err, Message: fmt.Sprintf("create role '%s'", config.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return roleCfg, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return roleCfg, &errors.MarshalError{Err: err}
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
		Put(filepath.Join(RolesEndpoint, config.Name))
	if err != nil {
		return roleCfg, &errors.APIError{Err: err, Message: fmt.Sprintf("update role '%s'", config.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return roleCfg, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &roleCfg); err != nil {
		return roleCfg, &errors.MarshalError{Err: err}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete role '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
