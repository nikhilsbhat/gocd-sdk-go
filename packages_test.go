package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/packages.json
	packagesJSON string
	//go:embed internal/fixtures/package.json
	packageJSON string
)

func Test_client_GetPackages(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to fetch all packages successfully", func(t *testing.T) {
		server := mockServer([]byte(packagesJSON), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.Package{
			{
				CommonConfig: gocd.CommonConfig{
					Name: "package",
					ID:   "f579bb13-bed3-4ad1-a547-ff9d9bcf56d4",
				},
				AutoUpdate: true,
				PackageRepos: gocd.CommonConfig{
					ID:   "273b246e-145d-49d2-a1a4-f0285af9cccc",
					Name: "foo",
				},
				Configuration: []gocd.PluginConfiguration{
					{
						Key:   "PACKAGE_NAME",
						Value: "bar",
					},
				},
			},
			{
				CommonConfig: gocd.CommonConfig{
					Name: "another_package",
					ID:   "c3e1d398-96c0-4edd-a577-9b5c769d449b",
				},
				AutoUpdate: true,
				PackageRepos: gocd.CommonConfig{
					ID:   "273b246e-145d-49d2-a1a4-f0285af9cccc",
					Name: "npm",
				},
				Configuration: []gocd.PluginConfiguration{
					{
						Key:   "PACKAGE_NAME",
						Value: "group",
					},
					{
						Key:   "VERSION_SPEC",
						Value: "1",
					},
					{
						Key:   "ARCHITECTURE",
						Value: "jar",
					},
				},
			},
		}

		actual, err := client.GetPackages()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all packages present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packagesJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Package

		actual, err := client.GetPackages()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/packages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all packages present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packagesJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Package

		actual, err := client.GetPackages()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/packages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all packages from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packagesJSON"), http.StatusOK, correctArtifactHeader,
			false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		var expected []gocd.Package

		actual, err := client.GetPackages()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all packages present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		var expected []gocd.Package

		actual, err := client.GetPackages()
		assert.EqualError(t, err, "call made to get all packages errored with: "+
			"Get \"http://localhost:8156/go/api/admin/packages\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetPackage(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	repositoryID := "dd8926c0-3b4a-4c9e-8012-957b179cec5b"

	t.Run("should be able to fetch a specific package successfully", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Package{
			CommonConfig: gocd.CommonConfig{
				ID:   "package.id",
				Name: "package",
			},
			AutoUpdate: true,
			PackageRepos: gocd.CommonConfig{
				ID:   "273b246e-145d-49d2-a1a4-f0285af9cccc",
				Name: "foo",
			},
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "PACKAGE_NAME",
					Value: "package_name",
				},
			},
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.GetPackage(repositoryID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Package{}

		actual, err := client.GetPackage(repositoryID)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Package{}

		actual, err := client.GetPackage(repositoryID)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Package{}

		actual, err := client.GetPackage(repositoryID)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific package present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.Package{}

		actual, err := client.GetPackage(repositoryID)
		assert.EqualError(t, err, "call made to get package 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreatePackage(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to create a specific package successfully", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{
			CommonConfig: gocd.CommonConfig{
				ID:   "package.id",
				Name: "package",
			},
			AutoUpdate: true,
			PackageRepos: gocd.CommonConfig{
				ID:   "273b246e-145d-49d2-a1a4-f0285af9cccc",
				Name: "foo",
			},
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "PACKAGE_NAME",
					Value: "package_name",
				},
			},
		}

		expected := packageCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.CreatePackage(packageCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{}
		expected := packageCfg

		actual, err := client.CreatePackage(packageCfg)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/packages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{}
		expected := packageCfg

		actual, err := client.CreatePackage(packageCfg)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/packages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{}
		expected := packageCfg

		actual, err := client.CreatePackage(packageCfg)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific package present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		packageCfg := gocd.Package{
			CommonConfig: gocd.CommonConfig{
				ID: "dd8926c0-3b4a-4c9e-8012-957b179cec5b",
			},
		}
		expected := gocd.Package{}

		actual, err := client.CreatePackage(packageCfg)
		assert.EqualError(t, err, "call made to create package 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/packages\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdatePackage(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	repositoryID := "dd8926c0-3b4a-4c9e-8012-957b179cec5b"

	t.Run("should be able to update a specific package successfully", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{
			CommonConfig: gocd.CommonConfig{
				ID:   "package.id",
				Name: "package",
			},
			AutoUpdate: true,
			PackageRepos: gocd.CommonConfig{
				ID:   "273b246e-145d-49d2-a1a4-f0285af9cccc",
				Name: "foo",
			},
			Configuration: []gocd.PluginConfiguration{
				{
					Key:   "PACKAGE_NAME",
					Value: "package_name",
				},
			},
			ETAG: "",
		}

		expected := packageCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdatePackage(packageCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{}
		expected := packageCfg

		actual, err := client.UpdatePackage(packageCfg)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/packages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(packageJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{}
		expected := packageCfg

		actual, err := client.UpdatePackage(packageCfg)
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/packages\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("packageJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		packageCfg := gocd.Package{}
		expected := packageCfg

		actual, err := client.UpdatePackage(packageCfg)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'p' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating a specific package present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		packageCfg := gocd.Package{
			CommonConfig: gocd.CommonConfig{
				ID: repositoryID,
			},
		}
		expected := gocd.Package{}

		actual, err := client.UpdatePackage(packageCfg)
		assert.EqualError(t, err, "call made to update package 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeletePackage(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}
	repositoryID := "dd8926c0-3b4a-4c9e-8012-957b179cec5b"

	t.Run("should be able to delete an appropriate package successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePackage(repositoryID)
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting a package due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePackage(repositoryID)
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting a package due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeletePackage(repositoryID)
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting a package as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeletePackage(repositoryID)
		assert.EqualError(t, err, "call made to delete package 'dd8926c0-3b4a-4c9e-8012-957b179cec5b' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/packages/dd8926c0-3b4a-4c9e-8012-957b179cec5b\": dial tcp [::1]:8156: connect: connection refused")
	})
}
