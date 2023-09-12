package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/materials.json
var materialsJSON string

//go:embed internal/fixtures/materials_usage.json
var materialUsageJSON string

func Test_client_GetMaterials(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionZero}
	t.Run("should be able to fetch all available materials present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Material{
			{
				CanTriggerUpdate: true,
				Config: gocd.MaterialConfig{
					Type:        "plugin",
					Fingerprint: "6bb6a100069baa6968671aa9f6c8ce4f50d9f6b14607e1820bb92b824d26e482",
					Attributes: gocd.Attribute{
						Ref:        "b054d0aa-4704-4369-b54c-7ebdd3394210",
						AutoUpdate: true,
						Origin: map[string]string{
							"type": "config_repo",
							"id":   "sample-repo",
						},
					},
				},
				Messages: []map[string]string{
					{
						"description": "Did not find 'scm' plugin with id 'git-path'. Looks like plugin is missing",
						"level":       "ERROR",
						"message": "Modification check failed for material: WARNING! Plugin missing. " +
							"[url=https://github.com/TWChennai/gocd-git-path-sample.git, path=integration, shallow_clone=true]\nAffected pipelines are both.",
					},
				},
			},
			{
				CanTriggerUpdate: true,
				Config: gocd.MaterialConfig{
					Type:        "git",
					Fingerprint: "ae0cc647060aae14b29e252c9a912a8980a1eef592f6cbfb72a66b5467f93d8e",
					Attributes: gocd.Attribute{
						URL:    "https://github.com/nikhilsbhat/helm-drift.git",
						Branch: "main",
					},
				},
				Messages: []map[string]string{},
			},
		}

		actual, err := client.GetMaterials()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/internal/materials\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/internal/materials\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("materialsJSON"), http.StatusOK, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'm' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all available materials present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		var expected []gocd.Material

		actual, err := client.GetMaterials()
		assert.EqualError(t, err, "call made to get all available materials errored with: "+
			"Get \"http://localhost:8156/go/api/internal/materials\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetMaterialUsage(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionZero}
	t.Run("should be able to fetch usage of a material present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(materialUsageJSON), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []string{
			"pipeline_1",
			"pipeline_2",
			"pipeline_3",
			"pipeline_4",
			"pipeline_5",
			"pipeline_6",
			"pipeline_7",
		}

		actual, err := client.GetMaterialUsage("2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching material usage from GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetMaterialUsage("2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/internal/materials/2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b/usages\nwith "+
			"BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching material usage from GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetMaterialUsage("2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/internal/materials/2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b/usages\nwith "+
			"BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching material usage from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("materialsJSON"), http.StatusOK, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetMaterialUsage("2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b")
		assert.EqualError(t, err, "reading response body errored with: invalid character 'm' looking for beginning of value")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching material usage from GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetMaterialUsage("2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b")
		assert.EqualError(t, err, "call made to get material usage '2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b' errored with: "+
			"Get \"http://localhost:8156/go/api/internal/materials/2faa648612e02becba2b6809fb375f1810acec79498fe50908cebdc5ba0a0a5b/usages\": "+
			"dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

func Test_client_NotifyMaterial(t *testing.T) {
	material := gocd.Material{
		Type:    "git",
		RepoURL: "https://github.com/nikhilsbhat/helm-images",
	}

	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to notify the material present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"message": "The material is now scheduled for an update. Please check relevant pipeline(s) for status."}`), http.StatusAccepted,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.NotifyMaterial(material)
		assert.NoError(t, err)
		assert.Equal(t, "The material is now scheduled for an update. Please check relevant pipeline(s) for status.", actual)
	})

	t.Run("should error out while notifying the material present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.NotifyMaterial(material)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/materials/git/notify\nwith "+
			"BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while notifying the material present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.NotifyMaterial(material)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/materials/git/notify\nwith "+
			"BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while notifying the material present in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"message": "The material is now scheduled for an update.Please check relevant pipeline(s) for status."`),
			http.StatusAccepted, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.NotifyMaterial(material)
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Equal(t, "", actual)
	})

	t.Run("should error out while fetching notifying the material present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.NotifyMaterial(material)
		assert.EqualError(t, err, "call made to notify material 'https://github.com/nikhilsbhat/helm-images' of type git errored with: "+
			"Post \"http://localhost:8156/go/api/admin/materials/git/notify\": "+
			"dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, "", actual)
	})
}

func Test_client_MaterialTriggerUpdate(t *testing.T) {
	materialID := "5fc2198707d4e5b7dfa8cc5c6e398b9ea4bcb17d3aa54f0146ccb361cf03bbd4" //nolint:gosec
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionZero}
	t.Run("should be able to trigger material update present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"message" : "OK"}`), http.StatusCreated,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := map[string]string{"message": "OK"}

		actual, err := client.MaterialTriggerUpdate(materialID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while triggering material update present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.MaterialTriggerUpdate(materialID)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/internal/materials/5fc2198707d4e5b7dfa8cc5c6e398b9ea4bcb17d3aa54f0146ccb361cf03bbd4/trigger_update\nwith "+
			"BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while triggering material update present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(materialsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.MaterialTriggerUpdate(materialID)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/internal/materials/5fc2198707d4e5b7dfa8cc5c6e398b9ea4bcb17d3aa54f0146ccb361cf03bbd4/trigger_update\nwith "+
			"BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while triggering material update present in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"message" : "OK"`),
			http.StatusCreated, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.MaterialTriggerUpdate(materialID)
		assert.EqualError(t, err, "reading response body errored with: unexpected end of JSON input")
		assert.Nil(t, actual)
	})

	t.Run("should error out while triggering material update material present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.MaterialTriggerUpdate(materialID)
		assert.EqualError(t, err, "call made to trigger update '5fc2198707d4e5b7dfa8cc5c6e398b9ea4bcb17d3aa54f0146ccb361cf03bbd4' errored with: "+
			"Post \"http://localhost:8156/go/api/internal/materials/5fc2198707d4e5b7dfa8cc5c6e398b9ea4bcb17d3aa54f0146ccb361cf03bbd4/trigger_update\": "+
			"dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}
