package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

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

	for _, plugin := range pluginInfosCfg.PluginsInfos.Plugins {
		correctKeys(plugin)
	}

	pluginInfosCfg.PluginsInfos.ETAG = resp.Header().Get("ETag")

	return pluginInfosCfg.PluginsInfos, nil
}

func (conf *client) GetPluginInfo(name string) (Plugin, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return Plugin{}, err
	}

	var pluginInfoCfg *Plugin

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

	correctKeys(pluginInfoCfg)

	pluginInfoCfg.ETAG = resp.Header().Get("ETag")

	return *pluginInfoCfg, nil
}

func correctKeys(plugin *Plugin) {
	for _, extension := range plugin.Extensions {
		if extension.AuthConfigSettings != nil {
			convertKeysToSnakeCase(extension.AuthConfigSettings.Configurations)
		}

		if extension.ArtifactConfigSettings != nil {
			convertKeysToSnakeCase(extension.ArtifactConfigSettings.Configurations)
		}

		if extension.ElasticAgentProfileSettings != nil {
			convertKeysToSnakeCase(extension.ElasticAgentProfileSettings.Configurations)
		}

		if extension.FetchArtifactSettings != nil {
			convertKeysToSnakeCase(extension.FetchArtifactSettings.Configurations)
		}

		if extension.ClusterProfileSettings != nil {
			convertKeysToSnakeCase(extension.ClusterProfileSettings.Configurations)
		}

		if extension.PluginSettings != nil {
			convertKeysToSnakeCase(extension.PluginSettings.Configurations)
		}

		if extension.PackageSettings != nil {
			convertKeysToSnakeCase(extension.PackageSettings.Configurations)
		}

		if extension.RepositorySettings != nil {
			convertKeysToSnakeCase(extension.RepositorySettings.Configurations)
		}

		if extension.ScmSettings != nil {
			convertKeysToSnakeCase(extension.ScmSettings.Configurations)
		}

		if extension.StoreConfigSettings != nil {
			convertKeysToSnakeCase(extension.StoreConfigSettings.Configurations)
		}

		if extension.SecretConfigSettings != nil {
			convertKeysToSnakeCase(extension.SecretConfigSettings.Configurations)
		}

		if extension.RoleSettings != nil {
			convertKeysToSnakeCase(extension.RoleSettings.Configurations)
		}

		if extension.TaskSettings != nil {
			convertKeysToSnakeCase(extension.TaskSettings.Configurations)
		}
	}
}

func convertKeysToSnakeCase(input []*PluginConfiguration) {
	for _, configuration := range input {
		configuration.Key = camelToSnake(configuration.Key)
	}
}

func camelToSnake(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")

	return strings.ToLower(snake)
}
