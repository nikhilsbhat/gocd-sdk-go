package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/pipeline_groups.json
var pipelineGroups string

func Test_client_GetPipelineGroupInfo(t *testing.T) {
	correctPipelineHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should error out while fetching all pipeline groups information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineGroups()
		assert.EqualError(t, err, "call made to get pipeline group information errored with "+
			"Get \"http://localhost:8156/go/api/admin/pipeline_groups\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all pipeline groups information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctPipelineHeader, false)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelineGroups()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all pipeline groups information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctPipelineHeader, false)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelineGroups()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to get information of all pipeline groups present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroups), http.StatusOK, correctPipelineHeader, false)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := []gocd.PipelineGroup{
			{
				Name:          "action-movies",
				PipelineCount: 2,
				Pipelines:     []gocd.Pipeline{{Name: "action-movies-auto"}, {Name: "action-movies-manual"}},
			},
			{
				Name:          "infrastructure",
				PipelineCount: 2,
				Pipelines:     []gocd.Pipeline{{Name: "gocd-prometheus-exporter"}, {Name: "helm-images"}},
			},
		}

		actual, err := client.GetPipelineGroups()
		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, actual)
	})
}

func TestGroups_Count(t *testing.T) {
	t.Run("should be able to fetch the total pipeline count", func(t *testing.T) {
		pipeGroup := gocd.Groups{
			{
				Name:          "action-movies",
				PipelineCount: 2,
				Pipelines:     []gocd.Pipeline{{Name: "action-movies-auto"}, {Name: "action-movies-manual"}},
			},
			{
				Name:          "infrastructure",
				PipelineCount: 2,
				Pipelines:     []gocd.Pipeline{{Name: "gocd-prometheus-exporter"}, {Name: "helm-images"}},
			},
		}

		acutal := pipeGroup.Count()
		assert.Equal(t, 4, acutal)
	})
}

func Test_client_DeletePipelineGroup(t *testing.T) {
	t.Run("should be able to delete pipeline group successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne}, false)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeletePipelineGroup("pipeline_group_1")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting pipeline group due to wrong headers passed", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeletePipelineGroup("pipeline_group_1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting pipeline group due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		err := client.DeletePipelineGroup("pipeline_group_1")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting pipeline group as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeletePipelineGroup("pipeline_group_1")
		assert.EqualError(t, err, "call made to delete pipeline group errored with Delete "+
			"\"http://localhost:8156/go/api/admin/pipeline_groups/pipeline_group_1\": dial tcp [::1]:8156: connect: connection refused")
	})
}
