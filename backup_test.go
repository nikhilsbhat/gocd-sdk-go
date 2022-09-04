package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetBackupInfo(t *testing.T) {
	t.Run("", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			nil,
			"info",
		)

		backup, err := client.GetBackupInfo()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(backup.Schedule))
	})
}
