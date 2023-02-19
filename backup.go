package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

// GetBackupConfig fetches information of backup configured in GoCD server.
func (conf *client) GetBackupConfig() (BackupConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return BackupConfig{}, err
	}

	var backUpConf BackupConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(BackupConfigEndpoint)
	if err != nil {
		return BackupConfig{}, &errors.APIError{Err: err, Message: "get backup information"}
	}

	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &backUpConf); err != nil {
		return BackupConfig{}, &errors.MarshalError{Err: err}
	}

	return backUpConf, nil
}

// CreateOrUpdateBackupConfig will either create or update the config repo, it creates one if not created else update the existing with newer configuration.
func (conf *client) CreateOrUpdateBackupConfig(backup BackupConfig) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(backup).
		Post(BackupConfigEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "create/update backup configuration"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// DeleteBackupConfig deletes the backup config configured in GoCD.
func (conf *client) DeleteBackupConfig() error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Delete(BackupConfigEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: "delete backup configuration"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
