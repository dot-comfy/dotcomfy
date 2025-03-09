package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var l *log.Logger

func Init(verbosity int) {
	l = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
	})
	if verbosity > 3 {
		verbosity = 3
	}
	switch verbosity {
	case 3:
		l.SetLevel(log.DebugLevel)
	case 2:
		l.SetLevel(log.InfoLevel)
	case 1:
		l.SetLevel(log.WarnLevel)
	default:
		l.SetLevel(log.ErrorLevel)
	}
}

func GetLogger() *log.Logger {
	return l
}
