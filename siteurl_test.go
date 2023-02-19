package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/site_url.json
var siteURLJSON string

func Test_client_GetSiteURL(t *testing.T) {
	t.Run("should be able to get the site url successfully", func(t *testing.T) {
		server := mockServer([]byte(siteURLJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{
			SiteURL:       "http://foo.com",
			SecureSiteURL: "https://foo.com",
		}

		actual, err := client.GetSiteURL()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while getting the site url due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(siteURLJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{}

		actual, err := client.GetSiteURL()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/site_urls\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while getting the site url due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(siteURLJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{}

		actual, err := client.GetSiteURL()
		assert.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/security/site_urls\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while getting site url as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("siteURLJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{}

		actual, err := client.GetSiteURL()
		assert.EqualError(t, err, "invalid character 's' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while getting site url as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.SiteURLConfig{}

		actual, err := client.GetSiteURL()
		assert.EqualError(t, err, "call made to get site url errored with: "+
			"Get \"http://localhost:8156/go/api/admin/security/site_urls\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateOrUpdateSiteURL(t *testing.T) {
	t.Run("should be able to create/update the site url successfully", func(t *testing.T) {
		server := mockServer([]byte(siteURLJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		site := gocd.SiteURLConfig{
			SiteURL:       "http://foo.com",
			SecureSiteURL: "https://foo.com",
		}

		expected := gocd.SiteURLConfig{
			SiteURL:       "http://foo.com",
			SecureSiteURL: "https://foo.com",
		}

		actual, err := client.CreateOrUpdateSiteURL(site)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating/updating the site url due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(siteURLJSON), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{}

		actual, err := client.CreateOrUpdateSiteURL(gocd.SiteURLConfig{})
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/security/site_urls\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while creating/updating the site url due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(siteURLJSON), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{}

		actual, err := client.CreateOrUpdateSiteURL(gocd.SiteURLConfig{})
		assert.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/security/site_urls\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while getting site url as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("siteURLJSON"), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.SiteURLConfig{}

		actual, err := client.CreateOrUpdateSiteURL(gocd.SiteURLConfig{})
		assert.EqualError(t, err, "invalid character 's' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while getting site url as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.SiteURLConfig{}

		actual, err := client.CreateOrUpdateSiteURL(gocd.SiteURLConfig{})
		assert.EqualError(t, err, "call made to create/update site url errored with: "+
			"Post \"http://localhost:8156/go/api/admin/security/site_urls\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}
