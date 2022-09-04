package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should be able retrieve config repo information", func(t *testing.T) {
		client := gocd.NewClient(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			nil,
			"info",
		)

		repos, err := client.GetConfigRepoInfo()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(repos))
	})
}
