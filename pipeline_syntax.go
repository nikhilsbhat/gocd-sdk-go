package gocd

import (
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/plugin"
)

func (conf *client) ValidatePipelineSyntax(pluginCfg plugin.Plugin, pipelines []string) (bool, error) {
	if err := pluginCfg.Type(pipelines); err != nil {
		return false, err
	}

	if _, err := pluginCfg.Download(); err != nil {
		return false, err
	}

	return pluginCfg.ValidatePlugin(pipelines)
}
