package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/scheduled_jobs.xml
var scheduledJobJSON string

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

func Test_client_RunJobs(t *testing.T) {
	correctJobsHeader := map[string]string{"Accept": gocd.HeaderVersionThree, gocd.HeaderConfirm: "true"}

	t.Run("should error out while running selected jobs of a pipeline from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.RunJobs(gocd.Stage{})
		assert.EqualError(t, err, "call made to run selected jobs errored with: "+
			"Post \"http://localhost:8156/go/api/stages/run-selected-jobs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while running selected jobs as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("scheduledJobJSON"), http.StatusBadGateway, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.RunJobs(gocd.Stage{})
		assert.EqualError(t, err, "got 502 from GoCD while making POST call for "+server.URL+
			"/api/stages/run-selected-jobs\nwith BODY:scheduledJobJSON")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while running selected jobs as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"`), http.StatusAccepted, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.RunJobs(gocd.Stage{})
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, "", actual)
	})

	t.Run("should be able to run jobs from a pipeline present in GoCD", func(t *testing.T) {
		runJobResponse := `{"message": "Request to rerun jobs accepted"}`
		server := mockServer([]byte(runJobResponse), http.StatusAccepted, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		stage := gocd.Stage{
			Pipeline:         "pipeline1",
			PipelineInstance: "2",
			Name:             "myStage1",
			StageCounter:     "3",
			Jobs:             []string{"myJob1", "myJob2"},
		}
		expected := "Request to rerun jobs accepted"

		actual, err := client.RunJobs(stage)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_RunFailedJobs(t *testing.T) {
	correctJobsHeader := map[string]string{"Accept": gocd.HeaderVersionThree, gocd.HeaderConfirm: "true"}

	t.Run("should error out while running failed jobs from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.RunFailedJobs(gocd.Stage{})
		assert.EqualError(t, err, "call made to run failed jobs errored with: "+
			"Post \"http://localhost:8156/go/api/stages/run-failed-jobs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while running failed jobs as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("scheduledJobJSON"), http.StatusBadGateway, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.RunFailedJobs(gocd.Stage{})
		assert.EqualError(t, err, "got 502 from GoCD while making POST call for "+server.URL+
			"/api/stages/run-failed-jobs\nwith BODY:scheduledJobJSON")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while running failed jobs as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"`), http.StatusAccepted, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.RunFailedJobs(gocd.Stage{})
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, "", actual)
	})

	t.Run("should be able to run failed jobs present in GoCD", func(t *testing.T) {
		runFailedJobResponse := `{"message": "Request to rerun jobs accepted"}`
		server := mockServer([]byte(runFailedJobResponse), http.StatusAccepted, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		stage := gocd.Stage{
			Pipeline:         "pipeline1",
			PipelineInstance: "2",
			Name:             "mystage",
			StageCounter:     "3",
			Jobs:             nil,
		}

		expected := "Request to rerun jobs accepted"
		actual, err := client.RunFailedJobs(stage)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
