package gocd

import (
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
	resp, err := newClient.httpClient.R().SetResult(&groupConf).Get(GoCdPipelineGroupEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get pipeline group information errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}

	updatedGroupConf := make([]PipelineGroup, 0)
	for _, group := range groupConf.PipelineGroups.PipelineGroups {
		updatedGroupConf = append(updatedGroupConf, PipelineGroup{
			Name:          group.Name,
			PipelineCount: len(group.Pipelines),
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
	resp, err := newClient.httpClient.R().SetResult(&pipelinesInfo).Get(GoCdAPIFeedPipelineEndpoint)
	if err != nil {
		return PipelinesInfo{}, fmt.Errorf("call made to get pipelines errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return PipelinesInfo{}, apiWithCodeError(resp.StatusCode())
	}

	return pipelinesInfo, nil
}

// GetPipelineState fetches status of selected pipelines.
func (conf *client) GetPipelineState(pipelines []string) ([]PipelineState, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})

	pipelinesStatus := make([]PipelineState, 0)
	for _, pipeline := range pipelines {
		var pipelineState PipelineState
		resp, err := newClient.httpClient.R().SetResult(&pipelineState).Get(fmt.Sprintf(GoCdPipelineStatus, pipeline))
		if err != nil {
			return nil, fmt.Errorf("call made to get pipeline state errored with %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, apiWithCodeError(resp.StatusCode())
		}

		pipelineState.Name = pipeline
		pipelinesStatus = append(pipelinesStatus, pipelineState)
	}

	return pipelinesStatus, nil
}

// Count return the number of pipelines present.
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
