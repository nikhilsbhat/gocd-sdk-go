package main

import (
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
	resp, err := newClient.httpClient.R().SetResult(&agentsConf).Get(GoCdAgentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("call made to get agents information errored with: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, apiWithCodeError(resp.StatusCode())
	}

	return agentsConf.Config.Config, nil
}

func (conf *client) GetAgentJobRunHistory(agents []string) ([]AgentJobHistory, error) {
	newClient := &client{}
	if err := copier.CopyWithOption(newClient, conf, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	newClient.httpClient.SetHeaders(map[string]string{
		"Accept": GoCdHeaderVersionOne,
	})
	newClient.httpClient.SetQueryParam("sort_order", "DESC")

	jobHistory := make([]AgentJobHistory, 0)
	for _, agent := range agents {
		var jobHistoryConf AgentJobHistory
		resp, err := newClient.httpClient.R().SetResult(&jobHistoryConf).Get(fmt.Sprintf(GoCdJobRunHistoryEndpoint, agent))
		if err != nil {
			return nil, fmt.Errorf("call made to get agent job run history errored with %w", err)
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, apiWithCodeError(resp.StatusCode())
		}
		jobHistory = append(jobHistory, jobHistoryConf)
	}

	return jobHistory, nil
}
