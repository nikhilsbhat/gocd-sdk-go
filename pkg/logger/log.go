package logger

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetLoglevel sets the loglevel to the kind of log asked for.
func GetLoglevel(level string) log.Level {
	switch strings.ToLower(level) {
	case log.WarnLevel.String():
		return log.WarnLevel
	case log.DebugLevel.String():
		return log.DebugLevel
	case log.TraceLevel.String():
		return log.TraceLevel
	case log.FatalLevel.String():
		return log.FatalLevel
	case log.ErrorLevel.String():
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
