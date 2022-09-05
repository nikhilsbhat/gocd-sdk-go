package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/pipelines.xml
	pipelines string
	//go:embed internal/fixtures/pipeline_state.json
	pipelineState string
	//go:embed internal/fixtures/pipeline_groups.json
	pipelineGroups string
)

func Test_client_GetPipelines(t *testing.T) {
	t.Run("should error out while fetching pipelines from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, "call made to get pipelines errored with Get \"http://localhost:8153/go/api/feed/pipelines.xml\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, "reading response body errored with: EOF")
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should be able to fetch the pipelines present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelines), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelines()
		assert.NoError(t, err)
		assert.Equal(t, 3, len(actual.Pipeline))
	})
}

func Test_client_getPipelineName(t *testing.T) {
	t.Run("should be able to fetch pipeline name from the href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "animation-and-action-movies", name)
	})
	t.Run("should fetch malformed pipeline name as malformed/invalid (prefix) href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8153/go/api/feed/pipelinesss/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "/go/api/feed/pipelinesss/animation-and-action-movies", name)
	})
	t.Run("should fetch malformed pipeline name as malformed/invalid (suffix) href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/test/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "animation-and-action-movies/test", name)
	})
	t.Run("should fail while fetching pipeline name from malformed/invalid href", func(t *testing.T) {
		pipelineLink := "://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.EqualError(t, err, "parsing URL errored with parse \"://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml\": missing protocol scheme")
		assert.Equal(t, "", name)
	})
}

func Test_client_GetPipelineStatus(t *testing.T) {
	pipeline := "action-movies-manual"

	t.Run("should error out while fetching pipeline statuses information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, "call made to get pipeline state errored with Get \"http://localhost:8153/go/api/pipelines/action-movies-manual/status\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should error out while fetching pipeline statuses information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should error out while fetching pipeline statuses information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should be able to get pipeline statuses of all pipeline present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelineState), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		expected := gocd.PipelineState{
			Name:        pipeline,
			Paused:      true,
			Locked:      false,
			Schedulable: false,
			PausedBy:    "admin",
			PausedCause: "Reason for pausing this pipeline goes here",
		}

		actual, err := client.GetPipelineState(pipeline)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetPipelineGroupInfo(t *testing.T) {
	t.Run("should error out while fetching all pipeline groups information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineGroupInfo()
		assert.EqualError(t, err, "call made to get pipeline group information errored with Get \"http://localhost:8153/go/api/admin/pipeline_groups\": dial tcp 127.0.0.1:8153: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all pipeline groups information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJson"), http.StatusBadGateway)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelineGroupInfo()
		assert.EqualError(t, err, gocd.ApiWithCodeError(http.StatusBadGateway).Error())
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all pipeline groups information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetPipelineGroupInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to get information of all pipeline groups present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroups), http.StatusOK)
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

		actual, err := client.GetPipelineGroupInfo()
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
