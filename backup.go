package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetBackupInfo fetches information of backup configured in GoCD server.
func (conf *client) GetBackupInfo() (BackupConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return BackupConfig{}, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})

	var backUpConf BackupConfig
	resp, err := newClient.httpClient.R().Get(GoCdBackupConfigEndpoint)
	if err != nil {
		return BackupConfig{}, fmt.Errorf("call made to get backup information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, APIWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &backUpConf); err != nil {
		return BackupConfig{}, ResponseReadError(err.Error())
	}

	return backUpConf, nil
}
