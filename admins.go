package gocd

import (
	"encoding/json"
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

	var adminsConf SystemAdmins
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(SystemAdminEndpoint)

	if err != nil {
		return SystemAdmins{}, fmt.Errorf("call made to get system admin errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return SystemAdmins{}, APIWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &adminsConf); err != nil {
		return adminsConf, ResponseReadError(err.Error())
	}

	return adminsConf, nil
}
