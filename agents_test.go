package gocd_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/fixtures/agent.json
	agentJSON string

	//go:embed internal/fixtures/agents.json
	agentsJSON string

	//go:embed internal/fixtures/agent_run_history.json
	agentRunHistoryJSON string

	//go:embed internal/fixtures/agents_update.json
	agentUpdateJSON string

	//go:embed internal/fixtures/agents_update_bulk.json
	agentsUpdateBulkJSON string
)

func Test_client_GetAgentsInfo(t *testing.T) {
	correctAgentsHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}

	t.Run("should error out as call made to server while fetching agents", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetAgents()
		require.EqualError(t, err, "call made to get agents information errored with: "+
			"Get \"http://localhost:8156/go/api/agents\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching agents information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("agentsJson"), http.StatusBadGateway, correctAgentsHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetAgents()
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+"/api/agents\nwith BODY:agentsJson")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching agents information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"_embedded": {"agents": [{`), http.StatusOK, correctAgentsHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetAgents()
		require.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Nil(t, actual)
	})

	t.Run("should be able to fetch the agents information from GoCD server", func(t *testing.T) {
		server := mockServer([]byte(agentsJSON), http.StatusOK, correctAgentsHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Agent{
			{
				IPAddress:          "10.12.20.47",
				Name:               "agent01.example.com",
				ID:                 "adb9540a-b954-4571-9d9b-2f330739d4da",
				Version:            "20.5.0",
				CurrentState:       "Idle",
				OS:                 "Mac OS X",
				ConfigState:        "Enabled",
				Sandbox:            "/Users/ketanpadegaonkar/projects/gocd/gocd/agent",
				DiskSpaceAvailable: 8.4983328768e+10,
				Resources:          []string{"java", "linux", "firefox"},
				BuildState:         "Idle",
			},
		}

		actual, err := client.GetAgents()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetAgent(t *testing.T) {
	correctAgentHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	agentID := "adb9540a-b954-4571-9d9b-2f330739d4da"

	t.Run("should be able to fetch a specific agent successfully", func(t *testing.T) {
		server := mockServer([]byte(agentJSON), http.StatusOK,
			correctAgentHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Agent{
			Name:               "ketanpkr.corporate.thoughtworks.com",
			IPAddress:          "10.12.20.47",
			ID:                 "adb9540a-b954-4571-9d9b-2f330739d4da",
			Version:            "20.5.0",
			CurrentState:       "Building",
			OS:                 "Mac OS X",
			ConfigState:        "Enabled",
			Sandbox:            "/Users/ketanpadegaonkar/projects/gocd/gocd/agent",
			DiskSpaceAvailable: 8.5890146304e+10,
			Resources:          []string{"java", "linux", "firefox"},
			Environments: []interface{}{
				map[string]interface{}{"name": "perf"},
				map[string]interface{}{"name": "UAT"},
			},
			BuildState: "Building",
			BuildDetails: gocd.BuildInfo{
				Pipeline: "up42",
				Stage:    "up42_stage",
				Job:      "up42_job",
			},
		}

		actual, err := client.GetAgent(agentID)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific agent present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(agentJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Agent{}

		actual, err := client.GetAgent(agentID)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific agent present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(agentJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Agent{}

		actual, err := client.GetAgent(agentID)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific agent from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("agentJSON"), http.StatusOK, correctAgentHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Agent{}

		actual, err := client.GetAgent(agentID)
		require.EqualError(t, err, "reading response body errored with: invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific agent present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.Agent{}

		actual, err := client.GetAgent(agentID)
		require.EqualError(t, err, "call made to get agent 'adb9540a-b954-4571-9d9b-2f330739d4da' information errored with: "+
			"Get \"http://localhost:8156/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetAgentJobRunHistory1(t *testing.T) {
	agentID := "adb9540a-b954-4571-9d9b-2f330739d4da" //nolint:goconst
	correctAgentsHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should error out as call made to server while fetching job run", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetAgentJobRunHistory(agentID)
		require.EqualError(t, err, "call made to get agent job run history errored with: "+
			"Get \"http://localhost:8156/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da/job_run_history?sort_order=DESC\": "+
			"dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.AgentJobHistory{}, actual)
	})

	t.Run("should error out while fetching job run history as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("agentRunHistoryJSON"), http.StatusBadGateway, correctAgentsHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetAgentJobRunHistory(agentID)
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da/job_run_history?sort_order=DESC\nwith BODY:agentRunHistoryJSON")
		assert.Equal(t, gocd.AgentJobHistory{}, actual)
	})

	t.Run("should error out while fetching agent job run history as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"_embedded": {"agents": [{`), http.StatusOK, correctAgentsHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetAgentJobRunHistory(agentID)
		require.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, gocd.AgentJobHistory{}, actual)
	})

	t.Run("should be able to fetch the agent job run history", func(t *testing.T) {
		server := mockServer([]byte(agentRunHistoryJSON), http.StatusOK, correctAgentsHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

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
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateAgent(t *testing.T) {
	agentID := "adb9540a-b954-4571-9d9b-2f330739d4da"
	correctAgentUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON}

	t.Run("should update agent with updated configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(agentUpdateJSON), http.StatusOK,
			correctAgentUpdateHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			ID:           agentID,
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    []string{"Java", "Linux"},
			Environments: []string{"Dev"},
		}

		err := client.UpdateAgent(agentUpdateInfo)
		require.NoError(t, err)
	})

	t.Run("should error out while updating agent information due to wrong headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodPatch, map[string]string{"Accept": gocd.HeaderVersionFour, "Content-Type": gocd.ContentJSON})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			ID:           agentID,
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgent(agentUpdateInfo)
		require.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while updating agent information due to missing headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodPatch, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			ID:           agentID,
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgent(agentUpdateInfo)
		require.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while updating agent information as malformed data sent to server", func(t *testing.T) {
		server := agentMockServer([]byte("agentsJSON"), http.MethodPatch, correctAgentUpdateHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			ID:           agentID,
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgent(agentUpdateInfo)
		require.EqualError(t, err, "got 500 from GoCD while making PATCH call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:json: cannot unmarshal string into Go value of type gocd.Agent")
	})

	t.Run("should error out while updating agent information as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		agentUpdateInfo := gocd.Agent{
			ID:           agentID,
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgent(agentUpdateInfo)
		require.EqualError(t, err, "call made to update agent02.example.com agent information errored with: "+
			"Patch \"http://localhost:8156/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_UpdateAgentBulk(t *testing.T) {
	correctAgentUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to bulk update the specified agents with updated configurations", func(t *testing.T) {
		var agentInfo gocd.Agent
		err := json.Unmarshal([]byte(agentsUpdateBulkJSON), &agentInfo)
		require.NoError(t, err)

		server := agentMockServer(agentInfo, http.MethodPatch, correctAgentUpdateHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err = client.UpdateAgentBulk(agentUpdateInfo)
		require.NoError(t, err)
	})

	t.Run("should error out while bulk updating agents information due to wrong headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodPatch, map[string]string{"Accept": gocd.HeaderVersionFour, "Content-Type": gocd.ContentJSON})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgentBulk(agentUpdateInfo)
		require.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/agents\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while bulk updating agents information due to missing headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodPatch, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgentBulk(agentUpdateInfo)
		require.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/agents\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while bulk updating agents information as malformed data sent to server", func(t *testing.T) {
		server := agentMockServer([]byte("agentsJSON"), http.MethodPatch, correctAgentUpdateHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agentUpdateInfo := gocd.Agent{
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgentBulk(agentUpdateInfo)
		require.EqualError(t, err, "got 500 from GoCD while making PATCH call for "+server.URL+
			"/api/agents\nwith BODY:json: cannot unmarshal string into Go value of type gocd.Agent")
	})

	t.Run("should error out while bulk updating agents information as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		agentUpdateInfo := gocd.Agent{
			Name:         "agent02.example.com",
			ConfigState:  "Enabled",
			Resources:    nil,
			Environments: nil,
		}

		err := client.UpdateAgentBulk(agentUpdateInfo)
		require.EqualError(t, err, "call made to bulk update [] agents information errored with: "+
			"Patch \"http://localhost:8156/go/api/agents\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_DeleteAgent(t *testing.T) {
	correctAgentUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	agentID := "adb9540a-b954-4571-9d9b-2f330739d4da"

	t.Run("should be able to delete the specified agent successfully", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodDelete, correctAgentUpdateHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.DeleteAgent(agentID)
		require.NoError(t, err)
		assert.JSONEq(t, `{"message": "Deleted 1 agent(s)."}`, actual)
	})

	t.Run("should error out while deleting agent due to wrong headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionFour})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.DeleteAgent(agentID)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while deleting agent due to missing headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodDelete, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.DeleteAgent(agentID)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while deleting agent as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.DeleteAgent(agentID)
		require.EqualError(t, err, "call made to delete agent adb9540a-b954-4571-9d9b-2f330739d4da errored with: "+
			"Delete \"http://localhost:8156/go/api/agents/adb9540a-b954-4571-9d9b-2f330739d4da\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})
}

func Test_client_DeleteAgentBulk(t *testing.T) {
	correctAgentUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionSeven}
	agent := gocd.Agent{
		UUIDS: []string{"adb9540a-b954-4571-9d9b-2f330739d4da", "adb9540a-5hfh-6453-9d9b-2f37467739d4da"},
	}

	t.Run("should be able bulk delete the specified agents successfully", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodDelete, correctAgentUpdateHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.DeleteAgentBulk(agent)
		require.NoError(t, err)
		assert.JSONEq(t, `{"message": "Deleted 1 agent(s)."}`, actual)
	})

	t.Run("should error out while bulk deleting agents due to wrong headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodDelete, map[string]string{"Accept": gocd.HeaderVersionFour})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.DeleteAgentBulk(agent)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/agents\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while bulk deleting agents due to missing headers", func(t *testing.T) {
		server := agentMockServer(nil, http.MethodDelete, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.DeleteAgentBulk(agent)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/agents\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while bulk deleting agent as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.DeleteAgentBulk(agent)
		require.EqualError(t, err, "call made to delete agents [adb9540a-b954-4571-9d9b-2f330739d4da adb9540a-5hfh-6453-9d9b-2f37467739d4da] errored with: "+
			"Delete \"http://localhost:8156/go/api/agents\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})
}

func Test_client_AgentKillTask(t *testing.T) {
	correctTaskKillHeader := map[string]string{
		"Accept":           gocd.HeaderVersionSeven,
		gocd.HeaderConfirm: "true",
	}

	t.Run("should be able to cancel the tasks running on an agent successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroupUpdate), http.StatusOK, correctTaskKillHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agent := gocd.Agent{ID: "adb9540a-5hfh-6453-9d9b-2f37467739d4da"}

		err := client.AgentKillTask(agent)
		require.NoError(t, err)
	})

	t.Run("should error out while canceling the tasks running on an agent due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroupUpdate), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo, gocd.HeaderConfirm: "true"},
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agent := gocd.Agent{ID: "adb9540a-5hfh-6453-9d9b-2f37467739d4da"}

		err := client.AgentKillTask(agent)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/agents/adb9540a-5hfh-6453-9d9b-2f37467739d4da/kill_running_tasks\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while canceling the tasks running on an agent due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroupUpdate), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo},
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		agent := gocd.Agent{ID: "adb9540a-5hfh-6453-9d9b-2f37467739d4da"}

		err := client.AgentKillTask(agent)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/agents/adb9540a-5hfh-6453-9d9b-2f37467739d4da/kill_running_tasks\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while canceling the tasks running on an agent since server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		agent := gocd.Agent{ID: "adb9540a-5hfh-6453-9d9b-2f37467739d4da"}

		err := client.AgentKillTask(agent)
		require.EqualError(t, err, "call made to kill tasks from agent adb9540a-5hfh-6453-9d9b-2f37467739d4da errored with: "+
			"Post \"http://localhost:8156/go/api/agents/adb9540a-5hfh-6453-9d9b-2f37467739d4da/kill_running_tasks\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func agentMockServer(request interface{}, method string, header map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if header == nil {
			writer.WriteHeader(http.StatusNotFound)

			if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
				log.Fatalln(err)
			}

			return
		}

		for key, value := range header {
			if req.Header.Get(key) != value {
				writer.WriteHeader(http.StatusNotFound)

				if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
					log.Fatalln(err)
				}

				return
			}
		}

		if method == http.MethodDelete {
			writer.WriteHeader(http.StatusOK)

			if _, err := writer.Write([]byte(`{"message": "Deleted 1 agent(s)."}`)); err != nil {
				log.Fatalln(err)
			}
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			if _, err = writer.Write([]byte(fmt.Sprintf("%s %s", string(requestByte), err.Error()))); err != nil {
				log.Fatalln(err)
			}

			return
		}

		var gocdAgent gocd.Agent

		if err = json.Unmarshal(requestByte, &gocdAgent); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			if _, err = writer.Write([]byte(err.Error())); err != nil {
				log.Fatalln(err)
			}

			return
		}

		writer.WriteHeader(http.StatusOK)

		if _, err = writer.Write([]byte("")); err != nil {
			log.Fatalln(err)
		}
	}))
}
