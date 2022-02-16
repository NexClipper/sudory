package logs

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/macro/logs/internal/serialize"
)

func MKV(msg string, keysAndValues ...interface{}) string {

	buf := bytes.Buffer{}

	buf.WriteString(msg)

	serialize.KVListFormat(&buf, keysAndValues...)

	return buf.String()
}

func PrintS(msg string, keysAndValues ...interface{}) string {

	buf := bytes.Buffer{}

	buf.WriteString(strconv.Quote(msg))

	serialize.KVListFormat(&buf, keysAndValues...)

	return buf.String()
}

func PrintE(err error, msg string, keysAndValues ...interface{}) string {

	buf := bytes.Buffer{}

	buf.WriteString(strconv.Quote(msg))

	serialize.KVListFormat(&buf, "err", err)

	serialize.KVListFormat(&buf, keysAndValues...)

	return buf.String()
}

func DebugS(msg string, keysAndValues ...interface{}) {
	logger.DebugDepth(1, PrintS(msg, keysAndValues...))
}
func Debug(v ...interface{}) {
	logger.DebugDepth(1, fmt.Sprint(v...))
}
func Debugln(v ...interface{}) {
	logger.DebugDepth(1, fmt.Sprintln(v...))
}
func DebugDepth(depth int, v ...interface{}) {
	logger.DebugDepth(depth+1, v...)
}
func Debugf(format string, v ...interface{}) {
	logger.DebugDepth(1, fmt.Sprintf(format, v...))
}

func InfoS(msg string, keysAndValues ...interface{}) {
	logger.InfoDepth(1, PrintS(msg, keysAndValues...))
}
func Info(v ...interface{}) {
	logger.InfoDepth(1, fmt.Sprint(v...))
}
func Infoln(v ...interface{}) {
	logger.InfoDepth(1, fmt.Sprintln(v...))
}
func InfoDepth(depth int, v ...interface{}) {
	logger.InfoDepth(depth+1, fmt.Sprint(v...))
}
func Infof(format string, v ...interface{}) {
	logger.InfoDepth(1, fmt.Sprintf(format, v...))
}

func WarningS(msg string, keysAndValues ...interface{}) {
	logger.WarningDepth(1, PrintS(msg, keysAndValues...))
}
func Warning(v ...interface{}) {
	logger.WarningDepth(1, fmt.Sprint(v...))
}
func Warningln(v ...interface{}) {
	logger.WarningDepth(1, fmt.Sprintln(v...))
}
func WarningDepth(depth int, v ...interface{}) {
	logger.WarningDepth(depth+1, fmt.Sprint(v...))
}
func Warningf(format string, v ...interface{}) {
	logger.WarningDepth(1, fmt.Sprintf(format, v...))
}

func ErrorS(err error, msg string, keysAndValues ...interface{}) {
	logger.ErrorDepth(1, PrintE(err, msg, keysAndValues...))
}
func Error(v ...interface{}) {
	logger.ErrorDepth(1, fmt.Sprint(v...))
}
func Errorln(v ...interface{}) {
	logger.ErrorDepth(1, fmt.Sprintln(v...))
}
func ErrorDepth(depth int, v ...interface{}) {
	logger.ErrorDepth(depth+1, fmt.Sprint(v...))
}
func Errorf(format string, v ...interface{}) {
	logger.ErrorDepth(1, fmt.Sprintf(format, v...))
}

func FatalS(msg string, keysAndValues ...interface{}) {
	logger.FatalDepth(1, PrintS(msg, keysAndValues...))
}
func Fatal(v ...interface{}) {
	logger.FatalDepth(1, fmt.Sprint(v...))
}
func Fatalln(v ...interface{}) {
	logger.FatalDepth(1, fmt.Sprintln(v...))
}
func FatalDepth(depth int, v ...interface{}) {
	logger.FatalDepth(depth+1, fmt.Sprint(v...))
}
func Fatalf(format string, v ...interface{}) {
	logger.FatalDepth(1, fmt.Sprintf(format, v...))
}
