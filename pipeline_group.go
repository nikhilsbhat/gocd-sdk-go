package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
)

// Groups implements methods that help in fetching several other information from PipelineGroup.
type Groups []PipelineGroup

// CreatePipelineGroup will create pipeline group with provided configurations.
func (conf *client) CreatePipelineGroup(group PipelineGroup) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			HeaderConfirm:  "true",
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		SetBody(group).
		Post(PipelineGroupEndpoint)
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("create pipeline group '%s'", group.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// GetPipelineGroups fetches information of backup configured in GoCD server.
func (conf *client) GetPipelineGroups() ([]PipelineGroup, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var groupConf PipelineGroupsConfig

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(PipelineGroupEndpoint)
	if err != nil {
		return nil, &errors.APIError{Err: err, Message: "get pipeline groups information"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &groupConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
	}

	updatedGroupConf := make([]PipelineGroup, 0)
	for _, group := range groupConf.PipelineGroups.PipelineGroups {
		updatedGroupConf = append(updatedGroupConf, PipelineGroup{
			Name:          group.Name,
			PipelineCount: len(group.Pipelines),
			Pipelines:     group.Pipelines,
		})
	}

	return updatedGroupConf, nil
}

// DeletePipelineGroup deletes the specified pipeline group present in GoCD.
func (conf *client) DeletePipelineGroup(name string) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Delete(filepath.Join(PipelineGroupEndpoint, name))
	if err != nil {
		return &errors.APIError{Err: err, Message: "delete pipeline group"}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}

// Count return the total number of pipelines present.
func (conf Groups) Count() int {
	var pipelines int
	for _, i := range conf {
		pipelines += i.PipelineCount
	}

	return pipelines
}

// GetPipelineGroup fetches information of a specific pipeline group.
func (conf *client) GetPipelineGroup(name string) (PipelineGroup, error) {
	var pipelineGroup PipelineGroup

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return pipelineGroup, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		Get(filepath.Join(PipelineGroupEndpoint, name))
	if err != nil {
		return pipelineGroup, &errors.APIError{Err: err, Message: "fetch pipeline group"}
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineGroup, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pipelineGroup); err != nil {
		return pipelineGroup, &errors.MarshalError{Err: err}
	}

	pipelineGroup.ETAG = resp.Header().Get("ETag")

	return pipelineGroup, nil
}

// UpdatePipelineGroup updates the specified pipeline group with the latest config provided.
func (conf *client) UpdatePipelineGroup(group PipelineGroup) (PipelineGroup, error) {
	var pipelineGroup PipelineGroup

	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return pipelineGroup, err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
			"If-Match":     group.ETAG,
		}).
		SetBody(group).
		Put(filepath.Join(PipelineGroupEndpoint, group.Name))
	if err != nil {
		return pipelineGroup, &errors.APIError{Err: err, Message: fmt.Sprintf("update pipeline group '%s'", group.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineGroup, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &pipelineGroup); err != nil {
		return pipelineGroup, &errors.MarshalError{Err: err}
	}

	pipelineGroup.ETAG = resp.Header().Get("ETag")

	return pipelineGroup, nil
}
