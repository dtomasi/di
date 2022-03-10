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

type Service struct {
	logger logr.Logger
}

func (s *Service) GetUnderlying() logr.Logger {
	return s.logger
}

func (s *Service) Log(lvl Level, msg string, keysAndValues ...interface{}) {
	s.logger.V(int(lvl)).Info(msg, keysAndValues)
}

func (s *Service) Warning(msg string, keysAndValues ...interface{}) {
	s.Log(Warning, msg, keysAndValues)
}

func (s *Service) Notice(msg string, keysAndValues ...interface{}) {
	s.Log(Notice, msg, keysAndValues)
}

func (s *Service) Info(msg string, keysAndValues ...interface{}) {
	s.Log(Info, msg, keysAndValues)
}

func (s *Service) Debug(msg string, keysAndValues ...interface{}) {
	s.Log(Debug, msg, keysAndValues)
}

func (s *Service) Error(err error, msg string, keysAndValues ...interface{}) {
	s.logger.Error(err, msg, keysAndValues)
}
