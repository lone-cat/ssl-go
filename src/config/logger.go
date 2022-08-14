package config

import (
	"fmt"
)

type loggerInterface interface {
	Infof(format string, args ...interface{})
}

type defaultLogger struct{}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

var logger loggerInterface

func init() {
	SetLogger(&defaultLogger{})
}

func SetLogger(loggr loggerInterface) {
	logger = loggr
}
