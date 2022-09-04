package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetConfigRepoInfo(t *testing.T) {
	t.Run("should be able retrieve config repo information", func(t *testing.T) {
		client := NewClient(
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
