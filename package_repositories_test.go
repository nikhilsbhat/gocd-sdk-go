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
	//go:embed internal/fixtures/package_repositories.json
	packageRepositoriesJSON string
	//go:embed internal/fixtures/package_repository.json
	packageRepositoryJSON string
)

func Test_client_GetPackageRepositories(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should be able to fetch the package repositories successfully", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoriesJSON), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.PackageRepository{
			{
				ID:   "dd8926c0-3b4a-4c9e-8012-957b179cec5b",
				Name: "repository",
				PluginMetaData: map[string]string{
					"id":      "deb",
					"version": "1",
				},
				Configuration: []gocd.PluginConfiguration{
					{
						Key:   "REPO_URL",
						Value: "http://sample-repo",
					},
				},
				Packages: struct {
					Packages []gocd.CommonConfig `json:"packages,omitempty" yaml:"packages,omitempty"`
				}(struct{ Packages []gocd.CommonConfig }{
					Packages: []gocd.CommonConfig{
						{
							Name: "package",
							ID:   "6bba891e-2675-49af-b16d-200bd6c6801e",
						},
					},
				}),
			},
		}

		actual, err := client.GetPackageRepositories()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all package repositories present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoriesJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.PackageRepository

		actual, err := client.GetPackageRepositories()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/repositories\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all package repositories present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoriesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.PackageRepository

		actual, err := client.GetPackageRepositories()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/repositories\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all package repositories from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageRepositoriesJSON"), http.StatusOK, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.PackageRepository

		actual, err := client.GetPackageRepositories()
		require.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all package repositories present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		var expected []gocd.PackageRepository

		actual, err := client.GetPackageRepositories()
		require.EqualError(t, err, "call made to get package repositories errored with: "+
			"Get \"http://localhost:8156/go/api/admin/repositories\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetPackageRepository(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	repositoryID := "dd8926c0-3b4a-4c9e-8012-957b179cec5b"

	t.Run("should be able to fetch a specific package repository successfully", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PackageRepository{
			ID:   "dd8926c0-3b4a-4c9e-8012-957b179cec5b",
			Name: "repository",
			PluginMetaData: map[string]string{
				"id":      "deb",
				"version": "1",
			},
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "REPO_URL",
					Value: "http://sample-repo",
				},
			},
			Packages: struct {
				Packages []gocd.CommonConfig `json:"packages,omitempty" yaml:"packages,omitempty"`
			}{
				Packages: []gocd.CommonConfig{
					{
						Name: "package",
						ID:   "6bba891e-2675-49af-b16d-200bd6c6801e",
					},
				},
			},
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.GetPackageRepository(repositoryID)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package repository present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PackageRepository{}

		actual, err := client.GetPackageRepository(repositoryID)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package repository present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PackageRepository{}

		actual, err := client.GetPackageRepository(repositoryID)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package repository from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageRepositoryJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.PackageRepository{}

		actual, err := client.GetPackageRepository(repositoryID)
		require.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package repository present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.PackageRepository{}

		actual, err := client.GetPackageRepository(repositoryID)
		require.EqualError(t, err, "call made to get package repository 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreatePackageRepository(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should be able to create a specific package repository successfully", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{
			ID:   "dd8926c0-3b4a-4c9e-8012-957b179cec5b",
			Name: "repository",
			PluginMetaData: map[string]string{
				"id":      "deb",
				"version": "1",
			},
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "REPO_URL",
					Value: "http://sample-repo",
				},
			},
			Packages: struct {
				Packages []gocd.CommonConfig `json:"packages,omitempty" yaml:"packages,omitempty"`
			}{
				Packages: []gocd.CommonConfig{
					{
						Name: "package",
						ID:   "6bba891e-2675-49af-b16d-200bd6c6801e",
					},
				},
			},
		}

		expected := repositoryCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.CreatePackageRepository(repositoryCfg)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package repository present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{}
		expected := repositoryCfg

		actual, err := client.UpdatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/repositories\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package repository present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{}
		expected := repositoryCfg

		actual, err := client.CreatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/repositories\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package repository from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageRepositoryJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{}
		expected := repositoryCfg

		actual, err := client.CreatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package repository present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		repositoryCfg := gocd.PackageRepository{ID: "dd8926c0-3b4a-4c9e-8012-957b179cec5b"}
		expected := gocd.PackageRepository{}

		actual, err := client.CreatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "call made to create package repository 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/repositories\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdatePackageRepository(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	repositoryID := "dd8926c0-3b4a-4c9e-8012-957b179cec5b"

	t.Run("should be able to update a specific package repository successfully", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			ID:   "dd8926c0-3b4a-4c9e-8012-957b179cec5b",
			Name: "repository",
			PluginMetaData: map[string]string{
				"id":      "deb",
				"version": "1",
			},
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "REPO_URL",
					Value: "http://sample-repo",
				},
			},
			Packages: struct {
				Packages []gocd.CommonConfig `json:"packages,omitempty" yaml:"packages,omitempty"`
			}{
				Packages: []gocd.CommonConfig{
					{
						Name: "package",
						ID:   "6bba891e-2675-49af-b16d-200bd6c6801e",
					},
				},
			},
		}

		expected := repositoryCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdatePackageRepository(repositoryCfg)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package repository present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{}
		expected := repositoryCfg

		actual, err := client.UpdatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/repositories\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package repository present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageRepositoryJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{}
		expected := repositoryCfg

		actual, err := client.UpdatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/repositories\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package repository from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageRepositoryJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		repositoryCfg := gocd.PackageRepository{}
		expected := repositoryCfg

		actual, err := client.UpdatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package repository present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		repositoryCfg := gocd.PackageRepository{ID: repositoryID}
		expected := gocd.PackageRepository{}

		actual, err := client.UpdatePackageRepository(repositoryCfg)
		require.EqualError(t, err, "call made to update package repository 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeletePackageRepository(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	repositoryID := "dd8926c0-3b4a-4c9e-8012-957b179cec5b"

	t.Run("should be able to delete an appropriate package repository successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePackageRepository(repositoryID)
		require.NoError(t, err)
	})

	t.Run("should error out while deleting a package repository due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePackageRepository(repositoryID)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting a package repository due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePackageRepository(repositoryID)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting a package repository as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeletePackageRepository(repositoryID)
		require.EqualError(t, err, "call made to delete package repository 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/repositories/dd8926c0-3b4a-4c9e-8012-957b179cec5b\": dial tcp [::1]:8156: connect: connection refused")
	})
}
