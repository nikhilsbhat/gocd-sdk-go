package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetUsers() ([]User, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var usersObj Users
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(UsersEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get all users"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &usersObj); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return usersObj.GoCDUsers.Users, nil
}

func (conf *client) GetUser(user string) (User, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return User{}, err
	}

	var userObj User
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Get(filepath.Join(UsersEndpoint, user))
	if err != nil {
		return userObj, &errors.APIError{Err: err, Message: "get user information"}
	}

	if resp.StatusCode() != http.StatusOK {
		return userObj, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &userObj); err != nil {
		return userObj, &errors.MarshalError{Err: err}
	}

	return userObj, nil
}

func (conf *client) CreateUser(user User) (User, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return User{}, err
	}

	var userConfig User
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(user).
		Post(UsersEndpoint)
	if err != nil {
		return userConfig, &errors.APIError{Err: err, Message: fmt.Sprintf("create user '%s'", user.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return userConfig, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &userConfig); err != nil {
		return userConfig, &errors.MarshalError{Err: err}
	}

	return userConfig, nil
}

func (conf *client) UpdateUser(user User) (User, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return User{}, err
	}

	var userConfig User
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(user).
		Patch(filepath.Join(UsersEndpoint, user.Name))
	if err != nil {
		return userConfig, &errors.APIError{Err: err, Message: fmt.Sprintf("update user '%s'", user.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return userConfig, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &userConfig); err != nil {
		return userConfig, &errors.MarshalError{Err: err}
	}

	return userConfig, nil
}

func (conf *client) DeleteUser(user string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionThree,
		}).
		Delete(filepath.Join(UsersEndpoint, user))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete user '%s'", user)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) BulkDeleteUsers(users map[string]interface{}) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(users).
		Delete(UsersEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "bulk delete users"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) BulkEnableDisableUsers(users map[string]interface{}) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionThree,
			"Content-Type": ContentJSON,
		}).
		SetBody(users).
		Patch(AdminOperationStateEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "bulk enable/disable users"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
