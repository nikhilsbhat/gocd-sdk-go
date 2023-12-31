package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestGetGoCDMethodNames(t *testing.T) {
	t.Run("should list all method names", func(t *testing.T) {
		response := gocd.GetGoCDMethodNames()
		assert.Equal(t, 142, len(response))
		assert.Equal(t, "AgentKillTask", response[0])
		assert.Equal(t, "UpdatePluginSettings", response[136])
	})
}
