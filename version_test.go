package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

func Test_config_GetVersionInfo(t *testing.T) {
	t.Run("should be able to fetch the version info", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			nil,
			"info",
		)
		actual, err := client.GetVersionInfo()
		assert.NoError(t, err)
		assert.Equal(t, gocd.VersionInfo{}, actual)
	})
}
