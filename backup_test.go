package gocd_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/backup.json
var backupJSON string

func TestConfig_GetBackupInfo(t *testing.T) {
	correctBackupHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should error out while fetching latest backup configuration information from server", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"admin",
			"admin",
			"info",
			nil,
		)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetBackupConfig()
		assert.EqualError(t, err, "call made to get backup information errored with "+
			"Get \"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should error out while fetching latest backup configuration information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctBackupHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetBackupConfig()
		assert.EqualError(t, err, "body: backupJSON httpcode: 502")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should error out while fetching latest backup configuration information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctBackupHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetBackupConfig()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should be able to fetch the latest backup configuration information available in GoCD", func(t *testing.T) {
		server := mockServer([]byte(backupJSON), http.StatusOK, correctBackupHeader, false, nil)

		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		expected := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		actual, err := client.GetBackupConfig()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateOrUpdateBackup(t *testing.T) {
	correctBackupHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}
	t.Run("should be able to create/update the backup configurations successfully", func(t *testing.T) {
		var backupInfo gocd.BackupConfig
		err := json.Unmarshal([]byte(backupJSON), &backupInfo)
		assert.NoError(t, err)

		server := backupMockServer(backupInfo, http.MethodPost, correctBackupHeader)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err = client.CreateOrUpdateBackupConfig(backupObj)
		assert.NoError(t, err)
	})

	t.Run("should error while creating or updating backup configuration due to wrong headers", func(t *testing.T) {
		server := backupMockServer([]byte("backupJSON"), http.MethodPost, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON})
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error while creating or updating backup configuration due to missing headers", func(t *testing.T) {
		server := backupMockServer([]byte("backupJSON"), http.MethodPost, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error while creating or updating backup configuration as malformed data sent to server", func(t *testing.T) {
		server := backupMockServer([]byte("backupJSON"), http.MethodPost, correctBackupHeader)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		assert.EqualError(t, err, "body: json: cannot unmarshal string into Go value of type gocd.BackupConfig httpcode: 500")
	})

	t.Run("should error while creating or updating backup configuration as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		assert.EqualError(t, err, "call made to create/udpate backup configuration errored with Post "+
			"\"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_DeleteBackupConfig(t *testing.T) {
	correctBackupHeader := map[string]string{"Accept": gocd.HeaderVersionOne}
	t.Run("should be able to delete the backup configurations successfully ", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctBackupHeader, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.DeleteBackupConfig()
		assert.NoError(t, err)
	})

	t.Run("should error out while deleting the backup configurations due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(
			server.URL,
			"",
			"",
			"info",
			nil,
		)

		err := client.DeleteBackupConfig()
		assert.EqualError(t, err, "body: <html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html> httpcode: 404")
	})

	t.Run("should error out while deleting the backup configurations as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8156/go",
			"",
			"",
			"info",
			nil,
		)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteBackupConfig()
		assert.EqualError(t, err, "call made to get backup information errored with: "+
			"Delete \"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func backupMockServer(request interface{}, method string, header map[string]string) *httptest.Server { //nolint:unparam
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if header == nil {
			writer.WriteHeader(http.StatusNotFound)
			if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
				log.Fatalln(err)
			}

			return
		}

		for key, value := range header {
			val := req.Header.Get(key)
			_ = val
			if req.Header.Get(key) != value {
				writer.WriteHeader(http.StatusNotFound)
				if _, err := writer.Write([]byte(`<html>
<body>
	<h2>404 Not found</h2>
</body>

</html>`)); err != nil {
					log.Fatalln(err)
				}

				return
			}
		}

		if method == http.MethodDelete {
			writer.WriteHeader(http.StatusOK)
			if _, err := writer.Write([]byte(`{"message": "Backup config was deleted successfully!"}`)); err != nil {
				log.Fatalln(err)
			}
		}

		requestByte, err := json.Marshal(request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err = writer.Write([]byte(fmt.Sprintf("%s %s", string(requestByte), err.Error()))); err != nil {
				log.Fatalln(err)
			}

			return
		}

		var backupCfg gocd.BackupConfig
		if err = json.Unmarshal(requestByte, &backupCfg); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			if _, err = writer.Write([]byte(err.Error())); err != nil {
				log.Fatalln(err)
			}

			return
		}

		writer.WriteHeader(http.StatusOK)
		if _, err = writer.Write([]byte("")); err != nil {
			log.Fatalln(err)
		}
	}))
}
