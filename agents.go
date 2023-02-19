package gocd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"

	"github.com/jinzhu/copier"
)

// GetAgents implements method that fetches the details of all the agents present in GoCD server.
func (conf *client) GetAgents() ([]Agent, error) {
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
		return nil, &errors.APIError{Err: err, Message: "get agents information"}
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &agentsConf); err != nil {
		return nil, &errors.MarshalError{Err: err}
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
		return AgentJobHistory{}, &errors.APIError{Err: err, Message: "get agent job run history"}
	}
	if resp.StatusCode() != http.StatusOK {
		return AgentJobHistory{}, &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	if err = json.Unmarshal(resp.Body(), &jobHistoryConf); err != nil {
		return AgentJobHistory{}, &errors.MarshalError{Err: err}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("update %s agent information", agent.Name)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		return &errors.APIError{Err: err, Message: fmt.Sprintf("bulk update %v agents information", agent.UUIDS)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		}).
		Delete(filepath.Join(AgentsEndpoint, agentID))
	if err != nil {
		return "", &errors.APIError{Err: err, Message: fmt.Sprintf("delete agent %s", agentID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
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
		return "", &errors.APIError{Err: err, Message: fmt.Sprintf("delete agents %s", agent.UUIDS)}
	}

	if resp.StatusCode() != http.StatusOK {
		return "", &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return resp.String(), nil
}

// AgentKillTask will kill running tasks from an selected agent.
func (conf *client) AgentKillTask(agent Agent) error {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}

	resp, err := newClient.httpClient.R().
		SetHeaders(map[string]string{
			"Accept":      HeaderVersionSeven,
			HeaderConfirm: "true",
		}).
		Post(filepath.Join(AgentsEndpoint, agent.ID, "kill_running_tasks"))
	if err != nil {
		return &errors.APIError{Err: err, Message: fmt.Sprintf("kill tasks from agent %s", agent.ID)}
	}

	if resp.StatusCode() != http.StatusOK {
		return &errors.NonOkError{Code: resp.StatusCode(), Response: resp}
	}

	return nil
}
