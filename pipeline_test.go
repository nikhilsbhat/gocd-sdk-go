package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"
	"time"

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
	//go:embed internal/fixtures/scheduled_jobs.xml
	scheduledJobJSON string
	//go:embed internal/fixtures/pipeline_schedules.json
	pipelineSchedulesJSON string
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
		assert.EqualError(t, err, "call made to get pipelines errored with: "+
			"Get \"http://localhost:8156/go/api/feed/pipelines.xml\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.PipelinesInfo{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelines()
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/feed/pipelines.xml\nwith BODY:backupJSON")
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

func Test_client_GetPipelineSchedules(t *testing.T) {
	t.Run("should error out while fetching pipeline schedules from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineSchedules("helm-images", "0", "2")
		assert.EqualError(t, err, "call made to get pipeline schedules helm-images errored with: Get "+
			"\"http://localhost:8156/go/pipelineHistory.json?pipelineName=helm-images&perPage=2&start=0\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.PipelineSchedules{}, actual)
	})

	t.Run("should error out while fetching pipeline schedules as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("pipelineSchedulesJSON"), http.StatusBadGateway, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineSchedules("helm-images", "0", "2")
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/pipelineHistory.json?pipelineName=helm-images&perPage=2&start=0\nwith BODY:pipelineSchedulesJSON")
		assert.Equal(t, gocd.PipelineSchedules{}, actual)
	})

	t.Run("should error out while fetching pipeline schedules as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"pipelineSchedulesJSON"}`), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineSchedules("helm-images", "0", "2")
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.PipelineSchedules{}, actual)
	})

	t.Run("should be able to fetch the pipeline schedules present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedulesJSON), http.StatusOK, map[string]string{
			"Accept":       gocd.HeaderVersionZero,
			"Content-Type": gocd.ContentJSON,
		}, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineSchedules{
			Name:  "helm-images",
			Count: 8, Groups: []gocd.PipelineSchedulesGroups{
				{
					History: []gocd.PipelineSchedulesHistory{
						{
							Label:              "8",
							ScheduledDate:      "25 Jun, 2023 at 12:48:41 [+0530]",
							ScheduledTimestamp: 1.687677521085e+12,
							ModificationDate:   "about 1 month ago",
							BuildCause:         "Triggered by admin",
						},
						{
							Label:              "7",
							ScheduledDate:      "18 Jun, 2023 at 18:39:04 [+0530]",
							ScheduledTimestamp: 1.68709374451e+12,
							ModificationDate:   "about 1 month ago",
							BuildCause:         "Triggered by changes",
						},
					},
				},
			},
		}

		actual, err := client.GetPipelineSchedules("helm-images", "0", "2")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetPipelineHistory(t *testing.T) {
	t.Run("should error out while fetching pipeline run history from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineRunHistory("helm-images", "0", time.Duration(2)*time.Second)
		assert.EqualError(t, err, "call made to get pipeline helm-images errored with: "+
			"Get \"http://localhost:8156/go/api/pipelines/helm-images/history?after=0&page_size=0\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching pipeline run history as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("pipelineRunHistoryJSON"), http.StatusBadGateway, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineRunHistory("helm-images", "0", time.Duration(2)*time.Second)
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/pipelines/helm-images/history?after=0&page_size=0\nwith BODY:pipelineRunHistoryJSON")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching pipeline run history as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"pipelineRunHistoryJSON"}`), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineRunHistory("helm-images", "0", time.Duration(2)*time.Second)
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	// t.Run("should be able to fetch the pipeline run history present in GoCD", func(t *testing.T) {
	//	server := mockServer([]byte(pipelineRunHistoryJSON), http.StatusOK, map[string]string{
	//		"Accept":       gocd.HeaderVersionOne,
	//		"Content-Type": gocd.ContentJSON,
	//	}, true, nil)
	//	client := gocd.NewClient(server.URL, auth, "info", nil)
	//
	//	expected := []gocd.PipelineRunHistory{
	//		{
	//			Name:          "helm-images",
	//			Counter:       3,
	//			ScheduledDate: 1678470766332,
	//			BuildCause:    gocd.PipelineBuildCause{Message: "Forced by admin", Approver: "admin", TriggerForced: true},
	//		},
	//		{
	//			Name:          "helm-images",
	//			Counter:       2,
	//			ScheduledDate: 1677128882155,
	//			BuildCause:    gocd.PipelineBuildCause{Message: "modified by nikhilsbhat <nikhilsbhat93@gmail.com>", Approver: "changes", TriggerForced: false},
	//		},
	//		{
	//			Name:          "helm-images",
	//			Counter:       1,
	//			ScheduledDate: 1672544013154,
	//			BuildCause:    gocd.PipelineBuildCause{Message: "Forced by admin", Approver: "admin", TriggerForced: true},
	//		},
	//	}
	//
	//	actual, err := client.GetPipelineRunHistory("helm-images", "0")
	//	assert.NoError(t, err)
	//	assert.Equal(t, expected, actual)
	// })
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
		assert.EqualError(t, err, "parsing URL errored with:"+
			" parse \"://localhost:8153/go/api/feed/pipelines/animation-and-action-movies/stages.xml\": missing protocol scheme")
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
		assert.EqualError(t, err, "call made to get pipeline state errored with:"+
			" Get \"http://localhost:8156/go/api/pipelines/action-movies-manual/status\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.PipelineState{}, actual)
	})

	t.Run("should error out while fetching pipeline statuses information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineState(pipeline)
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/pipelines/action-movies-manual/status\nwith BODY:backupJSON")
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
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/pause\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("Should error out while pausing pipeline due missing header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelinePause("first_pipeline", "pausing the pipeline")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/pause\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("Should error out while pausing pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.PipelinePause("first_pipeline", "pausing the pipeline")
		assert.EqualError(t, err, "call made to pause pipeline errored with: "+
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
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/unpause\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("Should error out while un pausing pipeline due missing header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnPause("first_pipeline")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/unpause\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("Should error out while un pausing pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.PipelineUnPause("first_pipeline")
		assert.EqualError(t, err, "call made to unpause pipeline errored with: "+
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
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/unlock\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("Should error out while unlocking pipeline due missing header", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.PipelineUnlock("first_pipeline")
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/unlock\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("Should error out while unlocking pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.PipelineUnlock("first_pipeline")
		assert.EqualError(t, err, "call made to unlock pipeline errored with: "+
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
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/schedule\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while scheduling pipeline due to missing header", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		schedule := gocd.Schedule{}

		err := client.SchedulePipeline("first_pipeline", schedule)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/first_pipeline/schedule\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
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
		assert.EqualError(t, err, "call made to schedule pipeline 'first_pipeline' errored with: "+
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
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/pipeline1/1/comment\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while commenting on pipeline due to missing header", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.CommentOnPipeline(comment)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/pipelines/pipeline1/1/comment\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while commenting on pipeline due to missing fields in Comment object", func(t *testing.T) {
		server := mockServer([]byte(pipelineSchedule), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		comments := gocd.PipelineObject{
			Name:    "pipeline1",
			Counter: 1,
		}

		err := client.CommentOnPipeline(comments)
		assert.EqualError(t, err, "comment message cannot be empty <nil>")
	})

	t.Run("should error out while commenting on pipeline as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.CommentOnPipeline(comment)
		assert.EqualError(t, err, "call made to comment on pipeline 'pipeline1' errored with: "+
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
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/pipelines/pipeline1/1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching pipeline instance due to missing header", func(t *testing.T) {
		server := mockServer([]byte(pipelineInstance), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineInstance(pipelineObj)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/pipelines/pipeline1/1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
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
		assert.EqualError(t, err, "call made to fetch pipeline instance 'pipeline1' errored with: "+
			"Get \"http://localhost:8156/go/api/pipelines/pipeline1/1\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

func Test_client_ScheduledJobs(t *testing.T) {
	t.Run("should error out while fetching scheduled jobs from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetScheduledJobs()
		assert.EqualError(t, err, "call made to get scheduled jobs errored with: "+
			"Get \"http://localhost:8156/go/api/feed/jobs/scheduled.xml\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.ScheduledJobs{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("scheduledJobJSON"), http.StatusBadGateway, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetScheduledJobs()
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/feed/jobs/scheduled.xml\nwith BODY:scheduledJobJSON")
		assert.Equal(t, gocd.ScheduledJobs{}, actual)
	})

	t.Run("should error out while fetching pipelines as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"`), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetScheduledJobs()
		assert.EqualError(t, err, "reading response body errored with: EOF")
		assert.Equal(t, gocd.ScheduledJobs{}, actual)
	})

	t.Run("should be able to fetch the pipelines present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(scheduledJobJSON), http.StatusOK, nil, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.ScheduledJobs{
			Job: []gocd.Job{
				{
					Name:         "job1",
					ID:           "6",
					BuildLocator: "mypipeline/5/defaultStage/1/job1",
					Environment:  "sample_environment",
				},
			},
		}

		actual, err := client.GetScheduledJobs()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
