package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

// EnableMaintenanceMode enables maintenance mode.
func (conf *client) EnableMaintenanceMode() error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":         HeaderVersionOne,
			"X-GoCD-Confirm": "true",
		}).
		Post(filepath.Join(MaintenanceEndpoint, "enable"))
	if err != nil {
		return fmt.Errorf("call made to enable maintenance mode errored with %w", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// GetMaintenanceModeInfo fetches the latest information of server maintenance mode information.
func (conf *client) GetMaintenanceModeInfo() (Maintenance, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Maintenance{}, err
	}

	var maintenanceInfo Maintenance
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).Get(filepath.Join(MaintenanceEndpoint, "info"))
	if err != nil {
		return Maintenance{}, fmt.Errorf("call made to enable maintenance mode errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return Maintenance{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &maintenanceInfo); err != nil {
		return Maintenance{}, ResponseReadError(err.Error())
	}

	return maintenanceInfo, nil
}

// DisableMaintenanceMode disables the maintenance mode.
func (conf *client) DisableMaintenanceMode() error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":         HeaderVersionOne,
			"X-GoCD-Confirm": "true",
		}).
		Post(filepath.Join(MaintenanceEndpoint, "disable"))
	if err != nil {
		return fmt.Errorf("call made to enable maintenance mode errored with %w", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
