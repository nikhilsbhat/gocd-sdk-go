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
	//go:embed internal/fixtures/role_config.json
	roleConfigJSON string
	//go:embed internal/fixtures/roles_config.json
	rolesConfigsJSON string
	//go:embed internal/fixtures/roles_configs_by_type.json
	rolesConfigsByTypeJSON string
)

func Test_client_GetRoles(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to fetch all the roles successfully", func(t *testing.T) {
		server := mockServer([]byte(rolesConfigsJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Role: []gocd.Role{
				{
					Name: "spacetiger",
					Type: "gocd",
					Attributes: gocd.RoleAttribute{
						Users: []string{"alice", "bob", "robin"},
					},
					Policy: []map[string]string{
						{
							"permission": "allow",
							"action":     "view",
							"type":       "environment",
							"resource":   "env_A_*",
						},
					},
				},
				{
					Name: "blackbird",
					Type: "plugin",
					Attributes: gocd.RoleAttribute{
						AuthConfigID: "ldap",
						Properties: []gocd.PluginConfiguration{
							{
								Key:   "UserGroupMembershipAttribute",
								Value: "memberOf",
							},
							{
								Key:   "GroupIdentifiers",
								Value: "ou=admins,ou=groups,ou=system,dc=example,dc=com",
							},
						},
					},
					Policy: []map[string]string{
						{
							"permission": "allow",
							"action":     "view",
							"type":       "environment",
							"resource":   "env_B_*",
						},
					},
				},
			},
		}

		actual, err := client.GetRoles()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(rolesConfigsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRoles()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/roles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(rolesConfigsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRoles()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/roles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("rolesConfigsJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRoles()
		require.EqualError(t, err, "reading response body errored with: invalid character 'r' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRoles()
		require.EqualError(t, err, "call made to get all roles errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/roles\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetRole(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	roleName := "blackbird"

	t.Run("should be able to fetch a specific role successfully", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Role{
			ETAG:         "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Name:         "blackbird",
			Type:         "plugin",
			AuthConfigID: "ldap",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "MemberOf",
					Value: "ou=blackbird,ou=area51,dc=example,dc=com",
				},
			},
			Policy: []map[string]string{
				{
					"permission": "allow",
					"action":     "view",
					"type":       "environment",
					"resource":   "env_A_*",
				},
			},
		}

		actual, err := client.GetRole(roleName)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific role present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Role{}

		actual, err := client.GetRole(roleName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/roles/blackbird\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific role present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Role{}

		actual, err := client.GetRole(roleName)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/roles/blackbird\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific role from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("roleConfigJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.Role{}

		actual, err := client.GetRole(roleName)
		require.EqualError(t, err, "reading response body errored with: invalid character 'r' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific role present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.Role{}

		actual, err := client.GetRole(roleName)
		require.EqualError(t, err, "call made to get role 'blackbird' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/roles/blackbird\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeleteRole(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	roleName := "blackbird"

	t.Run("should be able to delete a role profile successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteRole(roleName)
		require.NoError(t, err)
	})

	t.Run("should error out while deleting role due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteRole(roleName)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/security/roles/blackbird\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting role due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteRole(roleName)
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/security/roles/blackbird\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting role as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteRole(roleName)
		require.EqualError(t, err, "call made to delete role 'blackbird' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/security/roles/blackbird\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_GetRolesByType(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	roleType := "gocd"

	t.Run("should be able to fetch all roles by type successfully", func(t *testing.T) {
		server := mockServer([]byte(rolesConfigsByTypeJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Role: []gocd.Role{
				{
					Name: "spacetiger",
					Type: "gocd",
					Attributes: gocd.RoleAttribute{
						Users: []string{"alice", "bob", "robin"},
					},
					Policy: []map[string]string{
						{
							"permission": "allow",
							"action":     "view",
							"type":       "environment",
							"resource":   "env_A_*",
						},
					},
				},
			},
		}

		actual, err := client.GetRolesByType(roleType)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles by type present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(rolesConfigsJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRolesByType(roleType)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/roles?type=gocd\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles by type present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(rolesConfigsJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRolesByType(roleType)
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/roles?type=gocd\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles by type from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("rolesConfigsJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRolesByType(roleType)
		require.EqualError(t, err, "reading response body errored with: invalid character 'r' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all roles by type present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.RolesConfig{}

		actual, err := client.GetRolesByType(roleType)
		require.EqualError(t, err, "call made to get role by type 'gocd' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/roles?type=gocd\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateRole(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to create role successfully", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{
			Name:         "blackbird",
			Type:         "plugin",
			AuthConfigID: "ldap",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "MemberOf",
					Value: "ou=blackbird,ou=area51,dc=example,dc=com",
				},
			},
			Policy: []map[string]string{
				{
					"permission": "allow",
					"action":     "view",
					"type":       "environment",
					"resource":   "env_A_*",
				},
			},
		}
		expected := role
		expected.ETAG = "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"

		actual, err := client.CreateRole(role)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{}
		expected := role

		actual, err := client.CreateRole(role)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/security/roles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{}
		expected := role

		actual, err := client.CreateRole(role)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/security/roles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("roleConfigJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{}
		expected := role

		actual, err := client.CreateRole(role)
		require.EqualError(t, err, "reading response body errored with: invalid character 'r' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		role := gocd.Role{
			Name: "blackbird",
		}
		expected := gocd.Role{}

		actual, err := client.CreateRole(role)
		require.EqualError(t, err, "call made to create role 'blackbird' errored with: "+
			"Post \"http://localhost:8156/go/api/admin/security/roles\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateRole(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to create role successfully", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			correctArtifactHeader, false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{
			ETAG:         "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
			Name:         "blackbird",
			Type:         "plugin",
			AuthConfigID: "ldap",
			Properties: []gocd.PluginConfiguration{
				{
					Key:   "MemberOf",
					Value: "ou=blackbird,ou=area51,dc=example,dc=com",
				},
			},
			Policy: []map[string]string{
				{
					"permission": "allow",
					"action":     "view",
					"type":       "environment",
					"resource":   "env_A_*",
				},
			},
		}
		expected := role
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateRole(role)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{}
		expected := role

		actual, err := client.UpdateRole(role)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/security/roles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(roleConfigJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{}
		expected := role

		actual, err := client.UpdateRole(role)
		require.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/security/roles\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("roleConfigJSON"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		role := gocd.Role{}
		expected := role

		actual, err := client.UpdateRole(role)
		require.EqualError(t, err, "reading response body errored with: invalid character 'r' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating role in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		role := gocd.Role{
			Name: "blackbird",
		}
		expected := gocd.Role{}

		actual, err := client.UpdateRole(role)
		require.EqualError(t, err, "call made to update role 'blackbird' errored with: "+
			"Put \"http://localhost:8156/go/api/admin/security/roles/blackbird\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
