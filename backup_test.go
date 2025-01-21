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
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/fixtures/backup.json
	backupJSON string
	//go:embed internal/fixtures/backup_stats.json
	backupStats string
)

func TestConfig_GetBackupInfo(t *testing.T) {
	correctBackupHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should error out while fetching latest backup configuration information from server", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)
		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.GetBackupConfig()
		require.EqualError(t, err, "call made to get backup information errored with: "+
			"Get \"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should error out while fetching latest backup configuration information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, correctBackupHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetBackupConfig()
		require.EqualError(t, err, "got 502 from GoCD while making GET call for "+server.URL+"/api/config/backup\nwith BODY:backupJSON")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should error out while fetching latest backup configuration information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, correctBackupHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.GetBackupConfig()
		require.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should be able to fetch the latest backup configuration information available in GoCD", func(t *testing.T) {
		server := mockServer([]byte(backupJSON), http.StatusOK, correctBackupHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		actual, err := client.GetBackupConfig()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_client_CreateOrUpdateBackup(t *testing.T) {
	correctBackupHeader := map[string]string{"Accept": gocd.HeaderVersionOne, "Content-Type": gocd.ContentJSON}

	t.Run("should be able to create/update the backup configurations successfully", func(t *testing.T) {
		var backupInfo gocd.BackupConfig
		err := json.Unmarshal([]byte(backupJSON), &backupInfo)
		require.NoError(t, err)

		server := backupMockServer(backupInfo, http.MethodPost, correctBackupHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err = client.CreateOrUpdateBackupConfig(backupObj)
		require.NoError(t, err)
	})

	t.Run("should error while creating or updating backup configuration due to wrong headers", func(t *testing.T) {
		server := backupMockServer([]byte("backupJSON"), http.MethodPost, map[string]string{"Accept": gocd.HeaderVersionTwo, "Content-Type": gocd.ContentJSON})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/config/backup\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error while creating or updating backup configuration due to missing headers", func(t *testing.T) {
		server := backupMockServer([]byte("backupJSON"), http.MethodPost, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/config/backup\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error while creating or updating backup configuration as malformed data sent to server", func(t *testing.T) {
		server := backupMockServer([]byte("backupJSON"), http.MethodPost, correctBackupHeader)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		require.EqualError(t, err, "got 500 from GoCD while making POST call for "+server.URL+
			"/api/config/backup\nwith BODY:json: cannot unmarshal string into Go value of type gocd.BackupConfig")
	})

	t.Run("should error while creating or updating backup configuration as server was not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		backupObj := gocd.BackupConfig{
			EmailOnSuccess:   false,
			EmailOnFailure:   false,
			Schedule:         "0 0 2 * * ?",
			PostBackupScript: "/usr/local/bin/copy-gocd-backup-to-s3",
		}

		err := client.CreateOrUpdateBackupConfig(backupObj)
		require.EqualError(t, err, "call made to create/update backup configuration errored with: "+
			"Post \"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_DeleteBackupConfig(t *testing.T) {
	correctBackupHeader := map[string]string{"Accept": gocd.HeaderVersionOne}

	t.Run("should be able to delete the backup configurations successfully ", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK, correctBackupHeader, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteBackupConfig()
		require.NoError(t, err)
	})

	t.Run("should error out while deleting the backup configurations due to wrong headers", func(t *testing.T) {
		server := mockServer(nil, http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionTwo}, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		err := client.DeleteBackupConfig()
		require.EqualError(t, err, "got 404 from GoCD while making DELETE call for "+server.URL+
			"/api/config/backup\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
	})

	t.Run("should error out while deleting the backup configurations as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		err := client.DeleteBackupConfig()
		require.EqualError(t, err, "call made to delete backup configuration errored with: "+
			"Delete \"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
	})
}

func Test_client_GetBackup(t *testing.T) {
	correctArtifactHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to fetch the backup stats successfully", func(t *testing.T) {
		server := mockServer([]byte(backupStats), http.StatusOK,
			correctArtifactHeader, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.BackupStats{
			Time:           "2015-08-07T10:07:19.868Z",
			Path:           "/var/lib/go-server/serverBackups/backup_20150807-153719",
			Status:         "COMPLETED",
			ProgressStatus: "BACKUP_DATABASE",
			Message:        "Backup was generated successfully.",
		}

		actual, err := client.GetBackup("backup_20150807-153719")
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching backup stats present in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(backupStats), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.BackupStats{}

		actual, err := client.GetBackup("backup_20150807-153719")
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/backups/backup_20150807-153719\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching backup stats present in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(backupStats), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.BackupStats{}

		actual, err := client.GetBackup("backup_20150807-153719")
		require.EqualError(t, err, "got 404 from GoCD while making GET call for "+server.URL+
			"/api/backups/backup_20150807-153719\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching backup stats from GoCD as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte("backupStats"), http.StatusOK, correctArtifactHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		expected := gocd.BackupStats{}

		actual, err := client.GetBackup("backup_20150807-153719")
		require.EqualError(t, err, "reading response body errored with: invalid character 'b' looking for beginning of value")
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while fetching backup stats present in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		expected := gocd.BackupStats{}

		actual, err := client.GetBackup("backup_20150807-153719")
		require.EqualError(t, err, "call made to get backup stats errored with: "+
			"Get \"http://localhost:8156/go/api/backups/backup_20150807-153719\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, expected, actual)
	})
}

func Test_client_ScheduleBackup(t *testing.T) {
	scheduleBackupHeader := map[string]string{"Accept": gocd.HeaderVersionTwo}

	t.Run("should be able to schedule the backup successfully", func(t *testing.T) {
		server := mockServer([]byte(backupStats), http.StatusAccepted,
			scheduleBackupHeader, false,
			map[string]string{gocd.LocationHeader: "/var/lib/go-server/serverBackups/backup_20150807-153719", "Retry-After": "10"})

		expected := map[string]string{
			"BackUpID":   "backup_20150807-153719",
			"RetryAfter": "10",
		}

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ScheduleBackup()
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should error out while scheduling backup in GoCD due to wrong headers", func(t *testing.T) {
		server := mockServer([]byte(backupStats), http.StatusOK,
			map[string]string{"Accept": gocd.HeaderVersionThree}, false, nil)

		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ScheduleBackup()
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/backups\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while scheduling backup in GoCD due to missing headers", func(t *testing.T) {
		server := mockServer([]byte(backupStats), http.StatusOK,
			nil, false, nil)
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ScheduleBackup()
		require.EqualError(t, err, "got 404 from GoCD while making POST call for "+server.URL+
			"/api/backups\nwith BODY:<html>\n<body>\n\t<h2>404 Not found</h2>\n</body>\n\n</html>")
		assert.Nil(t, actual)
	})

	t.Run("should error out while scheduling backup in GoCD as headers missing from server response", func(t *testing.T) {
		server := mockServer([]byte("backupStats"), http.StatusAccepted, scheduleBackupHeader,
			false, map[string]string{"ETag": "cbc5f2d5b9c13a2cc1b1efb3d8a6155d"})
		client := gocd.NewClient(server.URL, auth, "info", nil)

		actual, err := client.ScheduleBackup()
		require.EqualError(t, err, "header Location not set, this will impact while getting backup stats")
		assert.Nil(t, actual)
	})

	t.Run("should error out while scheduling backup in GoCD as server is not reachable", func(t *testing.T) {
		client := gocd.NewClient("http://localhost:8156/go", auth, "info", nil)

		client.SetRetryCount(1)
		client.SetRetryWaitTime(1)

		actual, err := client.ScheduleBackup()
		require.EqualError(t, err, "call made to schedule backup errored with: Post "+
			"\"http://localhost:8156/go/api/backups\": dial tcp [::1]:8156: connect: connection refused")
		assert.Nil(t, actual)
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
