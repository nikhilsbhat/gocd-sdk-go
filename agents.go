package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/jinzhu/copier"
)

// GetAgentsInfo implements method that fetches the details of all the agents present in GoCD server.
func (conf *client) GetAgentsInfo() ([]Agent, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	var agentsConf AgentsConfig
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		Get(AgentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get agents information errored with: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, APIWithCodeError(resp.StatusCode())
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

	var jobHistoryConf AgentJobHistory
	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionOne,
		}).
		SetQueryParam("sort_order", "DESC").
		Get(fmt.Sprintf(JobRunHistoryEndpoint, agentID))
	if err != nil {
		return AgentJobHistory{}, fmt.Errorf("call made to get agent job run history errored with %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return AgentJobHistory{}, APIWithCodeError(resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &jobHistoryConf); err != nil {
		return AgentJobHistory{}, ResponseReadError(err.Error())
	}

	return jobHistoryConf, nil
}

// UpdateAgent updates specific agent with updated configuration passed.
func (conf *client) UpdateAgent(agentID string, agent Agent) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionSeven,
			"Content-Type": ContentJSON,
		}).
		SetBody(agent).
		Patch(filepath.Join(AgentsEndpoint, agentID))
	if err != nil {
		return fmt.Errorf("call made to update %s agent information errored with: %w", agent.Name, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// UpdateAgentBulk will bulk update the specified agents with updated configurations.
func (conf *client) UpdateAgentBulk(agent Agent) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":       HeaderVersionSeven,
			"Content-Type": ContentJSON,
		}).
		SetBody(agent).
		Patch(AgentsEndpoint)
	if err != nil {
		return fmt.Errorf("call made to bulk update %v agents information errored with: %w", agent.UUIDS, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return nil
}

// DeleteAgent deletes the specified agent.
func (conf *client) DeleteAgent(agentID string) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).Delete(filepath.Join(AgentsEndpoint, agentID))
	if err != nil {
		return "", fmt.Errorf("call made delete agent %s errored with: %w", agentID, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return resp.String(), nil
}

// DeleteAgentBulk bulk deletes the specified agents.
func (conf *client) DeleteAgentBulk(agent Agent) (string, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return "", err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept": HeaderVersionSeven,
		}).
		SetBody(agent).
		Delete(AgentsEndpoint)
	if err != nil {
		return "", fmt.Errorf("call made delete agents %s errored with: %w", agent.UUIDS, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", APIErrorWithBody(resp.String(), resp.StatusCode())
	}

	return resp.String(), nil
}
