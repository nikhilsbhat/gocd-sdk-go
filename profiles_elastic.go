package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
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
		return ProfilesConfig{}, &errors.APIError{Err: err, Message: "get elastic agent profiles"}
	}

	if resp.StatusCode() != http.StatusOK {
		return ProfilesConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return ProfilesConfig{}, &errors.MarshalError{Err: err}
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
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get elastic agent profile '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
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
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("create elastic agent profile '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
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
		return CommonConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("update elastic agent profile '%s'", config.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return CommonConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &elasticAgentCfg); err != nil {
		return CommonConfig{}, &errors.MarshalError{Err: err}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete elastic agent profile '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) GetElasticAgentProfileUsage(profileID string) ([]ElasticProfileUsage, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var elasticProfileUsage []ElasticProfileUsage

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(fmt.Sprintf(ElasticProfileUsageEndpoint, profileID))
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("get elastic agent profile usage '%s'", profileID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &elasticProfileUsage); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	return elasticProfileUsage, nil
}
