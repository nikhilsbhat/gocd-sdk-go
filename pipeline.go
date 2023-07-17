package gocd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

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
		return PipelinesInfo{}, &errors.APIError{Err: err, Message: "get pipelines"}
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelinesInfo{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = xml.Unmarshal(resp.Body(), &pipelinesInfo); err != nil {
		return PipelinesInfo{}, &errors.MarshalError{Err: err}
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
		return PipelineState{}, &errors.APIError{Err: err, Message: "get pipeline state"}
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelineState{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pipelinesStatus); err != nil {
		return PipelineState{}, &errors.MarshalError{Err: err}
	}
	pipelinesStatus.Name = pipeline

	return pipelinesStatus, nil
}

// GetPipelineRunHistory fetches all run history of selected pipeline from GoCD server.
// This would be an expensive operation; make sure to run it during non-peak hours.
func (conf *client) GetPipelineRunHistory(pipeline, pageSize string, delay time.Duration) ([]PipelineRunHistory, error) {
	type runHistory struct {
		Links     map[string]interface{} `json:"_links,omitempty" yaml:"_links,omitempty"`
		Pipelines []PipelineRunHistory   `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
	}

	pipelineRunHistories := make([]PipelineRunHistory, 0)

	after := "0"

	for {
		newClient := &client{}
		if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
			return nil, err
		}

		var pipelineRunHistory runHistory

		resp, err := newClient.httpClient.R().
			SetHeaders(map[string]string{
				"Accept":       HeaderVersionOne,
				"Content-Type": ContentJSON,
			}).
			SetQueryParams(map[string]string{
				"page_size": pageSize,
				"after":     after,
			}).
			Get(filepath.Join(PipelinesEndpoint, pipeline, "history"))
		if err != nil {
			return nil, &errors.APIError{Err: err, Message: fmt.Sprintf("get pipeline %s", pipeline)}
		}

		if resp.StatusCode() != http.StatusOK {
			return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
		}

		if err = json.Unmarshal(resp.Body(), &pipelineRunHistory); err != nil {
			return nil, &errors.MarshalError{Err: err}
		}

		if nextLnk := pipelineRunHistory.Links["next"]; nextLnk == nil {
			pipelineRunHistories = append(pipelineRunHistories, pipelineRunHistory.Pipelines...)

			break
		}

		nextLink := pipelineRunHistory.Links["next"].(map[string]interface{})["href"].(string)
		after = strings.Split(nextLink, "after=")[1]

		pipelineRunHistories = append(pipelineRunHistories, pipelineRunHistory.Pipelines...)

		time.Sleep(delay)
	}

	return pipelineRunHistories, nil
}

// GetPipelineSchedules fetches the last X schedules of the selected pipeline from GoCD server.
func (conf *client) GetPipelineSchedules(pipeline, start, perPage string) (PipelineSchedules, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PipelineSchedules{}, err
	}

	var pipelineSchedules PipelineSchedules
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionZero,
			"Content-Type": ContentJSON,
		}).
		SetQueryParams(map[string]string{
			"start":   start,
			"perPage": perPage,
		}).
		Get(fmt.Sprintf(LastXPipelineScheduledDates, pipeline))
	if err != nil {
		return PipelineSchedules{}, &errors.APIError{Err: err, Message: fmt.Sprintf("get pipeline schedules %s", pipeline)}
	}

	if resp.StatusCode() != http.StatusOK {
		return PipelineSchedules{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pipelineSchedules); err != nil {
		return PipelineSchedules{}, &errors.MarshalError{Err: err}
	}

	return pipelineSchedules, nil
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
		return &errors.APIError{Err: err, Message: "pause pipeline"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		return &errors.APIError{Err: err, Message: "unpause pipeline"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		return &errors.APIError{Err: err, Message: "unlock pipeline"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("schedule pipeline '%s'", name)}
	}

	if resp.StatusCode() != http.StatusAccepted {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// GetPipelineName parses pipeline url to fetch the pipeline name.
func GetPipelineName(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", &errors.GoCDError{Message: "parsing URL errored with:", Err: err}
	}

	return strings.TrimSuffix(strings.TrimPrefix(parsedURL.Path, PipelinePrefix), PipelineSuffix), nil
}

// CommentOnPipeline publishes comment on specified pipeline.
func (conf *client) CommentOnPipeline(comment PipelineObject) error {
	if len(comment.Message) == 0 {
		return &errors.GoCDError{Message: "comment message cannot be empty"}
	}

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(map[string]string{"comment": comment.Message}).
		Post(filepath.Join(PipelinesEndpoint, comment.Name, strconv.Itoa(comment.Counter), "comment"))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("comment on pipeline '%s'", comment.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// GetPipelineInstance fetches the instance of a selected pipeline with counter.
func (conf *client) GetPipelineInstance(pipeline PipelineObject) (map[string]interface{}, error) {
	var pipelineInstance map[string]interface{}
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return pipelineInstance, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(PipelinesEndpoint, pipeline.Name, strconv.Itoa(pipeline.Counter)))
	if err != nil {
		return pipelineInstance, &errors.APIError{Err: err, Message: fmt.Sprintf("fetch pipeline instance '%s'", pipeline.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineInstance, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pipelineInstance); err != nil {
		return pipelineInstance, &errors.MarshalError{Err: err}
	}

	return pipelineInstance, nil
}

func (conf *client) ExportPipelineToConfigRepoFormat(pipelineName, pluginID string) (PipelineExport, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PipelineExport{}, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		SetQueryParams(map[string]string{
			"plugin_id": pluginID,
		}).
		Get(filepath.Join(PipelineExportEndpoint, pipelineName))
	if err != nil {
		return PipelineExport{}, &errors.APIError{Err: err, Message: fmt.Sprintf("export pipeline '%s' to format '%s'", pipelineName, pluginID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return PipelineExport{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	disposition := strings.Split(resp.Header().Get("Content-Disposition"), "=")

	pipelineExport := PipelineExport{
		PluginID:         pluginID,
		PipelineFileName: strings.Trim(disposition[1], `"`),
		PipelineContent:  resp.String(),
		ETAG:             resp.Header().Get("ETag"),
	}

	return pipelineExport, nil
}
