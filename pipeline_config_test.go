package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/pipeline_confg.json
var pipelineConfigJSON string

//nolint:funlen
func Test_client_GetPipelineConfig(t *testing.T) {
	correctPipelineConfigHeader := map[string]string{"Accept": gocd.HeaderVersionEleven}
	t.Run("should be able to fetch pipeline configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			correctPipelineConfigHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Config: map[string]interface{}{
				"environment_variables": []interface{}{},
				"group":                 "new_group",
				"label_template":        "${COUNT}",
				"lock_behavior":         "lockOnFailure",
				"materials": []interface{}{map[string]interface{}{
					"attributes": map[string]interface{}{
						"auto_update":      true,
						"branch":           "master",
						"destination":      "dest",
						"filter":           interface{}(nil),
						"invert_filter":    false,
						"name":             interface{}(nil),
						"shallow_clone":    false,
						"submodule_folder": interface{}(nil),
						"url":              "git@github.com:sample_repo/example.git",
					},
					"type": "git",
				}},
				"name":       "new_pipeline",
				"parameters": []interface{}{},
				"stages": []interface{}{
					map[string]interface{}{
						"approval": map[string]interface{}{
							"authorization": map[string]interface{}{
								"roles": []interface{}{},
								"users": []interface{}{},
							},
							"type": "success",
						},
						"clean_working_directory": false,
						"environment_variables":   []interface{}{},
						"fetch_materials":         true,
						"jobs": []interface{}{map[string]interface{}{
							"artifacts": []interface{}{map[string]interface{}{
								"artifact_id": "docker-image",
								"configuration": []interface{}{
									map[string]interface{}{
										"key":   "Image",
										"value": "gocd/gocd-server",
									},
									map[string]interface{}{
										"key":   "Tag",
										"value": "v${GO_PIPELINE_LABEL}",
									},
								},
								"store_id": "dockerhub",
								"type":     "external",
							}},
							"environment_variables": []interface{}{},
							"name":                  "defaultJob",
							"resources":             []interface{}{},
							"run_instance_count":    interface{}(nil),
							"tabs":                  []interface{}{},
							"tasks": []interface{}{map[string]interface{}{
								"attributes": map[string]interface{}{
									"args":    "",
									"command": "ls",
									"run_if":  []interface{}{"passed"},
								},
								"type": "exec",
							}},
							"timeout": interface{}(nil),
						}},
						"name":                    "defaultStage",
						"never_cleanup_artifacts": false,
					},
					map[string]interface{}{
						"approval": map[string]interface{}{
							"authorization": map[string]interface{}{
								"roles": []interface{}{},
								"users": []interface{}{},
							},
							"type": "success",
						},
						"clean_working_directory": false,
						"environment_variables":   []interface{}{},
						"fetch_materials":         true,
						"jobs": []interface{}{map[string]interface{}{
							"artifacts":             []interface{}{},
							"environment_variables": []interface{}{},
							"name":                  "j2",
							"resources":             []interface{}{},
							"run_instance_count":    interface{}(nil),
							"tabs":                  []interface{}{},
							"tasks": []interface{}{map[string]interface{}{
								"attributes": map[string]interface{}{
									"artifact_id":     "docker-image",
									"artifact_origin": "external",
									"job":             "defaultJob",
									"pipeline":        "",
									"run_if":          []interface{}{},
									"stage":           "defaultStage",
								},
								"type": "fetch",
							}},
							"timeout": interface{}(nil),
						}},
						"name":                    "s2",
						"never_cleanup_artifacts": false,
					},
				},
				"template":      interface{}(nil),
				"timer":         interface{}(nil),
				"tracking_tool": interface{}(nil),
			},
		}

		actual, err := client.GetPipelineConfig("new_pipeline")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline configuration present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{}

		actual, err := client.GetPipelineConfig("new_pipeline")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/pipelines/new_pipeline\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline configuration present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{}

		actual, err := client.GetPipelineConfig("new_pipeline")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/pipelines/new_pipeline\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline configuration from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pipelineConfigJSON"), http.StatusOK, correctPipelineConfigHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{}

		actual, err := client.GetPipelineConfig("new_pipeline")
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline configuration present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.PipelineConfig{}

		actual, err := client.GetPipelineConfig("new_pipeline")
		assert.EqualError(t, err, "call made to get pipeline config 'new_pipeline' errored with: Get "+
			"\"http://localhost:8156/go/api/admin/pipelines/new_pipeline\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeletePipeline(t *testing.T) {
	correctPipelineConfigHeader := map[string]string{"Accept": gocd.HeaderVersionEleven}
	t.Run("should be able to delete pipeline successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctPipelineConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePipeline("pipeline_group_1")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting pipeline due to wrong headers passed", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePipeline("pipeline_group_1")
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/pipelines/pipeline_group_1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting pipeline due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePipeline("pipeline_group_1")
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/pipelines/pipeline_group_1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeletePipeline("pipeline_group_1")
		assert.EqualError(t, err, "call made to delete pipeline config 'pipeline_group_1' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/pipelines/pipeline_group_1\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_ExtractTemplatePipeline(t *testing.T) {
	correctPipelineConfigHeader := map[string]string{"Accept": gocd.HeaderVersionEleven}
	t.Run("should be able to extract template from pipeline successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK, correctPipelineConfigHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.ExtractTemplatePipeline("pipeline_group_1", "my_template")
		assert.NotNil(t, response)
		assert.NoError(t, err)
	})

	t.Run("should error out while extracting template from pipeline due to wrong headers passed", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.ExtractTemplatePipeline("pipeline_group_1", "my_template")
		assert.NotNil(t, response)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/pipelines/pipeline_group_1/extract_to_template\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while extracting template from pipeline due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.ExtractTemplatePipeline("pipeline_group_1", "my_template")
		assert.NotNil(t, response)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/pipelines/pipeline_group_1/extract_to_template\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while extracting template from pipeline as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pipelineConfigJSON"), http.StatusOK, correctPipelineConfigHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		response, err := client.ExtractTemplatePipeline("pipeline_group_1", "my_template")
		assert.NotNil(t, response)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
	})

	t.Run("should error out while extracting template from pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		response, err := client.ExtractTemplatePipeline("pipeline_group_1", "my_template")
		assert.NotNil(t, response)
		assert.EqualError(t, err, "call made to extracting template from pipeline 'pipeline_group_1' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/pipelines/pipeline_group_1/extract_to_template\": dial tcp [::1]:8156: connect: connection refused")
	})
}

