package gocd

import (
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
	resp, err := newClient.httpClient.R().SetResult(&backUpConf).Get(GoCdBackupConfigEndpoint)
	if err != nil {
		return BackupConfig{}, fmt.Errorf("call made to get backup information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return BackupConfig{}, apiWithCodeError(resp.StatusCode())
	}

	return backUpConf, nil
}
