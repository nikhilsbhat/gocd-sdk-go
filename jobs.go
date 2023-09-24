package gocd

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// GetScheduledJobs returns all scheduled jobs from GoCD.
func (conf *client) GetScheduledJobs() (ScheduledJobs, error) {
	var scheduledJobs ScheduledJobs

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return scheduledJobs, err
	}

	resp, err := newClient.httpClient.R().
		Get(APIJobFeedEndpoint)
	if err != nil {
		return scheduledJobs, &errors.APIError{Err: err, Message: "get scheduled jobs"}
	}

	if resp.StatusCode() != http.StatusOK {
		return scheduledJobs, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = xml.Unmarshal(resp.Body(), &scheduledJobs); err != nil {
		return scheduledJobs, &errors.MarshalError{Err: err}
	}

	return scheduledJobs, nil
}

// RunFailedJobs runs all failed jobs from a selected pipeline.
func (conf *client) RunFailedJobs(stage Stage) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionThree,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(StageEndpoint, stage.Pipeline, stage.PipelineInstance, stage.Name, stage.StageCounter, "run-failed-jobs"))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: "run failed jobs"}
	}

	if resp.StatusCode() != http.StatusAccepted {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var message map[string]string

	if err = json.Unmarshal(resp.Body(), &message); err != nil {
		return "", &errors.MarshalError{Err: err}
	}

	return message["message"], nil
}

// RunJobs runs all selected jobs from a selected pipeline.
func (conf *client) RunJobs(stage Stage) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetBody(stage).
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionThree,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(StageEndpoint, stage.Pipeline, stage.PipelineInstance, stage.Name, stage.StageCounter, "run-selected-jobs"))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: "run selected jobs"}
	}

	if resp.StatusCode() != http.StatusAccepted {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var message map[string]string

	if err = json.Unmarshal(resp.Body(), &message); err != nil {
		return "", &errors.MarshalError{Err: err}
	}

	return message["message"], nil
}
