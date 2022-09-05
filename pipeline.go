package gocd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jinzhu/copier"
)

// Groups implements methods that help in fetching several other information from PipelineGroup.
type Groups []PipelineGroup

// GetPipelineGroupInfo fetches information of backup configured in GoCD server.
func (conf *client) GetPipelineGroupInfo() ([]PipelineGroup, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})

	var groupConf PipelineGroupsConfig
	resp, err := newClient.httpClient.R().Get(GoCdPipelineGroupEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get pipeline group information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, ApiWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &groupConf); err != nil {
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

// GetPipelines fetches all pipelines configured in GoCD server.
func (conf *client) GetPipelines() (PipelinesInfo, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return PipelinesInfo{}, err
	}

	var pipelinesInfo PipelinesInfo
	resp, err := newClient.httpClient.R().Get(GoCdAPIFeedPipelineEndpoint)
	if err != nil {
		return PipelinesInfo{}, fmt.Errorf("call made to get pipelines errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelinesInfo{}, ApiWithCodeError(resp.StatusCode())
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

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})

	var pipelinesStatus PipelineState
	resp, err := newClient.httpClient.R().Get(fmt.Sprintf(GoCdPipelineStatus, pipeline))
	if err != nil {
		return PipelineState{}, fmt.Errorf("call made to get pipeline state errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelineState{}, ApiWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &pipelinesStatus); err != nil {
		return PipelineState{}, ResponseReadError(err.Error())
	}
	pipelinesStatus.Name = pipeline

	return pipelinesStatus, nil
}

// Count return the total number of pipelines present.
func (conf Groups) Count() int {
	var pipelines int
	for _, i := range conf {
		pipelines += i.PipelineCount
	}

	return pipelines
}

// GetPipelineName parses pipeline url to fetch the pipeline name.
func GetPipelineName(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("parsing URL errored with %w", err)
	}

	return strings.TrimSuffix(strings.TrimPrefix(parsedURL.Path, "/go/api/feed/pipelines/"), "/stages.xml"), nil
}
