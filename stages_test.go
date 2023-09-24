package gocd_test

import (
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

func Test_client_RunStage(t *testing.T) {
	correctJobsHeader := map[string]string{"Accept": gocd.HeaderVersionTwo, gocd.HeaderConfirm: "true"}

	t.Run("should error out while running stage from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.RunStage(gocd.Stage{})
		assert.EqualError(t, err, "call made to run stage errored with: "+
			"Post \"http://localhost:8156/go/api/stages/run\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while running stage from a pipeline as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("scheduledJobJSON"), http.StatusBadGateway, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.RunStage(gocd.Stage{})
		assert.EqualError(t, err, "got 502 from GoCD while making POST call for "+server.URL+
			"/api/stages/run\nwith BODY:scheduledJobJSON")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while running stage from a pipeline as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"`), http.StatusAccepted, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.RunStage(gocd.Stage{})
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, "", actual)
	})

	t.Run("should be able to run stage from a pipeline present in GoCD", func(t *testing.T) {
		runFailedJobResponse := `{"message": "Request to schedule stage pipeline1/2/myStage accepted"}`
		server := mockServer([]byte(runFailedJobResponse), http.StatusAccepted, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		stage := gocd.Stage{
			Pipeline:         "pipeline1",
			PipelineInstance: "2",
			Name:             "myStage",
			Jobs:             nil,
		}

		expected := "Request to schedule stage pipeline1/2/myStage accepted"
		actual, err := client.RunStage(stage)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CancelStage(t *testing.T) {
	correctJobsHeader := map[string]string{"Accept": gocd.HeaderVersionThree, gocd.HeaderConfirm: "true"}

	t.Run("should error out while cancelling stage from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.CancelStage(gocd.Stage{})
		assert.EqualError(t, err, "call made to cancel stage errored with: "+
			"Post \"http://localhost:8156/go/api/stages/cancel\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while cancelling stage from a pipeline as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("scheduledJobJSON"), http.StatusBadGateway, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.CancelStage(gocd.Stage{})
		assert.EqualError(t, err, "got 502 from GoCD while making POST call for "+server.URL+
			"/api/stages/cancel\nwith BODY:scheduledJobJSON")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while cancelling stage from a pipeline as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"`), http.StatusOK, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.CancelStage(gocd.Stage{})
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, "", actual)
	})

	t.Run("should be able to cancel stage from a pipeline present in GoCD", func(t *testing.T) {
		runFailedJobResponse := `{"message": "Request to schedule stage pipeline1/2/myStage accepted"}`
		server := mockServer([]byte(runFailedJobResponse), http.StatusOK, correctJobsHeader, true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		stage := gocd.Stage{
			Pipeline:         "pipeline1",
			PipelineInstance: "2",
			Name:             "myStage",
			StageCounter:     "2",
			Jobs:             nil,
		}

		expected := "Request to schedule stage pipeline1/2/myStage accepted"
		actual, err := client.CancelStage(stage)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
