package logger_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_getLoglevel(t *testing.T) {
	t.Run("should return warn level", func(t *testing.T) {
		actual := logger.GetLoglevel("warning")
		assert.Equal(t, log.WarnLevel, actual)
	})
	t.Run("should return trace level", func(t *testing.T) {
		actual := logger.GetLoglevel("trace")
		assert.Equal(t, log.TraceLevel, actual)
	})
	t.Run("should return debug level", func(t *testing.T) {
		actual := logger.GetLoglevel("debug")
		assert.Equal(t, log.DebugLevel, actual)
	})
	t.Run("should return fatal level", func(t *testing.T) {
		actual := logger.GetLoglevel("fatal")
		assert.Equal(t, log.FatalLevel, actual)
	})
	t.Run("should return error level", func(t *testing.T) {
		actual := logger.GetLoglevel("error")
		assert.Equal(t, log.ErrorLevel, actual)
	})
	t.Run("should return error level", func(t *testing.T) {
		actual := logger.GetLoglevel("")
		assert.Equal(t, log.InfoLevel, actual)
	})
}
