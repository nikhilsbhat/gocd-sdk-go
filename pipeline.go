package gocd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"
)

// GetPipelines fetches all pipelines configured in GoCD server.
func (conf *client) GetPipelines() (PipelinesInfo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PipelinesInfo{}, err
	}

	var pipelinesInfo PipelinesInfo
	resp, err := newClient.httpClient.R().
		Get(APIFeedPipelineEndpoint)
	if err != nil {
		return PipelinesInfo{}, fmt.Errorf("call made to get pipelines errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelinesInfo{}, APIWithCodeError(resp.StatusCode())
	}

	if err := xml.Unmarshal(resp.Body(), &pipelinesInfo); err != nil {
		return PipelinesInfo{}, ResponseReadError(err.Error())
	}

	return pipelinesInfo, nil
}

// GetPipelineState fetches status of selected pipelines.
func (conf *client) GetPipelineState(pipeline string) (PipelineState, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PipelineState{}, err
	}

	var pipelinesStatus PipelineState
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(fmt.Sprintf(PipelineStatus, pipeline))
	if err != nil {
		return PipelineState{}, fmt.Errorf("call made to get pipeline state errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelineState{}, APIWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &pipelinesStatus); err != nil {
		return PipelineState{}, ResponseReadError(err.Error())
	}
	pipelinesStatus.Name = pipeline

	return pipelinesStatus, nil
}

// PipelinePause pauses specified pipeline with valid message passed.
func (conf *client) PipelinePause(name string, message any) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	msg := fmt.Sprintf("pausing pipeline %s", name)
	if message != nil {
		msg = message.(string)
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(map[string]string{
			"pause_cause": msg,
		}).
		Post(filepath.Join(PipelinesEndpoint, name, "pause"))
	if err != nil {
		return fmt.Errorf("call made to pause pipeline errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// PipelineUnPause unpauses specified pipeline.
func (conf *client) PipelineUnPause(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionOne,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(PipelinesEndpoint, name, "unpause"))
	if err != nil {
		return fmt.Errorf("call made to unpause pipeline errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// PipelineUnlock unlocks the specified locked pipeline.
func (conf *client) PipelineUnlock(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionOne,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(PipelinesEndpoint, name, "unlock"))
	if err != nil {
		return fmt.Errorf("call made to unlock pipeline errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// SchedulePipeline schedules the specified pipeline with specified configurations.
func (conf *client) SchedulePipeline(name string, schedule Schedule) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(schedule).
		Post(filepath.Join(PipelinesEndpoint, name, "schedule"))
	if err != nil {
		return fmt.Errorf("call made to schedule pipeline '%s' errored with %w", name, err)
	}

	if resp.StatusCode() != http.StatusAccepted {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// GetPipelineName parses pipeline url to fetch the pipeline name.
func GetPipelineName(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("parsing URL errored with %w", err)
	}

	return strings.TrimSuffix(strings.TrimPrefix(parsedURL.Path, "/go/api/feed/pipelines/"), "/stages.xml"), nil
}
