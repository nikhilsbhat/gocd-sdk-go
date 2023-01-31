package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

func (conf *client) GetPluginsInfo() (PluginsInfo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PluginsInfo{}, err
	}

	var pluginInfosCfg PluginsInfos
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Get(PluginInfoEndpoint)
	if err != nil {
		return PluginsInfo{}, fmt.Errorf("call made to get all plugins info errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return PluginsInfo{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &pluginInfosCfg); err != nil {
		return PluginsInfo{}, ResponseReadError(err.Error())
	}

	pluginInfosCfg.PluginsInfos.ETAG = resp.Header().Get("ETag")

	return pluginInfosCfg.PluginsInfos, nil
}

func (conf *client) GetPluginInfo(name string) (Plugin, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Plugin{}, err
	}

	var pluginInfoCfg Plugin
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Get(filepath.Join(PluginInfoEndpoint, name))
	if err != nil {
		return Plugin{}, fmt.Errorf("call made to get plugin info '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return Plugin{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &pluginInfoCfg); err != nil {
		return Plugin{}, ResponseReadError(err.Error())
	}

	pluginInfoCfg.ETAG = resp.Header().Get("ETag")

	return pluginInfoCfg, nil
}
