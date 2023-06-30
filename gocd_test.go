package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_getLoglevel(t *testing.T) {
	t.Run("should return warn level", func(t *testing.T) {
		actual := gocd.GetLoglevel("warning")
		assert.Equal(t, log.WarnLevel, actual)
	})
	t.Run("should return trace level", func(t *testing.T) {
		actual := gocd.GetLoglevel("trace")
		assert.Equal(t, log.TraceLevel, actual)
	})
	t.Run("should return debug level", func(t *testing.T) {
		actual := gocd.GetLoglevel("debug")
		assert.Equal(t, log.DebugLevel, actual)
	})
	t.Run("should return fatal level", func(t *testing.T) {
		actual := gocd.GetLoglevel("fatal")
		assert.Equal(t, log.FatalLevel, actual)
	})
	t.Run("should return error level", func(t *testing.T) {
		actual := gocd.GetLoglevel("error")
		assert.Equal(t, log.ErrorLevel, actual)
	})
}

func TestGetGoCDMethodNames(t *testing.T) {
	t.Run("should list all method names", func(t *testing.T) {
		response := gocd.GetGoCDMethodNames()
		assert.Equal(t, 131, len(response))
		assert.Equal(t, "AgentKillTask", response[0])
		assert.Equal(t, "ValidatePipelineSyntax", response[130])
	})
}
