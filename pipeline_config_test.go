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
			Group:         "new_group",
			LabelTemplate: "${COUNT}",
			LockBehavior:  "lockOnFailure",
			Name:          "new_pipeline",
			Template:      "",
			Origin: gocd.PipelineOrigin{
				Type: "config_repo", ID: "sample_config",
			},
			Parameters:           []gocd.PipelineEnvironmentVariables{},
			EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
			Materials: []gocd.Material{
				{
					Type: "git", Fingerprint: "",
					Attributes: gocd.Attribute{
						URL:                 "git@github.com:sample_repo/example.git",
						Username:            "",
						Password:            "",
						EncryptedPassword:   "",
						Branch:              "master",
						AutoUpdate:          true,
						CheckExternals:      false,
						UseTickets:          false,
						View:                "",
						Port:                "",
						ProjectPath:         "",
						Domain:              "",
						Ref:                 "",
						Name:                "",
						Stage:               "",
						Pipeline:            "",
						IgnoreForScheduling: false,
						Destination:         "dest",
						InvertFilter:        false,
						SubmoduleFolder:     "",
						ShallowClone:        false,
						Origin:              map[string]string(nil),
						Filter: struct {
							Ignore []string "json:\"ignore,omitempty\" yaml:\"ignore,omitempty\""
						}{Ignore: []string(nil)},
					},
					RepoURL: "",
					Config: gocd.MaterialConfig{
						Type:        "",
						Fingerprint: "",
						Attributes: gocd.Attribute{
							URL:                 "",
							Username:            "",
							Password:            "",
							EncryptedPassword:   "",
							Branch:              "",
							AutoUpdate:          false,
							CheckExternals:      false,
							UseTickets:          false,
							View:                "",
							Port:                "",
							ProjectPath:         "",
							Domain:              "",
							Ref:                 "",
							Name:                "",
							Stage:               "",
							Pipeline:            "",
							IgnoreForScheduling: false,
							Destination:         "",
							InvertFilter:        false,
							SubmoduleFolder:     "",
							ShallowClone:        false,
							Origin:              map[string]string(nil),
							Filter: struct {
								Ignore []string "json:\"ignore,omitempty\" yaml:\"ignore,omitempty\""
							}{Ignore: []string(nil)},
						},
					},
					CanTriggerUpdate:         false,
					MaterialUpdateInProgress: false,
					Messages:                 []map[string]string(nil),
				},
			},
			Stages: []gocd.PipelineStageConfig{
				{
					Name:                  "defaultStage",
					FetchMaterials:        true,
					CleanWorkingDirectory: false,
					NeverCleanupArtifacts: false,
					Approval: gocd.PipelineApprovalConfig{
						Type: "success", AllowOnlyOnSuccess: false,
						Authorization: gocd.AuthorizationConfig{
							Roles: []string{},
							Users: []string{},
						},
					},
					EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
					Jobs: []gocd.PipelineJobConfig{
						{
							Name:                 "defaultJob",
							RunInstanceCount:     interface{}(nil),
							ElasticProfileID:     "",
							Timeout:              0,
							EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
							Resources:            []string{},
							Tasks: []gocd.PipelineTaskConfig{
								{
									Type: "exec",
									Attributes: gocd.TaskAttributeConfig{
										IsSourceAFile:    false,
										ArtifactID:       "",
										ArtifactOrigin:   "",
										Command:          "ls",
										WorkingDirectory: "",
										Pipeline:         "",
										Stage:            "",
										Job:              "",
										Source:           "",
										Destination:      "",
										Configuration:    []gocd.PluginConfiguration(nil),
										RunIf:            []string{"passed"},
										Arguments:        []string(nil),
										PluginConfiguration: struct {
											ID      string "json:\"id,omitempty\" yaml:\"id,omitempty\""
											Version string "json:\"version,omitempty\" yaml:\"version,omitempty\""
										}{ID: "", Version: ""},
										OnCancel: struct {
											RunIf            []string "json:\"run_if,omitempty\" yaml:\"run_if,omitempty\""
											Command          string   "json:\"command,omitempty\" yaml:\"command,omitempty\""
											Arguments        []string "json:\"arguments,omitempty\" yaml:\"arguments,omitempty\""
											WorkingDirectory string   "json:\"working_directory,omitempty\" yaml:\"working_directory,omitempty\""
										}{RunIf: []string(nil), Command: "", Arguments: []string(nil), WorkingDirectory: ""},
									},
								},
							},
							Tabs: []gocd.PipelineTab{},
							Artifacts: []gocd.PipelineArtifact{
								{
									Type:          "external",
									Source:        "",
									Destination:   "",
									ArtifactID:    "",
									StoreID:       "dockerhub",
									Configuration: []map[string]string{{"key": "Image", "value": "gocd/gocd-server"}, {"key": "Tag", "value": "v${GO_PIPELINE_LABEL}"}},
								},
							},
						},
					},
				},
				{
					Name:                  "s2",
					FetchMaterials:        true,
					CleanWorkingDirectory: false,
					NeverCleanupArtifacts: false,
					Approval: gocd.PipelineApprovalConfig{
						Type:               "success",
						AllowOnlyOnSuccess: false,
						Authorization: gocd.AuthorizationConfig{
							Roles: []string{},
							Users: []string{},
						},
					},
					EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
					Jobs: []gocd.PipelineJobConfig{
						{
							Name: "j2", RunInstanceCount: interface{}(nil),
							ElasticProfileID:     "",
							Timeout:              0,
							EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
							Resources:            []string{},
							Tasks: []gocd.PipelineTaskConfig{
								{
									Type: "fetch",
									Attributes: gocd.TaskAttributeConfig{
										IsSourceAFile:    false,
										ArtifactID:       "docker-image",
										ArtifactOrigin:   "external",
										Command:          "",
										WorkingDirectory: "",
										Pipeline:         "",
										Stage:            "defaultStage",
										Job:              "defaultJob",
										Source:           "",
										Destination:      "",
										Configuration:    []gocd.PluginConfiguration(nil),
										RunIf:            []string{}, Arguments: []string(nil), PluginConfiguration: struct {
											ID      string "json:\"id,omitempty\" yaml:\"id,omitempty\""
											Version string "json:\"version,omitempty\" yaml:\"version,omitempty\""
										}{ID: "", Version: ""},
										OnCancel: struct {
											RunIf            []string "json:\"run_if,omitempty\" yaml:\"run_if,omitempty\""
											Command          string   "json:\"command,omitempty\" yaml:\"command,omitempty\""
											Arguments        []string "json:\"arguments,omitempty\" yaml:\"arguments,omitempty\""
											WorkingDirectory string   "json:\"working_directory,omitempty\" yaml:\"working_directory,omitempty\""
										}{RunIf: []string(nil), Command: "", Arguments: []string(nil), WorkingDirectory: ""},
									},
								},
							},
							Tabs:      []gocd.PipelineTab{},
							Artifacts: []gocd.PipelineArtifact{},
						},
					},
				},
			},
			TrackingTool: gocd.PipelineTracingToolConfig{
				Type: "",
				Attributes: struct {
					URLPattern string "json:\"url_pattern,omitempty\" yaml:\"url_pattern,omitempty\""
					Regex      string "json:\"regex,omitempty\" yaml:\"regex,omitempty\""
				}{URLPattern: "", Regex: ""},
			},
			Timer: gocd.PipelineTimerConfig{
				Spec:          "",
				OnlyOnChanges: false,
			},
			CreateOptions: gocd.PipelineCreateOptions{
				PausePipeline: false,
				PauseReason:   "",
			},
			Config: map[string]interface{}(nil),
			ETAG:   "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.GetPipelineConfig("new_pipeline")
		assert.NoError(t, err)
		assert.EqualValues(t, expected, actual)
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
			Group:         "new_group",
			LabelTemplate: "${COUNT}",
			LockBehavior:  "lockOnFailure",
			Name:          "new_pipeline",
			Template:      "",
			Origin: gocd.PipelineOrigin{
				Type: "config_repo",
				ID:   "sample_config",
			},
			Parameters:           []gocd.PipelineEnvironmentVariables{},
			EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
			Materials: []gocd.Material{
				{
					Type:        "git",
					Fingerprint: "",
					Attributes: gocd.Attribute{
						URL:                 "git@github.com:sample_repo/example.git",
						Username:            "",
						Password:            "",
						EncryptedPassword:   "",
						Branch:              "master",
						AutoUpdate:          true,
						CheckExternals:      false,
						UseTickets:          false,
						View:                "",
						Port:                "",
						ProjectPath:         "",
						Domain:              "",
						Ref:                 "",
						Name:                "",
						Stage:               "",
						Pipeline:            "",
						IgnoreForScheduling: false,
						Destination:         "dest",
						InvertFilter:        false,
						SubmoduleFolder:     "",
						ShallowClone:        false,
						Origin:              map[string]string(nil),
						Filter: struct {
							Ignore []string "json:\"ignore,omitempty\" yaml:\"ignore,omitempty\""
						}{Ignore: []string(nil)},
					}, RepoURL: "", Config: gocd.MaterialConfig{Type: "", Fingerprint: "", Attributes: gocd.Attribute{URL: "", Username: "", Password: "", EncryptedPassword: "", Branch: "", AutoUpdate: false, CheckExternals: false, UseTickets: false, View: "", Port: "", ProjectPath: "", Domain: "", Ref: "", Name: "", Stage: "", Pipeline: "", IgnoreForScheduling: false, Destination: "", InvertFilter: false, SubmoduleFolder: "", ShallowClone: false, Origin: map[string]string(nil), Filter: struct {
						Ignore []string "json:\"ignore,omitempty\" yaml:\"ignore,omitempty\""
					}{Ignore: []string(nil)}}},
					CanTriggerUpdate:         false,
					MaterialUpdateInProgress: false,
					Messages:                 []map[string]string(nil),
				},
			},
			Stages: []gocd.PipelineStageConfig{
				{
					Name:                  "defaultStage",
					FetchMaterials:        true,
					CleanWorkingDirectory: false,
					NeverCleanupArtifacts: false,
					Approval: gocd.PipelineApprovalConfig{
						Type:               "success",
						AllowOnlyOnSuccess: false,
						Authorization: gocd.AuthorizationConfig{
							Roles: []string{},
							Users: []string{},
						},
					},
					EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
					Jobs: []gocd.PipelineJobConfig{
						{
							Name:                 "defaultJob",
							RunInstanceCount:     interface{}(nil),
							ElasticProfileID:     "",
							Timeout:              0,
							EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
							Resources:            []string{},
							Tasks: []gocd.PipelineTaskConfig{{
								Type: "exec",
								Attributes: gocd.TaskAttributeConfig{
									IsSourceAFile:    false,
									ArtifactID:       "",
									ArtifactOrigin:   "",
									Command:          "ls",
									WorkingDirectory: "",
									Pipeline:         "",
									Stage:            "",
									Job:              "",
									Source:           "",
									Destination:      "",
									Configuration:    []gocd.PluginConfiguration(nil),
									RunIf:            []string{"passed"},
									Arguments:        []string(nil),
									PluginConfiguration: struct {
										ID      string "json:\"id,omitempty\" yaml:\"id,omitempty\""
										Version string "json:\"version,omitempty\" yaml:\"version,omitempty\""
									}{ID: "", Version: ""},
									OnCancel: struct {
										RunIf            []string "json:\"run_if,omitempty\" yaml:\"run_if,omitempty\""
										Command          string   "json:\"command,omitempty\" yaml:\"command,omitempty\""
										Arguments        []string "json:\"arguments,omitempty\" yaml:\"arguments,omitempty\""
										WorkingDirectory string   "json:\"working_directory,omitempty\" yaml:\"working_directory,omitempty\""
									}{RunIf: []string(nil), Command: "", Arguments: []string(nil), WorkingDirectory: ""},
								},
							}},
							Tabs: []gocd.PipelineTab{},
							Artifacts: []gocd.PipelineArtifact{{
								Type:          "external",
								Source:        "",
								Destination:   "",
								ArtifactID:    "",
								StoreID:       "dockerhub",
								Configuration: []map[string]string{{"key": "Image", "value": "gocd/gocd-server"}, {"key": "Tag", "value": "v${GO_PIPELINE_LABEL}"}},
							}},
						},
					},
				},
				{
					Name:                  "s2",
					FetchMaterials:        true,
					CleanWorkingDirectory: false,
					NeverCleanupArtifacts: false,
					Approval: gocd.PipelineApprovalConfig{
						Type:               "success",
						AllowOnlyOnSuccess: false,
						Authorization: gocd.AuthorizationConfig{
							Roles: []string{},
							Users: []string{},
						},
					},
					EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
					Jobs: []gocd.PipelineJobConfig{
						{
							Name:                 "j2",
							RunInstanceCount:     interface{}(nil),
							ElasticProfileID:     "",
							Timeout:              0,
							EnvironmentVariables: []gocd.PipelineEnvironmentVariables{},
							Resources:            []string{},
							Tasks: []gocd.PipelineTaskConfig{
								{
									Type: "fetch",
									Attributes: gocd.TaskAttributeConfig{
										IsSourceAFile:    false,
										ArtifactID:       "docker-image",
										ArtifactOrigin:   "external",
										Command:          "",
										WorkingDirectory: "",
										Pipeline:         "",
										Stage:            "defaultStage",
										Job:              "defaultJob",
										Source:           "",
										Destination:      "",
										Configuration:    []gocd.PluginConfiguration(nil),
										RunIf:            []string{},
										Arguments:        []string(nil),
										PluginConfiguration: struct {
											ID      string "json:\"id,omitempty\" yaml:\"id,omitempty\""
											Version string "json:\"version,omitempty\" yaml:\"version,omitempty\""
										}{
											ID:      "",
											Version: "",
										},
										OnCancel: struct {
											RunIf            []string "json:\"run_if,omitempty\" yaml:\"run_if,omitempty\""
											Command          string   "json:\"command,omitempty\" yaml:\"command,omitempty\""
											Arguments        []string "json:\"arguments,omitempty\" yaml:\"arguments,omitempty\""
											WorkingDirectory string   "json:\"working_directory,omitempty\" yaml:\"working_directory,omitempty\""
										}{
											RunIf:            []string(nil),
											Command:          "",
											Arguments:        []string(nil),
											WorkingDirectory: "",
										},
									},
								},
							},
							Tabs:      []gocd.PipelineTab{},
							Artifacts: []gocd.PipelineArtifact{},
						},
					},
				},
			},
			TrackingTool: gocd.PipelineTracingToolConfig{
				Type: "",
				Attributes: struct {
					URLPattern string "json:\"url_pattern,omitempty\" yaml:\"url_pattern,omitempty\""
					Regex      string "json:\"regex,omitempty\" yaml:\"regex,omitempty\""
				}{URLPattern: "", Regex: ""},
			},
			Timer:         gocd.PipelineTimerConfig{Spec: "", OnlyOnChanges: false},
			CreateOptions: gocd.PipelineCreateOptions{PausePipeline: false, PauseReason: ""},
			Config:        map[string]interface{}(nil), ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		input := gocd.PipelineConfig{
			Name: "new_group",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
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
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
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
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
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
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
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
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
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
			correctPipelineConfigHeader, false, map[string]string{"ETag": "65dbc5f2d5b9c13a2cwxlfkjdlw23654eofixnwe3b3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Group: "new_group",
			Name:  "new_pipeline",
			ETAG:  "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
			CreateOptions: gocd.PipelineCreateOptions{
				PausePipeline: true,
			},
		}

		out, err := client.CreatePipeline(input)
		assert.NoError(t, err)
		assert.Equal(t, "new_group", out.Group)
	})

	t.Run("should error out while creating pipeline configuration in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		_, err := client.CreatePipeline(input)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/pipelines\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating pipeline configuration in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		_, err := client.CreatePipeline(input)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/pipelines\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating pipeline configuration as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pipelineConfigJSON"), http.StatusOK, correctPipelineConfigHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.PipelineConfig{
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		response, err := client.CreatePipeline(input)
		assert.NotNil(t, response)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
	})

	t.Run("should error out while creating pipeline configuration present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		input := gocd.PipelineConfig{
			Name: "new_pipeline",
			ETAG: "65dbc5f2d5b9c13a2ccdlw23654b3b3d8a6155d",
		}

		_, err := client.CreatePipeline(input)
		assert.EqualError(t, err, "call made to create pipeline config 'new_pipeline' errored with: Post "+
			"\"http://localhost:8156/go/api/admin/pipelines\": dial tcp [::1]:8156: connect: connection refused")
	})
}