//nolint:funlen
func Test_client_UpdatePipelineConfig(t *testing.T) {
	correctPipelineConfigHeader := map[string]string{"Accept": gocd.HeaderVersionEleven}
	t.Run("should be able to update pipeline configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			correctPipelineConfigHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Config: map[string]interface{}{
				"environment_variables": []interface{}{},
				"group":                 "new_group",
				"label_template":        "${COUNT}",
				"lock_behavior":         "lockOnFailure",
				"materials": []interface{}{map[string]interface{}{
					"attributes": map[string]interface{}{
						"auto_update":      true,
						"branch":           "master",
						"destination":      "dest",
						"filter":           interface{}(nil),
						"invert_filter":    false,
						"name":             interface{}(nil),
						"shallow_clone":    false,
						"submodule_folder": interface{}(nil),
						"url":              "git@github.com:sample_repo/example.git",
					},
					"type": "git",
				}},
				"name":       "new_pipeline",
				"parameters": []interface{}{},
				"stages": []interface{}{
					map[string]interface{}{
						"approval": map[string]interface{}{
							"authorization": map[string]interface{}{
								"roles": []interface{}{},
								"users": []interface{}{},
							},
							"type": "success",
						},
						"clean_working_directory": false,
						"environment_variables":   []interface{}{},
						"fetch_materials":         true,
						"jobs": []interface{}{map[string]interface{}{
							"artifacts": []interface{}{map[string]interface{}{
								"artifact_id": "docker-image",
								"configuration": []interface{}{
									map[string]interface{}{
										"key":   "Image",
										"value": "gocd/gocd-server",
									},
									map[string]interface{}{
										"key":   "Tag",
										"value": "v${GO_PIPELINE_LABEL}",
									},
								},
								"store_id": "dockerhub",
								"type":     "external",
							}},
							"environment_variables": []interface{}{},
							"name":                  "defaultJob",
							"resources":             []interface{}{},
							"run_instance_count":    interface{}(nil),
							"tabs":                  []interface{}{},
							"tasks": []interface{}{map[string]interface{}{
								"attributes": map[string]interface{}{
									"args":    "",
									"command": "ls",
									"run_if":  []interface{}{"passed"},
								},
								"type": "exec",
							}},
							"timeout": interface{}(nil),
						}},
						"name":                    "defaultStage",
						"never_cleanup_artifacts": false,
					},
					map[string]interface{}{
						"approval": map[string]interface{}{
							"authorization": map[string]interface{}{
								"roles": []interface{}{},
								"users": []interface{}{},
							},
							"type": "success",
						},
						"clean_working_directory": false,
						"environment_variables":   []interface{}{},
						"fetch_materials":         true,
						"jobs": []interface{}{map[string]interface{}{
							"artifacts":             []interface{}{},
							"environment_variables": []interface{}{},
							"name":                  "j2",
							"resources":             []interface{}{},
							"run_instance_count":    interface{}(nil),
							"tabs":                  []interface{}{},
							"tasks": []interface{}{map[string]interface{}{
								"attributes": map[string]interface{}{
									"artifact_id":     "docker-image",
									"artifact_origin": "external",
									"job":             "defaultJob",
									"pipeline":        "",
									"run_if":          []interface{}{},
									"stage":           "defaultStage",
								},
								"type": "fetch",
							}},
							"timeout": interface{}(nil),
						}},
						"name":                    "s2",
						"never_cleanup_artifacts": false,
					},
				},
				"template":      interface{}(nil),
				"timer":         interface{}(nil),
				"tracking_tool": interface{}(nil),
			},
		}

		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		actual, err := client.UpdatePipelineConfig(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline configuration present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{}
		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		actual, err := client.UpdatePipelineConfig(input)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/pipelines/new_pipeline\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline configuration present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{}
		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		actual, err := client.UpdatePipelineConfig(input)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/pipelines/new_pipeline\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline configuration from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pipelineConfigJSON"), http.StatusOK, correctPipelineConfigHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineConfig{}
		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		actual, err := client.UpdatePipelineConfig(input)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline configuration present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.PipelineConfig{}
		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		actual, err := client.UpdatePipelineConfig(input)
		assert.EqualError(t, err, "call made to update pipeline config 'new_pipeline' errored with: Put "+
			"\"http://localhost:8156/go/api/admin/pipelines/new_pipeline\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreatePipeline(t *testing.T) {
	correctPipelineConfigHeader := map[string]string{"Accept": gocd.HeaderVersionEleven}
	t.Run("should be able to create pipeline configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			correctPipelineConfigHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Config:        map[string]interface{}{"name": "new_pipeline"},
			ETAG:          "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
			PausePipeline: true,
		}

		err := client.CreatePipeline(input)
		assert.NoError(t, err)
	})

	t.Run("should error out while creating pipeline configuration present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		err := client.CreatePipeline(input)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/pipelines/new_pipeline\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating pipeline configuration present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		err := client.CreatePipeline(input)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/pipelines/new_pipeline\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating pipeline configuration present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		input := gocd.PipelineConfig{
			Config: map[string]interface{}{"name": "new_pipeline"},
			ETAG:   "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		err := client.CreatePipeline(input)
		assert.EqualError(t, err, "call made to create pipeline config 'new_pipeline' errored with: Post "+
			"\"http://localhost:8156/go/api/admin/pipelines/new_pipeline\": dial tcp [::1]:8156: connect: connection refused")
	})
}
