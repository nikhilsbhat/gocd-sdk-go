package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
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
			"Accept":       HeaderVersionOne,
			"Content-Type": ContentJSON,
		}).
		Post(PipelineGroupEndpoint)
	if err != nil {
		return fmt.Errorf("call made to create pipeline group '%s' information errored with %w", group.Name, err)
	}
	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
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
		return nil, fmt.Errorf("call made to get pipeline group information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIWithCodeError(resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &groupConf); err != nil {
		return nil, ResponseReadError(err.Error())
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
		return fmt.Errorf("call made to delete pipeline group errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
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
		return pipelineGroup, fmt.Errorf("call made to fetch pipeline group errored with %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return pipelineGroup, APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	if err = json.Unmarshal(resp.Body(), &pipelineGroup); err != nil {
		return pipelineGroup, ResponseReadError(err.Error())
	}

	pipelineGroup.ETAG = resp.Header().Get("ETag")

	return pipelineGroup, nil
}
