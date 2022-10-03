package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetSystemAdmins fetches information of all system admins present in GoCD server.
func (conf *client) GetSystemAdmins() (SystemAdmins, error) {
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

	if err = json.Unmarshal(resp.Body(), &adminsConf); err != nil {
		return adminsConf, ResponseReadError(err.Error())
	}

	return adminsConf, nil
}

// UpdateSystemAdmins should update system admins and replace the system admins with roles and users in the request.
func (conf *client) UpdateSystemAdmins(data SystemAdmins) (SystemAdmins, error) {
	var admins SystemAdmins
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return admins, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
			"If-Match":     data.ETAG,
		}).
		Put(SystemAdminEndpoint)
	if err != nil {
		return admins, fmt.Errorf("call made to update system admin errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return admins, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &admins); err != nil {
		return admins, ResponseReadError(err.Error())
	}

	admins.ETAG = resp.Header().Get("Etag")

	return admins, nil
}
