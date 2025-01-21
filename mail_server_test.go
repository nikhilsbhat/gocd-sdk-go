package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed internal/fixtures/mail_server_config.json
var mailServerJSON string

func Test_client_GetMailServerConfig(t *testing.T) {
	t.Run("should be able to fetch mail server configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.MailServerConfig{
			Hostname:          "smtp.example.com",
			Port:              465,
			UserName:          "user@example.com",
			EncryptedPassword: "AES:lzcCuNSe4vUx+CsWgN11Uw==:Q2OlnqIf9S++yMPqSCx8qg==",
			TLS:               true,
			SenderEmail:       "no-reply@example.com",
			AdminEmail:        "gocd-admins@example.com",
		}

		actual, err := client.GetMailServerConfig()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching mail server configuration due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.MailServerConfig{}

		actual, err := client.GetMailServerConfig()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/config/mailserver\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching mail server configuration due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.MailServerConfig{}

		actual, err := client.GetMailServerConfig()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/config/mailserver\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching mail server configuration as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("mailServerJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.MailServerConfig{}

		actual, err := client.GetMailServerConfig()
		require.EqualError(t, err, "reading response body errored with: invalid character 'm' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching mail server configuration as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.MailServerConfig{}

		actual, err := client.GetMailServerConfig()
		require.EqualError(t, err, "call made to get mail server config errored with: "+
			"Get \"http://localhost:8156/go/api/config/mailserver\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_DeleteMailServerConfig(t *testing.T) {
	t.Run("should be able to delete mail server configuration successfully", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteMailServerConfig()
		require.NoError(t, err)
	})

	t.Run("should error out while deleting mail server configuration due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteMailServerConfig()
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/config/mailserver\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting mail server configuration due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteMailServerConfig()
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/config/mailserver\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting mail server configuration as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteMailServerConfig()
		require.EqualError(t, err, "call made to delete mail server config errored with: "+
			"Delete \"http://localhost:8156/go/api/config/mailserver\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_CreateOrUpdateMailServerConfig(t *testing.T) {
	t.Run("should be create/update mail server configuration successfully", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.MailServerConfig{
			Hostname:          "smtp.example.com",
			Port:              465,
			UserName:          "user@example.com",
			EncryptedPassword: "AES:lzcCuNSe4vUx+CsWgN11Uw==:Q2OlnqIf9S++yMPqSCx8qg==",
			TLS:               true,
			SenderEmail:       "no-reply@example.com",
			AdminEmail:        "gocd-admins@example.com",
		}

		expected := input

		actual, err := client.CreateOrUpdateMailServerConfig(input)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating/updating mail server configuration due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.MailServerConfig{}
		expected := input

		actual, err := client.CreateOrUpdateMailServerConfig(input)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/config/mailserver\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating/updating mail server configuration due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(mailServerJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.MailServerConfig{}
		expected := input

		actual, err := client.CreateOrUpdateMailServerConfig(input)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/config/mailserver\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating/updating mail server configuration as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("mailServerJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		input := gocd.MailServerConfig{}
		expected := input

		actual, err := client.CreateOrUpdateMailServerConfig(input)
		require.EqualError(t, err, "reading response body errored with: invalid character 'm' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating/updating mail server configuration as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		input := gocd.MailServerConfig{}
		expected := input

		actual, err := client.CreateOrUpdateMailServerConfig(input)
		require.EqualError(t, err, "call made to create or update mail server config errored with: "+
			"Post \"http://localhost:8156/go/api/config/mailserver\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
