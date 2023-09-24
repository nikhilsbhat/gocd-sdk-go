package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
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
		SetQueryParam("include_bad", "true").
		Get(PluginInfoEndpoint)
	if err != nil {
		return PluginsInfo{}, &errors.APIError{Err: err, Message: "get all plugins info"}
	}

	if resp.StatusCode() != http.StatusOK {
		return PluginsInfo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pluginInfosCfg); err != nil {
		return PluginsInfo{}, &errors.MarshalError{Err: err}
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
		return Plugin{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get plugin info '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return Plugin{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pluginInfoCfg); err != nil {
		return Plugin{}, &errors.MarshalError{Err: err}
	}

	pluginInfoCfg.ETAG = resp.Header().Get("ETag")

	return pluginInfoCfg, nil
}
