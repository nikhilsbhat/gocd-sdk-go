package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// GetAgentsInfo implements method that fetches the details of all the agents present in GoCD server.
func (conf *client) GetAgentsInfo() ([]Agent, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionSeven,
	})

	var agentsConf AgentsConfig
	resp, err := newClient.httpClient.R().Get(GoCdAgentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get agents information errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, ApiWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &agentsConf); err != nil {
		return nil, ResponseReadError(err.Error())
	}
	return agentsConf.Config.Config, nil
}

// GetAgentJobRunHistory implements method that fetches job run history from selected agents.
func (conf *client) GetAgentJobRunHistory(agentID string) (AgentJobHistory, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return AgentJobHistory{}, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})
	newClient.httpClient.SetQueryParam("sort_order", "DESC")

	var jobHistoryConf AgentJobHistory
	resp, err := newClient.httpClient.R().Get(fmt.Sprintf(GoCdJobRunHistoryEndpoint, agentID))
	if err != nil {
		return AgentJobHistory{}, fmt.Errorf("call made to get agent job run history errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return AgentJobHistory{}, ApiWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &jobHistoryConf); err != nil {
		return AgentJobHistory{}, ResponseReadError(err.Error())
	}

	return jobHistoryConf, nil
}
