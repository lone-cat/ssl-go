package logger

import (
	"log"
)

type severity uint8

const (
	undefined severity = iota
	errorSeverity
	infoSeverity
)

func (s severity) String() string {
	switch s {
	case errorSeverity:
		return `ERROR`
	case infoSeverity:
		return `INFO`
	default:
		return `undefined`
	}
}

type commonLogger struct {
	loggers map[severity]*log.Logger
}

var mainLogger *commonLogger

func (l *commonLogger) log(module string, severity severity, format string, args ...interface{}) {
	prefix := ``
	realArgs := make([]interface{}, 0)
	if module != `` {
		prefix += `[%s] `
		realArgs = append(realArgs, module)
	}
	if severity != undefined {
		prefix += `[%s] `
		realArgs = append(realArgs, severity.String())

	}
	loggr, exist := l.loggers[severity]
	if !exist || loggr == nil {
		return
	}
	realArgs = append(realArgs, args...)
	loggr.Printf(prefix+format, realArgs...)
}

func (l *commonLogger) Infof(module string, format string, args ...interface{}) {
	l.log(module, infoSeverity, format, args...)
}

func (l *commonLogger) Errorf(module string, format string, args ...interface{}) {
	l.log(module, errorSeverity, format, args...)
}
