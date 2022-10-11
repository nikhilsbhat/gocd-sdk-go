package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		return BackupConfig{}, fmt.Errorf("call made to get backup information errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &backUpConf); err != nil {
		return BackupConfig{}, ResponseReadError(err.Error())
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
		return fmt.Errorf("call made to create/udpate backup configuration errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
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
		return fmt.Errorf("call made to get backup information errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
