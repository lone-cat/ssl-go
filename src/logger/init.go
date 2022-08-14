package logger

import (
	"log"
	"os"
)

func init() {
	loggers := make(map[severity]*log.Logger)
	loggers[errorSeverity] = log.New(os.Stderr, ``, log.LstdFlags)
	loggers[infoSeverity] = log.New(os.Stdout, ``, log.LstdFlags)
	mainLogger = &commonLogger{
		loggers: loggers,
	}
}
