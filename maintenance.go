package gocd

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// EnableMaintenanceMode enables maintenance mode.
func (conf *client) EnableMaintenanceMode() error {
	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionOne,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(MaintenanceEndpoint, "enable"))
	if err != nil {
		return &errors.APIError{Err: err, Message: "enable maintenance mode"}
	}

	if resp.StatusCode() != http.StatusNoContent {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// GetMaintenanceModeInfo fetches the latest information of server maintenance mode information.
func (conf *client) GetMaintenanceModeInfo() (Maintenance, error) {
	newClient := conf.clone()

	var maintenanceInfo Maintenance

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(MaintenanceEndpoint, "info"))
	if err != nil {
		return Maintenance{}, &errors.APIError{Err: err, Message: "get maintenance mode information"}
	}

	if resp.StatusCode() != http.StatusOK {
		return Maintenance{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &maintenanceInfo); err != nil {
		return Maintenance{}, &errors.MarshalError{Err: err}
	}

	return maintenanceInfo, nil
}

// DisableMaintenanceMode disables the maintenance mode.
func (conf *client) DisableMaintenanceMode() error {
	newClient := conf.clone()

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionOne,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(MaintenanceEndpoint, "disable"))
	if err != nil {
		return &errors.APIError{Err: err, Message: "disable maintenance mode"}
	}

	if resp.StatusCode() != http.StatusNoContent {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
