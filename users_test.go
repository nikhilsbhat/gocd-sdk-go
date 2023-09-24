package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

var (
	//go:embed internal/fixtures/user.json
	userJSON string
	//go:embed internal/fixtures/user_current.json
	currentUserJSON string
	//go:embed internal/fixtures/user_current_update.json
	currentUserUpdateJSON string
	//go:embed internal/fixtures/users.json
	usersJSON string
)

func Test_client_DeleteUser(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to delete an appropriate user successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctUserHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteUser("sample")
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting a user due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteUser("sample-user")
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/users/sample-user\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting a user due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteUser("sample-user")
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/users/sample-user\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting a user as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteUser("sample-user")
		assert.EqualError(t, err, "call made to delete user 'sample-user' errored with: "+
			"Delete \"http://localhost:8156/go/api/users/sample-user\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_BulkDeleteUsers(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to bulk delete users successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctUserHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		usersToDelete := map[string]interface{}{
			"users": []string{"jez", "tez"},
		}

		err := client.BulkDeleteUsers(usersToDelete)
		assert.NoError(t, err)
	})

	t.Run("should error out while bulk deleting users due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		usersToDelete := map[string]interface{}{
			"users": []string{"jez", "tez"},
		}

		err := client.BulkDeleteUsers(usersToDelete)
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/users\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while bulk deleting users due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		usersToDelete := map[string]interface{}{
			"users": []string{"jez", "tez"},
		}

		err := client.BulkDeleteUsers(usersToDelete)
		assert.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/users\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while bulk deleting users as GoCD server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		usersToDelete := map[string]interface{}{
			"users": []string{"jez", "tez"},
		}

		err := client.BulkDeleteUsers(usersToDelete)
		assert.EqualError(t, err, "call made to bulk delete users errored with: "+
			"Delete \"http://localhost:8156/go/api/users\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_BulkEnableDisableUsers(t *testing.T) {
	correctBulkUpdateHeader := map[string]string{"Accept": gocd.HeaderVersionThree, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to bulk enable/disable the users successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			correctBulkUpdateHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		users := map[string]interface{}{
			"users": []string{"jez", "tez"},
			"operations": map[string]interface{}{
				"enable": true,
			},
		}

		err := client.BulkEnableDisableUsers(users)
		assert.NoError(t, err)
	})

	t.Run("should error out while bulk enabling/disabling users present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.BulkEnableDisableUsers(nil)
		assert.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/admin/operations/state\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while bulk enabling/disabling users present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.BulkEnableDisableUsers(nil)
		assert.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/admin/operations/state\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while bulk enabling/disabling users present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.BulkEnableDisableUsers(nil)
		assert.EqualError(t, err, "call made to bulk enable/disable users errored with: Patch "+
			"\"http://localhost:8156/go/api/admin/operations/state\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_GetUser(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionThree}
	userName := "jdoe"

	t.Run("should be able to fetch user details successfully", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK,
			correctUserHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{
			Name:         "jdoe",
			LoginName:    "jdoe",
			Enabled:      true,
			Admin:        true,
			CheckInAlias: []string{},
			Roles: []gocd.UserRole{
				{
					Name: "role name",
					Type: "gocd",
				},
			},
		}

		actual, err := client.GetUser(userName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific user present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionSeven}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}

		actual, err := client.GetUser(userName)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/users/jdoe\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific user present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}

		actual, err := client.GetUser(userName)
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/users/jdoe\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific user from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("userJSON"), http.StatusOK, correctUserHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}

		actual, err := client.GetUser(userName)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'u' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching a specific user present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.User{}

		actual, err := client.GetUser(userName)
		assert.EqualError(t, err, "call made to get user information errored with: "+
			"Get \"http://localhost:8156/go/api/users/jdoe\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetCurrentUser(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should be able to fetch current user details successfully", func(t *testing.T) {
		server := mockServer([]byte(currentUserJSON), http.StatusOK,
			correctUserHeader, false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{
			Name:         "John Doe",
			LoginName:    "jdoe",
			Enabled:      true,
			CheckInAlias: []string{},
		}

		actual, err := client.GetCurrentUser()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching current user present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(currentUserJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionSeven}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}

		actual, err := client.GetCurrentUser()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/current_user\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching current user present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(currentUserJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}

		actual, err := client.GetCurrentUser()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/current_user\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching current user from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("currentUserJSON"), http.StatusOK, correctUserHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}

		actual, err := client.GetCurrentUser()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching current user present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.User{}

		actual, err := client.GetCurrentUser()
		assert.EqualError(t, err, "call made to get current user information errored with: "+
			"Get \"http://localhost:8156/go/api/current_user\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_GetUsers(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to fetch users details successfully", func(t *testing.T) {
		server := mockServer([]byte(usersJSON), http.StatusOK, correctUserHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := []gocd.User{{
			Name:         "John Doe",
			LoginName:    "jdoe",
			Enabled:      true,
			Admin:        false,
			CheckInAlias: []string{"jdoe", "johndoe"},
			Roles: []gocd.UserRole{
				{
					Name: "role name",
					Type: "gocd",
				},
			},
		}}

		actual, err := client.GetUsers()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching all users present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(usersJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionSeven}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetUsers()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/users\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all users present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(usersJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetUsers()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/users\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all users from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("usersJSON"), http.StatusOK, correctUserHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetUsers()
		assert.EqualError(t, err, "reading response body errored with: invalid character 'u' looking for beginning of value")
		assert.Nil(t, actual)
	})

	t.Run("should error out while fetching all users present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetUsers()
		assert.EqualError(t, err, "call made to get all users errored with: "+
			"Get \"http://localhost:8156/go/api/users\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
	})
}

func Test_client_CreateUser(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to create a specific user successfully", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK, correctUserHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{
			Name:         "jdoe",
			LoginName:    "jdoe",
			Enabled:      true,
			Admin:        true,
			CheckInAlias: []string{},
			Roles: []gocd.UserRole{
				{
					Name: "role name",
					Type: "gocd",
				},
			},
		}
		expected := user

		actual, err := client.CreateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.CreateUser(user)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/users\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.CreateUser(user)
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/users\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("clusterProfileJSON"), http.StatusOK, correctUserHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.CreateUser(user)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.CreateUser(user)
		assert.EqualError(t, err, "call made to create user 'jdoe' errored with: Post "+
			"\"http://localhost:8156/go/api/users\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateUser(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionThree}

	t.Run("should be able to create a specific user successfully", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK, correctUserHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{
			Name:         "jdoe",
			LoginName:    "jdoe",
			Enabled:      true,
			Admin:        true,
			CheckInAlias: []string{},
			Roles: []gocd.UserRole{
				{
					Name: "role name",
					Type: "gocd",
				},
			},
		}
		expected := user

		actual, err := client.UpdateUser(user)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(userJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.UpdateUser(user)
		assert.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/users/jdoe\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(clusterProfileJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.UpdateUser(user)
		assert.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/users/jdoe\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("clusterProfileJSON"), http.StatusOK, correctUserHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.UpdateUser(user)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating a specific user in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		user := gocd.User{Name: "jdoe"}
		expected := gocd.User{}

		actual, err := client.UpdateUser(user)
		assert.EqualError(t, err, "call made to update user 'jdoe' errored with: Patch "+
			"\"http://localhost:8156/go/api/users/jdoe\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_UpdateCurrentUser(t *testing.T) {
	correctUserHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to update current user details successfully", func(t *testing.T) {
		server := mockServer([]byte(currentUserUpdateJSON), http.StatusOK,
			correctUserHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{
			Name:      "jdoe",
			LoginName: "jdoe",
			Enabled:   false,
			EmailID:   "jdoe@example.com",
			EmailMe:   true,
			CheckInAlias: []string{
				"jdoe",
				"johndoe",
			},
		}

		user := gocd.User{
			EmailID: "jdoe@example.com",
			EmailMe: true,
			CheckInAlias: []string{
				"jdoe",
				"johndoe",
			},
		}

		actual, err := client.UpdateCurrentUser(user)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating current user present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(currentUserJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionSeven}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{}
		expected := gocd.User{}

		actual, err := client.UpdateCurrentUser(user)
		assert.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/current_user\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating current user present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(currentUserUpdateJSON), http.StatusOK, nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		user := gocd.User{}
		expected := gocd.User{}

		actual, err := client.UpdateCurrentUser(user)
		assert.EqualError(t, err, "got 404 from GoCD while making PATCH call for "+server.URL+
			"/api/current_user\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating current user from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("currentUserJSON"), http.StatusOK, correctUserHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.User{}
		user := gocd.User{}

		actual, err := client.UpdateCurrentUser(user)
		assert.EqualError(t, err, "reading response body errored with: invalid character 'c' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while updating current user present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.User{}
		user := gocd.User{}

		actual, err := client.UpdateCurrentUser(user)
		assert.EqualError(t, err, "call made to update current user information errored with: "+
			"Patch \"http://localhost:8156/go/api/current_user\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
