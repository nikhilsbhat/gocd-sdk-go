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

	if err = xml.Unmarshal(resp.Body(), &pipelinesInfo); err != nil {
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

	if err = json.Unmarshal(resp.Body(), &pipelinesStatus); err != nil {
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

// CommentOnPipeline publishes comment on specified pipeline.
func (conf *client) CommentOnPipeline(comment PipelineObject) error {
	if len(comment.Message) == 0 {
		return fmt.Errorf("comment message cannot be empty") //nolint:goerr113
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
		return fmt.Errorf("call made to comment on pipeline '%s' errored with %w", comment.Name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
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
		return pipelineInstance, fmt.Errorf("call made to fetch pipeline instance '%s' errored with %w", pipeline.Name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineInstance, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &pipelineInstance); err != nil {
		return pipelineInstance, ResponseReadError(err.Error())
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
			return history, fmt.Errorf("call made to fetch pipeline history '%s' errored with %w", name, err)
		}

		if resp.StatusCode() != http.StatusOK {
			return history, APIErrorWithBody(resp.String(), resp.StatusCode())
		}

		if err = json.Unmarshal(resp.Body(), &pipelineHistory); err != nil {
			return history, ResponseReadError(err.Error())
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
