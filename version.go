package gocd

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// GetVersionInfo fetches version information of the GoCD to which it is connected to.
func (conf *client) GetVersionInfo() (VersionInfo, error) {
	var version VersionInfo

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return version, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(VersionEndpoint)
	if err != nil {
		return version, &errors.APIError{Err: err, Message: "get version information"}
	}

	if resp.StatusCode() != http.StatusOK {
		return version, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &version); err != nil {
		return version, &errors.MarshalError{Err: err}
	}

	return version, nil
}
