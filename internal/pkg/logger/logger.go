package logger

import (
	"github.com/go-logr/logr"
)

//go:generate stringer -type=Level -output=zz_gen_level_string.go

type Level int

// Apply common unix syslog levels.
const (
	Warning Level = iota
	Notice
	Info
	Debug
)

const loggerName = "DI"

func NewInternalLogger(logger logr.Logger) *InternalLogger {
	return &InternalLogger{logger: logger.WithName(loggerName)}
}

type InternalLogger struct {
	logger logr.Logger
}

func (s *InternalLogger) GetUnderlying() logr.Logger {
	return s.logger
}

func (s *InternalLogger) Log(lvl Level, msg string, keysAndValues ...interface{}) {
	s.logger.V(int(lvl)).Info(msg, keysAndValues)
}

func (s *InternalLogger) Warning(msg string, keysAndValues ...interface{}) {
	s.Log(Warning, msg, keysAndValues)
}

func (s *InternalLogger) Notice(msg string, keysAndValues ...interface{}) {
	s.Log(Notice, msg, keysAndValues)
}

func (s *InternalLogger) Info(msg string, keysAndValues ...interface{}) {
	s.Log(Info, msg, keysAndValues)
}

func (s *InternalLogger) Debug(msg string, keysAndValues ...interface{}) {
	s.Log(Debug, msg, keysAndValues)
}

func (s *InternalLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	s.logger.Error(err, msg, keysAndValues)
}
