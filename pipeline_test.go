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
	//go:embed internal/fixtures/pipeline_schedule.json
	pipelineSchedule string
	//go:embed internal/fixtures/pipeline_instance.json
	pipelineInstance string
)

var pipelineMap = map[string]interface{}{
	"can_run":               true,
	"comment":               interface{}(nil),
	"counter":               float64(1),
	"label":                 "1",
	"name":                  "PipelineName",
	"natural_order":         float64(1),
	"preparing_to_schedule": false,
	"scheduled_date":        1.436519914578e+12,
	"build_cause": map[string]interface{}{
		"approver": "",
		"material_revisions": []interface{}{map[string]interface{}{
			"changed": true,
			"material": map[string]interface{}{
				"description": "URL: https://github.com/gocd/gocd, Branch: master",
				"fingerprint": "de08b34d116a1c0cf57cd76683bf21",
				"name":        "https://github.com/gocd/gocd",
				"type":        "Git",
			},
			"modifications": []interface{}{map[string]interface{}{
				"comment":       "some commit message.",
				"email_address": interface{}(nil),
				"modified_time": 1.436519914378e+12,
				"revision":      "40f0a7ef224a0a2fba438b158483b",
				"user_name":     "user <user@users.noreply.github.com>",
			}},
		}},
		"trigger_forced":  false,
		"trigger_message": "modified by user <user@users.noreply.github.com>",
	},
	"stages": []interface{}{
		map[string]interface{}{
			"approval_type": "success",
			"approved_by":   "changes",
			"can_run":       true,
			"counter":       "1",
			"jobs": []interface{}{map[string]interface{}{
				"name":           "job",
				"result":         "Passed",
				"scheduled_date": 1.436782534378e+12,
				"state":          "Completed",
			}},
			"name":               "stage",
			"operate_permission": true,
			"rerun_of_counter":   interface{}(nil),
			"result":             "Passed",
			"scheduled":          true,
			"status":             "Completed",
		},
	},
}

func Test_client_GetPipelines(t *testing.T) {
	t.Run("should error out while fetching pipelines from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, "call made to get pipelines errored with "+
			"Get \"http://localhost:8156/go/api/feed/pipelines.xml\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, "reading response body errored with: EOF")
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should be able to fetch the pipelines present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelines), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelines()
		assert.NoError(t, err)
		assert.Equal(t, 3, len(actual.Pipeline))
	})
}

func Test_client_getPipelineName(t *testing.T) {
	t.Run("should be able to fetch pipeline name from the href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8156/go/api/feed/pipelines/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "animation-and-action-movies", name)
	})

	t.Run("should fetch malformed pipeline name as malformed/invalid (prefix) href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8156/go/api/feed/pipelinesss/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "/go/api/feed/pipelinesss/animation-and-action-movies", name)
	})

	t.Run("should fetch malformed pipeline name as malformed/invalid (suffix) href passed", func(t *testing.T) {
		pipelineLink := "http://localhost:8156/go/api/feed/pipelines/animation-and-action-movies/test/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.NoError(t, err)
		assert.Equal(t, "animation-and-action-movies/test", name)
	})

	t.Run("should fail while fetching pipeline name from malformed/invalid href", func(t *testing.T) {
		pipelineLink := "://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml"
		name, err := gocd.GetPipelineName(pipelineLink)
		assert.EqualError(t, err, "parsing URL errored with parse "+
			"\"://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml\": missing protocol scheme")
		assert.Equal(t, "", name)
	})
}

