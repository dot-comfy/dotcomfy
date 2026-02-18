package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var l *log.Logger

func Init(verbosity int) {
	if verbosity == 0 {
		// Disable logging by default - redirect to /dev/null
		nullFile, _ := os.OpenFile("/dev/null", os.O_WRONLY, 0644)
		l = log.NewWithOptions(nullFile, log.Options{
			ReportCaller:    false,
			ReportTimestamp: false,
		})
		l.SetLevel(log.FatalLevel) // Only fatal logs would show, but we redirect to null
	} else {
		l = log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
		})
		if verbosity > 4 {
			verbosity = 4
		}
		switch verbosity {
		case 4:
			l.SetLevel(log.DebugLevel)
		case 3:
			l.SetLevel(log.InfoLevel)
		case 2:
			l.SetLevel(log.WarnLevel)
		default: // 1
			l.SetLevel(log.ErrorLevel)
		}
	}
}

func GetLogger() *log.Logger {
	return l
}
