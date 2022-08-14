package logger

type logger struct {
	module     string
	mainLogger *commonLogger
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.mainLogger.Infof(l.module, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.mainLogger.Errorf(l.module, format, args...)
}

func (l *logger) Error(err error) {
	l.Errorf(`%s`, err)
}

func Make(module string) *logger {
	return &logger{
		module:     module,
		mainLogger: mainLogger,
	}
}
