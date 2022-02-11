package log

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

var glogger *Logger

type level int

const (
	levelDebug level = 1 << iota
	levelInfo
	levelWarn
	levelError
)

type Logger struct {
	logger *log.Logger
	lv     level
	m      sync.RWMutex
}

func New(level string) *Logger {
	if glogger != nil {
		return glogger
	}

	glogger = &Logger{logger: log.Default()}

	glogger.SetLevel(level)

	return glogger
}

func (l *Logger) GetLevel() string {
	l.m.RLock()
	defer l.m.RUnlock()

	switch l.lv {
	case levelError | levelWarn | levelInfo | levelDebug:
		return "debug"
	case levelError | levelWarn | levelInfo:
		return "info"
	case levelError | levelWarn:
		return "warn"
	case levelError:
		return "error"
	}

	return ""
}

func (l *Logger) SetLevel(level string) error {
	lv := strings.ToLower(level)
	prevLv := l.GetLevel()

	if lv == prevLv {
		return fmt.Errorf("logger level(%s) you want to change is the same as the previous logger level(%s)", lv, prevLv)
	}

	l.m.Lock()
	defer l.m.Unlock()

	switch lv {
	case "debug":
		l.lv = levelError | levelWarn | levelInfo | levelDebug
	case "info":
		l.lv = levelError | levelWarn | levelInfo
	case "warn":
		l.lv = levelError | levelWarn
	case "error":
		l.lv = levelError
	}

	return nil
}

func GetLogger() *Logger {
	if glogger == nil {
		return New("debug")
	}

	return glogger
}

func Debugf(format string, v ...interface{}) {
	l := GetLogger()
	l.m.RLock()
	defer l.m.RUnlock()

	if l.lv&levelDebug != 0 {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	l := GetLogger()
	l.m.RLock()
	defer l.m.RUnlock()

	if l.lv&levelInfo != 0 {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	l := GetLogger()
	l.m.RLock()
	defer l.m.RUnlock()

	if l.lv&levelWarn != 0 {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	l := GetLogger()
	l.m.RLock()
	defer l.m.RUnlock()

	if l.lv&levelError != 0 {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	GetLogger().logger.Fatalf("[FATAL] "+format, v...)
}

// for go-retryablehttp
type RetryableHttpLogger struct{}

func (rhl *RetryableHttpLogger) Debug(msg string, keysAndValues ...interface{}) {
	msg, values := rhl.formatMsg(msg, keysAndValues...)
	Debugf(msg, values...)
}

func (rhl *RetryableHttpLogger) Info(msg string, keysAndValues ...interface{}) {
	msg, values := rhl.formatMsg(msg, keysAndValues...)
	Infof(msg, values...)
}

func (rhl *RetryableHttpLogger) Warn(msg string, keysAndValues ...interface{}) {
	msg, values := rhl.formatMsg(msg, keysAndValues...)
	Warnf(msg, values...)
}

func (rhl *RetryableHttpLogger) Error(msg string, keysAndValues ...interface{}) {
	msg, values := rhl.formatMsg(msg, keysAndValues...)
	Errorf(msg, values...)
}

func (rhl *RetryableHttpLogger) formatMsg(msg string, keysAndValues ...interface{}) (string, []interface{}) {
	if len(keysAndValues) == 0 {
		return msg, nil
	}

	msg += ": { "
	var values []interface{}
	for i, v := range keysAndValues {
		if i%2 == 0 {
			msg += v.(string) + ": { %v }"
		} else {
			values = append(values, v)
			if len(keysAndValues)-1 != i {
				msg += ", "
			}
		}
	}
	msg += " }"

	return msg, values
}
