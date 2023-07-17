package gocd

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// RunStage runs a selected stage from an appropriate pipeline.
func (conf *client) RunStage(stage Stage) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionTwo,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(StageEndpoint, stage.Pipeline, stage.PipelineInstance, stage.Name, "run"))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: "run stage"}
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

// CancelStage cancels the selected stage from a selected pipeline.
func (conf *client) CancelStage(stage Stage) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionThree,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(StageEndpoint, stage.Pipeline, stage.PipelineInstance, stage.Name, stage.StageCounter, "cancel"))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: "cancel stage"}
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	var message map[string]string

	if err = json.Unmarshal(resp.Body(), &message); err != nil {
		return "", &errors.MarshalError{Err: err}
	}

	return message["message"], nil
}
