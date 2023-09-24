package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
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
		return SystemAdmins{}, &errors.APIError{Err: err, Message: "get system admin"}
	}

	if resp.StatusCode() != http.StatusOK {
		return SystemAdmins{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &adminsConf); err != nil {
		return adminsConf, &errors.MarshalError{Err: err}
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
		return admins, &errors.APIError{Err: err, Message: "update system admin"}
	}

	if resp.StatusCode() != http.StatusOK {
		return admins, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &admins); err != nil {
		return admins, &errors.MarshalError{Err: err}
	}

	admins.ETAG = resp.Header().Get("Etag")

	return admins, nil
}
