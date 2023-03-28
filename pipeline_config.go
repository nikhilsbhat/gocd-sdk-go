package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

func (conf *client) GetPipelineConfig(name string) (PipelineConfig, error) {
	var pipelineConfig PipelineConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return pipelineConfig, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionEleven,
		}).
		Get(filepath.Join(PipelineConfigEndpoint, name))
	if err != nil {
		return pipelineConfig, &errors.APIError{Err: err, Message: fmt.Sprintf("get pipeline config '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineConfig, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var pipelineCfg map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &pipelineCfg); err != nil {
		return pipelineConfig, &errors.MarshalError{Err: err}
	}

	delete(pipelineCfg, "_links")
	delete(pipelineCfg, "origin")

	var config PipelineConfig

	config.Config = pipelineCfg
	config.ETAG = resp.Header().Get("ETag")

	return config, nil
}

func (conf *client) UpdatePipelineConfig(config PipelineConfig) (PipelineConfig, error) {
	var pipelineConfig PipelineConfig
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return pipelineConfig, err
	}

	pipelineName := config.Config["name"].(string)

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionEleven,
			"Content-Type": ContentJSON,
			"If-Match":     config.ETAG,
		}).
		SetBody(config.Config).
		Put(filepath.Join(PipelineConfigEndpoint, pipelineName))
	if err != nil {
		return pipelineConfig, &errors.APIError{Err: err, Message: fmt.Sprintf("update pipeline config '%s'", pipelineName)}
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineConfig, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var pipelineCfg map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &pipelineCfg); err != nil {
		return pipelineConfig, &errors.MarshalError{Err: err}
	}

	delete(pipelineCfg, "_links")
	delete(pipelineCfg, "origin")

	pipelineConfig.Config = pipelineCfg
	pipelineConfig.ETAG = resp.Header().Get("ETag")

	return pipelineConfig, nil
}

func (conf *client) CreatePipeline(config PipelineConfig) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	pipelineName := config.Config["name"].(string)

	defaultHeaders := map[string]string{
		"Accept":       HeaderVersionEleven,
		"Content-Type": ContentJSON,
	}

	if config.PausePipeline {
		defaultHeaders["X-pause-pipeline"] = "true"
	}

	if len(config.PauseReason) != 0 {
		defaultHeaders["X-pause-cause"] = config.PauseReason
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(defaultHeaders).
		SetBody(config.Config).
		Post(PipelineConfigEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("create pipeline config '%s'", pipelineName)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) DeletePipeline(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionEleven,
		}).
		Delete(filepath.Join(PipelineConfigEndpoint, name))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("delete pipeline config '%s'", name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

func (conf *client) ExtractTemplatePipeline(pipeline, template string) (PipelineConfig, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PipelineConfig{}, err
	}

	var pipelineConfig map[string]interface{}
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionEleven,
			"Content-Type": ContentJSON,
		}).
		SetBody(map[string]string{"template_name": template}).
		Put(filepath.Join(PipelineConfigEndpoint, pipeline, "extract_to_template"))
	if err != nil {
		return PipelineConfig{}, &errors.APIError{Err: err, Message: fmt.Sprintf("extracting template from pipeline '%s'", pipeline)}
	}

	if resp.StatusCode() != http.StatusOK {
		return PipelineConfig{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pipelineConfig); err != nil {
		return PipelineConfig{}, &errors.MarshalError{Err: err}
	}

	delete(pipelineConfig, "_links")
	delete(pipelineConfig, "origin")

	var config PipelineConfig

	config.Config = pipelineConfig
	config.ETAG = resp.Header().Get("ETag")

	return config, nil
}
