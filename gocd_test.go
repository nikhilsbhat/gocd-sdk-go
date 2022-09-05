package gocd

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_getLoglevel(t *testing.T) {
	t.Run("should return warn level", func(t *testing.T) {
		actual := getLoglevel("warning")
		assert.Equal(t, log.WarnLevel, actual)
	})
	t.Run("should return trace level", func(t *testing.T) {
		actual := getLoglevel("trace")
		assert.Equal(t, log.TraceLevel, actual)
	})
	t.Run("should return debug level", func(t *testing.T) {
		actual := getLoglevel("debug")
		assert.Equal(t, log.DebugLevel, actual)
	})
	t.Run("should return fatal level", func(t *testing.T) {
		actual := getLoglevel("fatal")
		assert.Equal(t, log.FatalLevel, actual)
	})
	t.Run("should return error level", func(t *testing.T) {
		actual := getLoglevel("error")
		assert.Equal(t, log.ErrorLevel, actual)
	})
}
