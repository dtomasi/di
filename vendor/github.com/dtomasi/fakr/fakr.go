package fakr

import (
	"github.com/go-logr/logr"
)

func New() logr.Logger {
	return logr.New(&fakeSink{}) //nolint:exhaustivestruct
}

type fakeSink struct {
	sink logr.LogSink // nolint:structcheck,unused
}

func (s fakeSink) Init(_ logr.RuntimeInfo) {

}

func (s fakeSink) Enabled(_ int) bool {
	return true
}
func (s fakeSink) Info(_ int, _ string, _ ...interface{}) {

}
func (s fakeSink) Error(_ error, _ string, _ ...interface{}) {

}
func (s fakeSink) WithValues(_ ...interface{}) logr.LogSink {
	return s
}
func (s fakeSink) WithName(_ string) logr.LogSink {
	return s
}
