package gocd_test

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

//go:embed internal/fixtures/backup.json
var backupJSON string

func TestConfig_GetBackupInfo(t *testing.T) {
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

		actual, err := client.GetBackupInfo()
		assert.EqualError(t, err, "call made to get backup information errored with "+
			"Get \"http://localhost:8156/go/api/config/backup\": dial tcp [::1]:8156: connect: connection refused")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should error out while fetching latest backup configuration information as server returned non 200 status code", func(t *testing.T) {
		server := mockServer([]byte("backupJSON"), http.StatusBadGateway, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetBackupInfo()
		assert.EqualError(t, err, gocd.APIWithCodeError(http.StatusBadGateway).Error())
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should error out while fetching latest backup configuration information as server returned malformed response", func(t *testing.T) {
		server := mockServer([]byte(`{"email_on_failure"}`), http.StatusOK, nil)
		client := gocd.NewClient(
			server.URL,
			"admin",
			"admin",
			"info",
			nil,
		)

		actual, err := client.GetBackupInfo()
		assert.EqualError(t, err, "reading response body errored with: invalid character '}' after object key")
		assert.Equal(t, gocd.BackupConfig{}, actual)
	})

	t.Run("should be able to fetch the latest backup configuration information available in GoCD", func(t *testing.T) {
		server := mockServer([]byte(backupJSON), http.StatusOK, nil)

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

		actual, err := client.GetBackupInfo()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
