package gocd

import (
	"strings"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/plugin"
)

func (conf *client) ValidatePipelineSyntax(pluginCfg plugin.Plugin, pipelines []string, fetchVersionFromServer bool) (bool, error) {
	if err := pluginCfg.SetType(pipelines); err != nil {
		return false, err
	}

	if fetchVersionFromServer {
		conf.logger.Info("since fetch version from server is enabled, fetching the plugin version from GoCD server")

		pluginsInfo, err := conf.GetPluginsInfo()
		if err != nil {
			return false, err
		}

		for _, pluginInfo := range pluginsInfo.Plugins {
			if strings.Contains(pluginInfo.ID, pluginCfg.GetType()) {
				pluginVersion := pluginInfo.About["version"].(string)

				pluginCfg.SetVersion(pluginVersion)

				conf.logger.Infof("identified the plugin as '%s' and the version installed in GoCD is '%s'",
					pluginCfg.GetType(), pluginCfg.GetVersion())
			}
		}
	}

	if _, err := pluginCfg.Download(); err != nil {
		return false, err
	}

	return pluginCfg.ValidatePlugin(pipelines)
}
