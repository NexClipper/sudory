package log

import (
	"log"
)

var glogger *Logger

type Logger struct {
	logger *log.Logger
}

func New() *Logger {
	if glogger != nil {
		return glogger
	}

	glogger = &Logger{logger: log.Default()}

	return glogger
}

func GetLogger() *Logger {
	if glogger == nil {
		return New()
	}

	return glogger
}

func Debugf(format string, v ...interface{}) {
	GetLogger().logger.Printf("[DEBUG] "+format, v...)
}

func Errorf(format string, v ...interface{}) {
	GetLogger().logger.Printf("[ERROR] "+format, v...)
}

func Fatalf(format string, v ...interface{}) {
	GetLogger().logger.Fatalf("[FATAL] "+format, v...)
}
