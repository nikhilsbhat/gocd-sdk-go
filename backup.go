package gocd

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
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

// GetBackup gets the information of the backup which was taken earlier.
func (conf *client) GetBackup(backupID string) (BackupStats, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return BackupStats{}, err
	}

	var backUpStats BackupStats

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(filepath.Join(BackupStatsEndpoint, backupID))
	if err != nil {
		return backUpStats, &errors.APIError{Err: err, Message: "get backup stats"}
	}

	if resp.StatusCode() != http.StatusOK {
		return backUpStats, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &backUpStats); err != nil {
		return backUpStats, &errors.MarshalError{Err: err}
	}

	return backUpStats, nil
}

func (conf *client) ScheduleBackup() (map[string]string, error) {
	var backupStats map[string]string

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return backupStats, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionTwo,
			HeaderConfirm: "true",
		}).
		Post(BackupStatsEndpoint)
	if err != nil {
		return backupStats, &errors.APIError{Err: err, Message: "schedule backup"}
	}

	if resp.StatusCode() != http.StatusAccepted {
		return backupStats, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if len(resp.Header().Get(LocationHeader)) == 0 {
		return backupStats, &errors.NilHeaderError{Header: "Location", Message: "getting backup stats"}
	}

	_, backUpID := filepath.Split(resp.Header().Get(LocationHeader))

	backupStats = map[string]string{
		"BackUpID":   backUpID,
		"RetryAfter": resp.Header().Get("Retry-After"),
	}

	return backupStats, nil
}
