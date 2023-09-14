package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/vsm.json
var pipelineVSM string

func Test_client_GetPipelineVSM(t *testing.T) {
	t.Run("Should be able to fetch the VSM for a selected instance of a pipeline", func(t *testing.T) {
		server := mockServer([]byte(pipelineVSM), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.VSM{
			Pipeline: "helm-images",
			Level: []gocd.PipelineLevels{
				{Nodes: []gocd.PipelineNode{{
					Parents: []string{}, Dependents: []string{"helm-images"},
					Name: "https://github.com/nikhilsbhat/helm-images", ID: "1982acfa1edbe518d3d4b866c722cd7a658b6b6cb2c1d667e5ce9829959ca491",
				}}},
				{Nodes: []gocd.PipelineNode{{
					Parents:    []string{"1982acfa1edbe518d3d4b866c722cd7a658b6b6cb2c1d667e5ce9829959ca491"},
					Dependents: []string{"api-performance-test", "deploy-helm-images-dev"}, Name: "helm-images", ID: "helm-images",
				}}},
				{Nodes: []gocd.PipelineNode{
					{
						Parents: []string{"helm-images"}, Dependents: []string{}, Name: "api-performance-test", ID: "api-performance-test",
					},
					{Parents: []string{"helm-images"}, Dependents: []string{"helm-images-tests"}, Name: "deploy-helm-images-dev", ID: "deploy-helm-images-dev"},
				}},
				{Nodes: []gocd.PipelineNode{
					{
						Parents: []string{"Deploy_HELM_IMAGES_Master"}, Dependents: []string{}, Name: "Deploy_HELM_IMAGES_Master", ID: "Deploy_HELM_IMAGES_Master",
					},
					{Parents: []string{"Deploy_HELM_IMAGES_Master"}, Dependents: []string{}, Name: "Deploy_HELM_IMAGES_Master", ID: "Deploy_HELM_IMAGES_Master"},
					{Parents: []string{"Deploy_HELM_IMAGES_Master"}, Dependents: []string{}, Name: "Deploy_HELM_IMAGES_PKS_DR_Master", ID: "Deploy_HELM_IMAGES_PKS_DR_Master"},
					{Parents: []string{"Deploy_HELM_IMAGES_Master"}, Dependents: []string{}, Name: "Deploy_HELM_IMAGES_DR_Master", ID: "Deploy_HELM_IMAGES_DR_Master"},
				}},
			},
		}

		response, err := client.GetPipelineVSM("helm-images", "20")
		assert.NoError(t, err)
		assert.Equal(t, expected, response)
	})

	t.Run("Should error out wile fetching the VSM for a selected instance of a pipeline as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`pipelineVSM`), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.GetPipelineVSM("helm-images", "20")
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, gocd.VSM{}, response)
	})

	t.Run("Should error out wile fetching the VSM for a selected instance of a pipeline as server returned non ok status code", func(t *testing.T) {
		server := mockServer([]byte(`pipelineVSM`), http.StatusInternalServerError, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.GetPipelineVSM("helm-images", "20")
		assert.EqualError(t, err, "got 500 from GoCD while making GET call for "+server.URL+
			"/pipelines/value_stream_map/helm-images/20.json\nwith BODY:pipelineVSM")
		assert.Equal(t, gocd.VSM{}, response)
	})

	t.Run("Should error out wile fetching the VSM for a selected instance of a pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineVSM("helm-images", "20")
		assert.EqualError(t, err, "call made to get vsm information for pipeline 'helm-images' of instance '20' errored with: "+
			"Get \"http://localhost:8156/go/pipelines/value_stream_map/helm-images/20.json\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.VSM{}, actual)
	})
}
