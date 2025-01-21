package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/fixtures/pipeline_groups.json
	pipelineGroups string
	//go:embed internal/fixtures/pipeline_group.json
	pipelineGroup string
	//go:embed internal/fixtures/pipeline_group_update.json
	pipelineGroupUpdate string
)

func Test_client_GetPipelineGroupInfo(t *testing.T) {
	correctPipelineHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should error out while fetching all pipeline groups information from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetPipelineGroups()
		require.EqualError(t, err, "call made to get pipeline groups information errored with: "+
			"Get \"http://localhost:8156/go/api/admin/pipeline_groups\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all pipeline groups information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineGroups()
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+
			"/api/admin/pipeline_groups\nwith BODY:backupJSON")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all pipeline groups information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetPipelineGroups()
		require.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Nil(t, actual)
	})

	t.Run("should be able to get information of all pipeline groups present in GoCD", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroups), http.StatusOK, correctPipelineHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.PipelineGroup{
			{
				Authorization: gocd.PipelineGroupAuthorizationConfig{
					View: gocd.AuthorizationConfig{
						Users: []string{"operate"},
						Roles: []string{},
					},
					Admins: gocd.AuthorizationConfig{
						Users: []string{"operate"},
						Roles: []string{},
					},
				},
				Name:          "action-movies",
				PipelineCount: 2,
				Pipelines:     []gocd.Pipeline{{Name: "action-movies-auto"}, {Name: "action-movies-manual"}},
			},
			{
				Authorization: gocd.PipelineGroupAuthorizationConfig{
					View: gocd.AuthorizationConfig{
						Users: []string{"operate"},
						Roles: []string{},
					},
					Admins: gocd.AuthorizationConfig{
						Users: []string{"operate"},
						Roles: []string{},
					},
				},
				Name:          "infrastructure",
				PipelineCount: 2,
				Pipelines:     []gocd.Pipeline{{Name: "gocd-prometheus-exporter"}, {Name: "helm-images"}},
			},
		}

		actual, err := client.GetPipelineGroups()
		require.NoError(t, err)
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
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePipelineGroup("pipeline_group_1")
		require.NoError(t, err)
	})

	t.Run("should error out while deleting pipeline group due to wrong headers passed", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePipelineGroup("pipeline_group_1")
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/pipeline_groups/pipeline_group_1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting pipeline group due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePipelineGroup("pipeline_group_1")
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/pipeline_groups/pipeline_group_1\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting pipeline group as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeletePipelineGroup("pipeline_group_1")
		require.EqualError(t, err, "call made to delete pipeline group errored with:"+
			" Delete \"http://localhost:8156/go/api/admin/pipeline_groups/pipeline_group_1\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_GetPipelineGroup(t *testing.T) {
	t.Run("should be able to fetch pipeline group successfully from GoCD server", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroup), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false,
			map[string]string{"ETag": "17f5a9edf150884e5fc4315b4a7814cd"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PipelineGroup{
			Name:          "first",
			PipelineCount: 0,
			Pipelines:     []gocd.Pipeline{{Name: "up42"}},
			Authorization: gocd.PipelineGroupAuthorizationConfig{
				View: gocd.AuthorizationConfig{
					Users: []string{"operate"},
					Roles: []string{},
				},
				Admins: gocd.AuthorizationConfig{
					Users: []string{"operate"},
					Roles: []string{},
				},
			},
			ETAG: "17f5a9edf150884e5fc4315b4a7814cd",
		}

		actual, err := client.GetPipelineGroup("first")
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline group from GoCD server due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroup), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		expected := gocd.PipelineGroup{}

		actual, err := client.GetPipelineGroup("first")
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/pipeline_groups/first\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline group from GoCD server due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroup), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		expected := gocd.PipelineGroup{}

		actual, err := client.GetPipelineGroup("first")
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/pipeline_groups/first\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline group as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pipelineGroup"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false,
			map[string]string{"ETag": "17f5a9edf150884e5fc4315b4a7814cd"})
		client := gocd.NewClient(server.URL, auth, "info", nil)
		expected := gocd.PipelineGroup{}

		actual, err := client.GetPipelineGroup("first")
		require.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching pipeline group as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.PipelineGroup{}

		actual, err := client.GetPipelineGroup("first")
		require.EqualError(t, err, "call made to fetch pipeline group errored with: "+
			"Get \"http://localhost:8156/go/api/admin/pipeline_groups/first\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreatePipelineGroup(t *testing.T) {
	t.Run("should be able to create pipeline group with specified configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroup), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		group := gocd.PipelineGroup{
			Name: "first",
			Pipelines: []gocd.Pipeline{{
				Name: "name",
			}},
			Authorization: gocd.PipelineGroupAuthorizationConfig{
				View: gocd.AuthorizationConfig{
					Users: []string{"operate"},
				},
				Admins: gocd.AuthorizationConfig{
					Users: []string{"operate"},
				},
			},
		}

		err := client.CreatePipelineGroup(group)
		require.NoError(t, err)
	})

	t.Run("should error out while creating pipeline group from GoCD server due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroup), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		group := gocd.PipelineGroup{}

		err := client.CreatePipelineGroup(group)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/pipeline_groups\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating pipeline group from GoCD server due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroup), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)
		group := gocd.PipelineGroup{}

		err := client.CreatePipelineGroup(group)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/pipeline_groups\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while creating pipeline group as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		group := gocd.PipelineGroup{Name: "first"}

		err := client.CreatePipelineGroup(group)
		require.EqualError(t, err, "call made to create pipeline group 'first' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/pipeline_groups\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_UpdatePipelineGroup(t *testing.T) {
	etag := "17f5a9edf150884e5fc4315b4a7814cd"
	correctPipelineGroupHeader := map[string]string{
		"Accept":       gocd.HeaderVersionOne,
		"Content-Type": gocd.ContentJSON,
		"If-Match":     etag,
	}

	t.Run("should be able to update the pipeline group successfully", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroupUpdate), http.StatusOK, correctPipelineGroupHeader,
			false, map[string]string{"ETag": "28f5a8edf130994e6fc4315b4a7814cd"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		group := gocd.PipelineGroup{
			Name: "first",
			Authorization: gocd.PipelineGroupAuthorizationConfig{
				Operate: gocd.AuthorizationConfig{
					Users: []string{"alice"},
					Roles: []string{},
				},
			},
			ETAG: etag,
		}

		expected := gocd.PipelineGroup{
			Name:          "first",
			PipelineCount: 0,
			Pipelines: []gocd.Pipeline{
				{Name: "up42"},
			},
			Authorization: gocd.PipelineGroupAuthorizationConfig{
				Admins: gocd.AuthorizationConfig{
					Roles: []string{},
					Users: []string{"operate"},
				},
				Operate: gocd.AuthorizationConfig{
					Roles: []string(nil),
					Users: []string{"alice"},
				},
				View: gocd.AuthorizationConfig{
					Roles: []string{},
					Users: []string{"operate"},
				},
			},
			ETAG: "28f5a8edf130994e6fc4315b4a7814cd",
		}

		actual, err := client.UpdatePipelineGroup(group)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline group due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroupUpdate), http.StatusOK, map[string]string{
			"Accept":       gocd.HeaderVersionThree,
			"Content-Type": gocd.ContentJSON,
			"If-Match":     etag,
		},
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		group := gocd.PipelineGroup{}
		expected := gocd.PipelineGroup{}

		actual, err := client.UpdatePipelineGroup(group)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/pipeline_groups\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline group due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(pipelineGroupUpdate), http.StatusOK, map[string]string{
			"Accept":       gocd.HeaderVersionTwo,
			"Content-Type": gocd.ContentJSON,
		},
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		group := gocd.PipelineGroup{}
		expected := gocd.PipelineGroup{}

		actual, err := client.UpdatePipelineGroup(group)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/pipeline_groups\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline group as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("pipelineGroupUpdate"), http.StatusOK, correctPipelineGroupHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		group := gocd.PipelineGroup{
			Name: "first",
			Authorization: gocd.PipelineGroupAuthorizationConfig{
				Operate: gocd.AuthorizationConfig{
					Users: []string{"alice"},
					Roles: []string{},
				},
			},
			ETAG: etag,
		}
		expected := gocd.PipelineGroup{}

		actual, err := client.UpdatePipelineGroup(group)
		require.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating pipeline group due to missing headers", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		group := gocd.PipelineGroup{Name: "first"}
		expected := gocd.PipelineGroup{}

		actual, err := client.UpdatePipelineGroup(group)
		require.EqualError(t, err, "call made to update pipeline group 'first' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/pipeline_groups/first\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
