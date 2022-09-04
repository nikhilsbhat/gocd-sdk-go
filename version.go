package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetVersionInfo fetches version information of the GoCD to which it is connected to.
func (conf *client) GetVersionInfo() (VersionInfo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return VersionInfo{}, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})

	var version VersionInfo
	resp, err := newClient.httpClient.R().Get(GoCdVersionEndpoint)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("call made to get version information errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return VersionInfo{}, apiWithCodeError(resp.StatusCode())
	}
	if err := json.Unmarshal(resp.Body(), &version); err != nil {
		return VersionInfo{}, responseReadError(err.Error())
	}

	return version, nil
}
