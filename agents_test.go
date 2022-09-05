package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/agents.json
var agentsJson string

//go:embed internal/fixtures/agent_run_history.json
var agentRunHistoryJson string

func Test_client_GetAgentsInfo(t *testing.T) {
	t.Run("should error out as call made to server while fetching agents", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetAgentsInfo()
		assert.EqualError(t, err, "call made to get agents information errored with: Get \"http://localhost:8153/go/api/agents\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching agents information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("agentsJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetAgentsInfo()
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching agents information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"_embedded": {"agents": [{`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetAgentsInfo()
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch the agents information from GoCD server", func(t *testing.T) {
		server := mockServer([]byte(agentsJson), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := []gocd.Agent{
			{
				IPAddress: "10.12.20.47", Name: "agent01.example.com", ID: "adb9540a-b954-4571-9d9b-2f330739d4da", Version: "20.5.0", CurrentState: "Idle", OS: "Mac OS X", ConfigState: "Enabled", Sandbox: "/Users/ketanpadegaonkar/projects/gocd/gocd/agent", DiskSpaceAvailable: 8.4983328768e+10,
			},
		}

		actual, err := client.GetAgentsInfo()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetAgentJobRunHistory1(t *testing.T) {
	agentID := "adb9540a-b954-4571-9d9b-2f330739d4da"

	t.Run("should error out as call made to server while fetching job run", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetAgentJobRunHistory(agentID)
		assert.EqualError(t, err, "call made to get agent job run history errored with Get \"http://localhost:8153/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da/job_run_history?sort_order=DESC\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Equal(t, gocd.AgentJobHistory{}, actual)
	})

	t.Run("should error out while fetching job run history as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("agentRunHistoryJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetAgentJobRunHistory(agentID)
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.AgentJobHistory{}, actual)
	})

	t.Run("should error out while fetching agent job run history as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"_embedded": {"agents": [{`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetAgentJobRunHistory(agentID)
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, gocd.AgentJobHistory{}, actual)
	})

	t.Run("should be able to fetch the agent job run history", func(t *testing.T) {
		server := mockServer([]byte(agentRunHistoryJson), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.AgentJobHistory{
			Jobs: []gocd.JobRunHistory{
				{
					Name:            "build-windows-PR",
					JobName:         "jasmine",
					StageName:       "build-non-server",
					StageCounter:    1,
					PipelineCounter: 5282,
					Result:          "Unknown",
				},
			},
			Pagination: gocd.Pagination{
				PageSize: 50,
				Offset:   812,
				Total:    813,
			},
		}

		actual, err := client.GetAgentJobRunHistory(agentID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
