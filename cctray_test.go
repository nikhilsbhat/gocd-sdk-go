package gocd_test

import (
	_ "embed"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

//go:embed internal/fixtures/cctray.xml
var ccTrayJSON string

func Test_client_GetCCTray(t *testing.T) {
	t.Run("should be able to fetch the cctray successfully", func(t *testing.T) {
		server := mockServer([]byte(ccTrayJSON), http.StatusOK,
			nil, true, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Project{
			{Activity: "Sleeping", LastBuildLabel: "1", LastBuildStatus: "Failure", LastBuildTime: "2023-09-11T03:09:59Z", Name: "action-movies :: build", WebUrl: "http://localhost:8153/go/pipelines/action-movies/1/build/1"},
			{Activity: "Sleeping", LastBuildLabel: "1", LastBuildStatus: "Failure", LastBuildTime: "2023-09-11T03:09:59Z", Name: "action-movies :: build :: build", WebUrl: "http://localhost:8153/go/tab/build/detail/action-movies/1/build/1/build"},
			{Activity: "Sleeping", LastBuildLabel: "1", LastBuildStatus: "Failure", LastBuildTime: "2023-09-11T03:09:59Z", Name: "animation-movies :: build", WebUrl: "http://localhost:8153/go/pipelines/animation-movies/1/build/1"},
		}

		actual, err := client.GetCCTray()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all cctray present in GoCD as GoCD returned non ok status code", func(t *testing.T) {
		server := mockServer([]byte(""), http.StatusBadGateway,
			nil, true, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetCCTray()
		assert.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/cctray.xml\nwith BODY:")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all cctray from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("ccTrayJSON"), http.StatusOK, nil,
			true, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetCCTray()
		assert.EqualError(t, err, "reading response body errored with: EOF")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all cctray present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetCCTray()
		assert.EqualError(t, err, "call made to get cctray errored with: "+
			"Get \"http://localhost:8156/go/cctray.xml\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}
