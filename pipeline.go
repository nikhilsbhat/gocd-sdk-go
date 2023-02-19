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

// GetPipelineHistory fetches the history of a selected pipeline with counter.
func (conf *client) GetPipelineHistory(name string, defaultSize, defaultAfter int) ([]map[string]interface{}, error) {
	var history []map[string]interface{}
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return history, err
	}

	paginate := true
	size := defaultSize
	after := defaultAfter

	for paginate {
		var pipelineHistory PipelineHistory
		resp, err := newClient.httpClient.R().
			SetQueryParams(map[string]string{
				"page_size": strconv.Itoa(size),
				"after":     strconv.Itoa(after),
			}).
			SetHeaders(map[string]string{
				"Accept": HeaderVersionOne,
			}).
			Get(filepath.Join(PipelinesEndpoint, name, "history"))
		if err != nil {
			return history, &errors.APIError{Err: err, Message: fmt.Sprintf("fetch pipeline history '%s'", name)}
		}

		if resp.StatusCode() != http.StatusOK {
			return history, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
		}

		if err = json.Unmarshal(resp.Body(), &pipelineHistory); err != nil {
			return history, &errors.MarshalError{Err: err}
		}

		if (len(pipelineHistory.Pipelines) == 0) || (pipelineHistory.Links["next"] == nil) {
			conf.logger.Debug("no more pages to paginate, moving out of loop")
			paginate = false
		}

		after = size
		size += defaultSize

		history = append(history, pipelineHistory.Pipelines...)
	}

	return history, nil
}
