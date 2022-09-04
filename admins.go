package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetAdminsInfo fetches information of all system admins present in GoCD server.
func (conf *client) GetAdminsInfo() (SystemAdmins, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return SystemAdmins{}, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionTwo,
	})

	var adminsConf SystemAdmins
	resp, err := newClient.httpClient.R().SetResult(&adminsConf).Get(GoCdSystemAdminEndpoint)
	if err != nil {
		return SystemAdmins{}, fmt.Errorf("call made to get system admin errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return SystemAdmins{}, apiWithCodeError(resp.StatusCode())
	}

	return adminsConf, nil
}
