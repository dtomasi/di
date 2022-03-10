package logger

import "github.com/go-logr/logr"

func NewService(logger logr.Logger) *Service {
	return &Service{logger: logger}
}
