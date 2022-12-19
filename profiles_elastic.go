package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

func (conf *client) GetElasticAgentProfiles() (ProfilesConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return ProfilesConfig{}, err
	}

	var elasticAgentCfg ProfilesConfigs
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(AgentProfileEndpoint)
	if err != nil {
		return ProfilesConfig{}, fmt.Errorf("call made to get elastic agent profiles errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return ProfilesConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return ProfilesConfig{}, ResponseReadError(err.Error())
	}

	elasticAgentCfg.ProfilesConfigs.ETAG = resp.Header().Get("ETag")

	return elasticAgentCfg.ProfilesConfigs, nil
}

func (conf *client) GetElasticAgentProfile(name string) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var elasticAgentCfg CommonConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Get(filepath.Join(AgentProfileEndpoint, name))
	if err != nil {
		return CommonConfig{}, fmt.Errorf("call made to get elastic agent profile '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
	}

	elasticAgentCfg.ETAG = resp.Header().Get("ETag")

	return elasticAgentCfg, nil
}

func (conf *client) CreateElasticAgentProfile(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var elasticAgentCfg CommonConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
		}).
		SetBody(config).
		Post(AgentProfileEndpoint)
	if err != nil {
		return CommonConfig{}, fmt.Errorf("call made to create elastic agent profile '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
	}

	elasticAgentCfg.ETAG = resp.Header().Get("ETag")

	return elasticAgentCfg, nil
}

func (conf *client) UpdateElasticAgentProfile(config CommonConfig) (CommonConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return CommonConfig{}, err
	}

	var elasticAgentCfg CommonConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionTwo,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config).
		Put(filepath.Join(AgentProfileEndpoint, config.ID))
	if err != nil {
		return CommonConfig{}, fmt.Errorf("call made to update elastic agent profile '%s' errored with: %w", config.ID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return CommonConfig{}, ResponseReadError(err.Error())
	}

	elasticAgentCfg.ETAG = resp.Header().Get("ETag")

	return elasticAgentCfg, nil
}

func (conf *client) DeleteElasticAgentProfile(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionTwo,
		}).
		Delete(filepath.Join(AgentProfileEndpoint, name))
	if err != nil {
		return fmt.Errorf("call made to delete elastic agent profile '%s' errored with: %w", name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}
