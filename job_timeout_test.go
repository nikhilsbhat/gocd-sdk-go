package gocd_test

import (
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_client_GetDefaultJobTimeout(t *testing.T) {
	t.Run("should be able fetch the default job timeout successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"default_job_timeout" : "0"}`), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetDefaultJobTimeout()
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"default_job_timeout": "0"}, actual)
	})

	t.Run("should error out while fetching default job timeout due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(`{"default_job_timeout" : "0"}`), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetDefaultJobTimeout()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config/server/default_job_timeout\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("should error out while fetching default job timeout due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(`{"default_job_timeout" : "0"}`), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetDefaultJobTimeout()
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/admin/config/server/default_job_timeout\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("should error out while fetching default job timeout as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`default_job_timeout" : "0"}`), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetDefaultJobTimeout()
		require.EqualError(t, err, "reading response body errored with: invalid character 'd' looking for beginning of value")
		assert.Equal(t, map[string]string(nil), actual)
	})

	t.Run("should error out while fetching default job timeout as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetDefaultJobTimeout()
		require.EqualError(t, err, "call made to get default job timeout errored with: "+
			"Get \"http://localhost:8156/go/api/admin/config/server/default_job_timeout\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, map[string]string(nil), actual)
	})
}

func Test_client_UpdateDefaultJobTimeout(t *testing.T) {
	t.Run("should be able update the default job timeout successfully", func(t *testing.T) {
		server := mockServer([]byte(`{"default_job_timeout" : "10"}`), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionOne}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.UpdateDefaultJobTimeout(10)
		require.NoError(t, err)
	})

	t.Run("should error out while updating default job timeout due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(`{"default_job_timeout" : "0"}`), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.UpdateDefaultJobTimeout(10)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/config/server/default_job_timeout\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while updating default job timeout due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(`{"default_job_timeout" : "0"}`), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.UpdateDefaultJobTimeout(10)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/admin/config/server/default_job_timeout\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while updating default job timeout as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.UpdateDefaultJobTimeout(10)
		require.EqualError(t, err, "call made to update default job timeout errored with: "+
			"Post \"http://localhost:8156/go/api/admin/config/server/default_job_timeout\": dial tcp [::1]:8156: connect: connection refused")
	})
}
