package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/authorization.json
	authConfigGetJSON string
	//go:embed internal/fixtures/authorizations.json
	authConfigsGetJSON string
)

func Test_client_GetAuthConfigs(t *testing.T) {
	t.Run("should be able to fetch all auth configs present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(authConfigsGetJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.CommonConfig{{
			ID:                  "ldap",
			PluginID:            "cd.go.authentication.ldap",
			AllowOnlyKnownUsers: false,
			Properties: []gocd.PluginConfiguration{
				{Key: "Url", Value: "ldap://ldap.server.url:389"},
			},
		}}

		actual, err := client.GetAuthConfigs()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all auth configs present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.CommonConfig(nil)

		actual, err := client.GetAuthConfigs()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/auth_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all auth configs present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.CommonConfig(nil)

		actual, err := client.GetAuthConfigs()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/auth_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a auth configs from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("authConfigsGetJSON"), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo},
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.CommonConfig(nil)

		actual, err := client.GetAuthConfigs()
		assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a auth configs present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := []gocd.CommonConfig(nil)

		actual, err := client.GetAuthConfigs()
		assert.EqualError(t, err, "call made to get auth configs errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/auth_configs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetAuthConfig(t *testing.T) {
	t.Run("should be able to fetch a selected auth config present in GoCD successfully", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{
			ID:                  "ldap",
			PluginID:            "cd.go.authentication.ldap",
			AllowOnlyKnownUsers: false,
			Properties: []gocd.PluginConfiguration{
				{Key: "Url", Value: "ldap://ldap.server.url:389"},
				{Key: "SearchBase", Value: "ou=users,ou=system"},
				{Key: "ManagerDN", Value: "uid=admin,ou=system"},
				{Key: "SearchFilter", Value: "uid"},
				{Key: "Password", EncryptedValue: "GLzARJ+Qr+M="},
				{Key: "DisplayNameAttribute", Value: "displayName"},
				{Key: "EmailAttribute", Value: "mail"},
			},
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}

		actual, err := client.GetAuthConfig("ldap")
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a auth config present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetAuthConfig("ldap")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/auth_configs/ldap\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a auth config present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK, nil,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetAuthConfig("ldap")
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/auth_configs/ldap\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a auth config present in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("authConfigGetJSON"), http.StatusOK, map[string]string{"Accept": gocd.HeaderVersionTwo},
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.GetAuthConfig("ldap")
		assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a auth config present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{}

		actual, err := client.GetAuthConfig("ldap")
		assert.EqualError(t, err, "call made to get auth config 'ldap' errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/auth_configs/ldap\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateAuthConfig(t *testing.T) {
	correctAuthHeader := map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to create auth config successfully", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK, correctAuthHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		authCfg := gocd.CommonConfig{
			ID:                  "ldap",
			PluginID:            "cd.go.authentication.ldap",
			AllowOnlyKnownUsers: false,
			Properties: []gocd.PluginConfiguration{
				{Key: "Url", Value: "ldap://ldap.server.url:389"},
				{Key: "SearchBase", Value: "ou=users,ou=system"},
				{Key: "ManagerDN", Value: "uid=admin,ou=system"},
				{Key: "SearchFilter", Value: "uid"},
				{Key: "Password", EncryptedValue: "GLzARJ+Qr+M="},
				{Key: "DisplayNameAttribute", Value: "displayName"},
				{Key: "EmailAttribute", Value: "mail"},
			},
		}
		expected := authCfg
		expected.ETAG = "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"

		actual, err := client.CreateAuthConfig(authCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating auth config due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.CreateAuthConfig(gocd.CommonConfig{})
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/security/auth_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating auth config due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			nil, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.CreateAuthConfig(gocd.CommonConfig{})
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/security/auth_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating auth config as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("authConfigGetJSON"), http.StatusOK, correctAuthHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.CreateAuthConfig(gocd.CommonConfig{})
		assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating auth config as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{}

		actual, err := client.CreateAuthConfig(gocd.CommonConfig{ID: "ldap"})
		assert.EqualError(t, err, "call made to create auth config 'ldap' errored with:"+
			" Post \"http://localhost:8156/go/api/admin/security/auth_configs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateAuthConfig(t *testing.T) {
	correctAuthHeader := map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON, "If-Match": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"}

	t.Run("should be able to update auth config successfully", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK, correctAuthHeader,
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		authCfg := gocd.CommonConfig{
			ID:                  "ldap",
			PluginID:            "cd.go.authentication.ldap",
			AllowOnlyKnownUsers: false,
			Properties: []gocd.PluginConfiguration{
				{Key: "Url", Value: "ldap://ldap.server.url:389"},
				{Key: "SearchBase", Value: "ou=users,ou=system"},
				{Key: "ManagerDN", Value: "uid=admin,ou=system"},
				{Key: "SearchFilter", Value: "uid"},
				{Key: "Password", EncryptedValue: "GLzARJ+Qr+M="},
				{Key: "DisplayNameAttribute", Value: "displayName"},
				{Key: "EmailAttribute", Value: "mail"},
			},
			ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d",
		}
		expected := authCfg
		expected.ETAG = "61406622382e51c2079c11dcbdb978fb"

		actual, err := client.UpdateAuthConfig(authCfg)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating auth config due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.UpdateAuthConfig(gocd.CommonConfig{})
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/security/auth_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating auth config due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(authConfigGetJSON), http.StatusOK,
			nil, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.UpdateAuthConfig(gocd.CommonConfig{})
		assert.EqualError(t, err, "got 404 from GoCD while making PUT call for "+server.URL+
			"/api/admin/security/auth_configs\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating auth config as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("authConfigGetJSON"), http.StatusOK, correctAuthHeader,
			false, map[string]string{"ETag": "61406622382e51c2079c11dcbdb978fb"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.CommonConfig{}

		actual, err := client.UpdateAuthConfig(gocd.CommonConfig{ETAG: "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		assert.EqualError(t, err, "invalid character 'a' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating auth config as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.CommonConfig{}

		actual, err := client.UpdateAuthConfig(gocd.CommonConfig{ID: "ldap"})
		assert.EqualError(t, err, "call made to update auth config 'ldap' errored with:"+
			" Put \"http://localhost:8156/go/api/admin/security/auth_configs\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeleteAuthConfig(t *testing.T) {
	t.Run("should be able to delete auth config successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteAuthConfig("ldap")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting auth config due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteAuthConfig("ldap")
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/security/auth_configs/ldap\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting auth config due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteAuthConfig("ldap")
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/admin/security/auth_configs/ldap\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting auth config as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteAuthConfig("ldap")
		assert.EqualError(t, err, "call made to delete auth config 'ldap' errored with: "+
			"Delete \"http://localhost:8156/go/api/admin/security/auth_configs/ldap\": dial tcp [::1]:8156: connect: connection refused")
	})
}
