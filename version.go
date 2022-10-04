package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
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
		return version, fmt.Errorf("call made to get version information errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return version, APIWithCodeError(resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &version); err != nil {
		return version, ResponseReadError(err.Error())
	}

	return version, nil
}