func Test_client_GetPipelineStatus(t *testing.T) {
	pipeline := "action-movies-manual"
	correctPipelineHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should error out while fetching pipeline statuses information from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, "call made to get pipeline state errored with Get "+
			"\"http://localhost:8156/go/api/pipelines/action-movies-manual/status\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should error out while fetching pipeline statuses information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should error out while fetching pipeline statuses information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should be able to get pipeline statuses of all pipeline present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelineState), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

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

func Test_client_PipelinePause(t *testing.T) {
	correctPipelinePauseHeader := map[string]string{
		"Accept":       gocd.HeaderVersionOne,
		"Content-Type": gocd.ContentJSON,
	}
	t.Run("Should be able to pause pipeline successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctPipelinePauseHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelinePause("first_pipeline", "pausing the pipeline")
		assert.NoError(t, err)
	})

	t.Run("Should error out while pausing pipeline due to wrong header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelinePause("first_pipeline", "pausing the pipeline")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("Should error out while pausing pipeline due missing header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelinePause("first_pipeline", "pausing the pipeline")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("Should error out while pausing pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.PipelinePause("first_pipeline", "pausing the pipeline")
		assert.EqualError(t, err, "call made to pause pipeline errored with "+
			"Post \"http://localhost:8156/go/api/pipelines/first_pipeline/pause\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_PipelineUnPause(t *testing.T) {
	correctPipelinePauseHeader := map[string]string{
		"Accept":           gocd.HeaderVersionOne,
		gocd.HeaderConfirm: "true",
	}
	t.Run("Should be able to unpause pipeline successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctPipelinePauseHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnPause("first_pipeline")
		assert.NoError(t, err)
	})

	t.Run("Should error out while un pausing pipeline due to wrong header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnPause("first_pipeline")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("Should error out while un pausing pipeline due missing header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnPause("first_pipeline")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("Should error out while un pausing pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.PipelineUnPause("first_pipeline")
		assert.EqualError(t, err, "call made to unpause pipeline errored with "+
			"Post \"http://localhost:8156/go/api/pipelines/first_pipeline/unpause\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_PipelineUnlock(t *testing.T) {
	correctPipelinePauseHeader := map[string]string{"Accept": gocd.HeaderVersionOne, gocd.HeaderConfirm: "true"}
	t.Run("Should be able to unlock pipeline successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctPipelinePauseHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnlock("first_pipeline")
		assert.NoError(t, err)
	})

	t.Run("Should error out while unlocking pipeline due to wrong header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnlock("first_pipeline")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("Should error out while unlocking pipeline due missing header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnlock("first_pipeline")
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("Should error out while unlocking pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.PipelineUnlock("first_pipeline")
		assert.EqualError(t, err, "call made to unlock pipeline errored with "+
			"Post \"http://localhost:8156/go/api/pipelines/first_pipeline/unlock\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_SchedulePipeline(t *testing.T) {
	correctPipelinePauseHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": "application/json"}
	t.Run("should be able to schedule pipeline successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusAccepted, correctPipelinePauseHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		schedule := gocd.Schedule{
			EnvVars: []map[string]interface{}{
				{
					"TEST_ENV13": "value_env6",
				},
			},
		}

		err := client.SchedulePipeline("first_pipeline", schedule)
		assert.NoError(t, err)
	})

	t.Run("should error out while scheduling pipeline due to wrong header", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": "application/json"}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		schedule := gocd.Schedule{}

		err := client.SchedulePipeline("first_pipeline", schedule)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while scheduling pipeline due to missing header", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		schedule := gocd.Schedule{}

		err := client.SchedulePipeline("first_pipeline", schedule)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while scheduling pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "debug", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)
		schedule := gocd.Schedule{
			EnvVars: []map[string]interface{}{
				{
					"name":   "TEST_ENV13",
					"value":  "value_env6",
					"secure": false,
				},
			},
			UpdateMaterial: true,
		}

		err := client.SchedulePipeline("first_pipeline", schedule)
		assert.EqualError(t, err, "call made to schedule pipeline 'first_pipeline' errored with "+
			"Post \"http://localhost:8156/go/api/pipelines/first_pipeline/schedule\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_CommentOnPipeline(t *testing.T) {
	correctPipelineHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}
	comment := gocd.PipelineObject{
		Name:    "pipeline1",
		Counter: 1,
		Message: "this is test comment",
	}

	t.Run("should be able to comment on selected pipeline successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "debug", nil)

		err := client.CommentOnPipeline(comment)
		assert.NoError(t, err)
	})

	t.Run("should error out while commenting on pipeline due to wrong header", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.CommentOnPipeline(comment)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while commenting on pipeline due to missing header", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.CommentOnPipeline(comment)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while commenting on pipeline due to missing fields in Comment object", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		comments := gocd.PipelineObject{
			Name:    "pipeline1",
			Counter: 1,
		}

		err := client.CommentOnPipeline(comments)
		assert.EqualError(t, err, "comment message cannot be empty")
	})

	t.Run("should error out while commenting on pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.CommentOnPipeline(comment)
		assert.EqualError(t, err, "call made to comment on pipeline 'pipeline1' errored with "+
			"Post \"http://localhost:8156/go/api/pipelines/pipeline1/1/comment\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_GetPipelineInstance(t *testing.T) {
	correctPipelineHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	pipelineObj := gocd.PipelineObject{
		Name:    "pipeline1",
		Counter: 1,
	}
	t.Run("should be able to fetch the pipeline instance from GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineInstance), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "debug", nil)

		expected := pipelineMap

		actual, err := client.GetPipelineInstance(pipelineObj)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline instance due to wrong header", func(t *testing.T) {
		server := mockServer([]byte(pipelineInstance), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineInstance(pipelineObj)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching pipeline instance due to missing header", func(t *testing.T) {
		server := mockServer([]byte(pipelineInstance), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineInstance(pipelineObj)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching pipeline as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("{pipelineInstance}"), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineInstance(pipelineObj)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of object key string")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineInstance(pipelineObj)
		assert.EqualError(t, err, "call made to fetch pipeline instance 'pipeline1' errored with "+
			"Get \"http://localhost:8156/go/api/pipelines/pipeline1/1\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

// func Test_client_GetPipelineHistory2(t *testing.T) {
//	t.Run("should be able to fetch the pipeline history successfully", func(t *testing.T) {
//		client := gocd.NewClient(
//			"http://localhost:8153/go",
//			"admin",
//			"admin",
//			"info",
//			nil,
//		)
//
//		actual, err := client.GetPipelineHistory("gocd-prometheus-exporter", 10, 0)
//		assert.NoError(t, err)
//
//		for _, pipeline := range actual {
//			log.Println(pipeline["name"], pipeline["counter"])
//		}
//		assert.Equal(t, "", actual)
//	})
// }

// func Test_client_GetPipelineHistory(t *testing.T) {
//	correctPipelineHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
//	server1 := mockServer([]byte(pipelineHistory), http.StatusOK, correctPipelineHeader, false, nil)
//	server2 := mockServer([]byte(pipelineHistory), http.StatusOK, correctPipelineHeader, false, nil)
//
//	type errorTestCases struct {
//		description   string
//		mockServer    *httptest.Server
//		expectedError bool
//		errorString   string
//		expected      []map[string]interface{}
//	}
//
//	expectOne := []map[string]interface{}{
//		{"name": "pipeline1", "counter": 1},
//		{"name": "pipeline1", "counter": 2},
//		{"name": "pipeline1", "counter": 3},
//		{"name": "pipeline1", "counter": 4},
//		{"name": "pipeline1", "counter": 5},
//		{"name": "pipeline1", "counter": 6},
//		{"name": "pipeline1", "counter": 7},
//		{"name": "pipeline1", "counter": 8},
//		{"name": "pipeline1", "counter": 9},
//		{"name": "pipeline1", "counter": 10},
//	}
//	expectTwo := []map[string]interface{}{
//		{"name": "pipeline1", "counter": 11},
//		{"name": "pipeline1", "counter": 12},
//		{"name": "pipeline1", "counter": 13},
//		{"name": "pipeline1", "counter": 14},
//		{"name": "pipeline1", "counter": 15},
//		{"name": "pipeline1", "counter": 16},
//		{"name": "pipeline1", "counter": 17},
//		{"name": "pipeline1", "counter": 18},
//		{"name": "pipeline1", "counter": 19},
//		{"name": "pipeline1", "counter": 20},
//	}
//
//	tests := []errorTestCases{
//		{
//			description: "should be able to paginate once successfully",
//			mockServer:  server1,
//			expected:    expectOne,
//		},
//		{
//			description: "should be able to paginate once successfully",
//			mockServer:  server2,
//			expected:    expectTwo,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.description, func(t *testing.T) {
//			client := gocd.NewClient(
//				tt.mockServer.URL,
//				"admin",
//				"admin",
//				"info",
//				nil,
//			)
//			got, err := client.GetPipelineHistory("pipeline1", 10, 0)
//			assert.NoError(t, err)
//			assert.Equal(t, tt.expected, got)
//		})
//	}
// }
