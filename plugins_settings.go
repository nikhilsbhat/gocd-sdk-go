package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

// GetPluginSettings fetches the plugins settings of a specified plugin from GoCD.
func (conf *client) GetPluginSettings(name string) (PluginSettings, error) {
	var setting PluginSettings
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return setting, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(PluginSettingsEndpoint, name))
	if err != nil {
		return setting, fmt.Errorf("call made to get '%s' plugin setting errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return setting, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &setting); err != nil {
		return setting, ResponseReadError(err.Error())
	}

	setting.ETAG = resp.Header().Get("ETag")

	return setting, nil
}

// CreatePluginSettings creates the plugins settings of a specified plugin in GoCD.
func (conf *client) CreatePluginSettings(settings PluginSettings) (PluginSettings, error) {
	var setting PluginSettings
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return setting, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(settings).
		Post(PluginSettingsEndpoint)
	if err != nil {
		return setting, fmt.Errorf("call made to create plugin setting of '%s' errored with: %w", settings.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return setting, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &setting); err != nil {
		return setting, ResponseReadError(err.Error())
	}

	return setting, nil
}

// UpdatePluginSettings updates the plugins settings of an already existing plugin in GoCD.
func (conf *client) UpdatePluginSettings(settings PluginSettings) (PluginSettings, error) {
	var setting PluginSettings
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return setting, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
			"If-Match":     settings.ETAG,
		}).
		SetBody(settings).
		Put(filepath.Join(PluginSettingsEndpoint, settings.ID))
	if err != nil {
		return setting, fmt.Errorf("call made to update plugin setting of '%s' errored with: %w", settings.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return setting, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &setting); err != nil {
		return setting, ResponseReadError(err.Error())
	}

	return setting, nil
}
