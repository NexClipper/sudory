package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Logger struct {
		Severity        string `env:"SUDORY_LOG_SEVERIY"           yaml:"severity,omitempty"`
		SystemEvent     bool   `env:"SUDORY_LOG_SYSTEM_EVENT"      yaml:"system-event,omitempty"`
		SystemEventName string `env:"SUDORY_LOG_SYSTEM_EVENT_NAME" yaml:"system-event-name,omitempty"`
		Verbose         bool   `env:"SUDORY_LOG_VERBOSE"           yaml:"verbose,omitempty"`
		Filename        string `env:"SUDORY_LOG_FILENAME"          yaml:"filename,omitempty"`
		MaxSize         int    `env:"SUDORY_LOG_MAXSIZE"           yaml:"max-size,omitempty"`
		MaxAge          int    `env:"SUDORY_LOG_MAXAGE"            yaml:"max-age,omitempty"`
		MaxBackups      int    `env:"SUDORY_LOG_MAXBACKUPS"        yaml:"max-backups,omitempty"`
		Compress        bool   `env:"SUDORY_LOG_COMPRESS"          yaml:"compress,omitempty"`
	}
}

var lazyinitlogger func(string) = func(string) {
	panic(errors.Errorf("call me after func init()"))
}

var onceinitlogger sync.Once

var LazyInitLogger func(string) = func(configfile string) {
	onceinitlogger.Do(func() {
		lazyinitlogger(configfile)
	})
}

var LoggerInfoOutput io.Writer = os.Stdout
var LoggerErrorOutput io.Writer = os.Stderr

// 로그 환경설정 초기화
func init() {

	//실행 명령에서 읽는다
	//logger config
	cfg := LoggerConfig{}

	flag.StringVar(&cfg.Logger.Severity, "log-severity", "debug", "severity of log severity=debug,[info|information],[warn|warning],error,fatal")
	flag.BoolVar(&cfg.Logger.SystemEvent, "log-system-event", false, "enabled system event")
	flag.StringVar(&cfg.Logger.SystemEventName, "log-system-eventname", "nexclipper.io/sudory", "system event name")
	flag.BoolVar(&cfg.Logger.Verbose, "log-verbose", false, "enabled verbose")

	//file rotator (for lumberjack)
	flag.StringVar(&cfg.Logger.Filename, "log-filename", "sudory.log", "log file name")
	flag.IntVar(&cfg.Logger.MaxSize, "log-max-size", 20, "maximum size in megabytes of the log file (MB)")
	flag.IntVar(&cfg.Logger.MaxAge, "log-max-age", 30, "maximum number of days to retain old log files (DAY)")
	flag.IntVar(&cfg.Logger.MaxBackups, "log-max-backups", 20, "maximum number of old log files to retain (COUNT)")
	flag.BoolVar(&cfg.Logger.Compress, "log-compress", false, "log compress option")

	//swap lazyinitlogger
	lazyinitlogger = func(configfile string) {

		//환경변수에서 읽는다 syscall.Getenv
		if err := configor.Load(&cfg, configfile); err != nil {
			fmt.Fprintln(os.Stderr, "ENV Unmarshal", err.Error())
		}

		rotate := &lumberjack.Logger{
			Filename:   cfg.Logger.Filename,
			MaxSize:    cfg.Logger.MaxSize,
			MaxAge:     cfg.Logger.MaxAge,
			MaxBackups: cfg.Logger.MaxBackups,
			Compress:   cfg.Logger.Compress,
		}

		LoggerInfoOutput = io.MultiWriter(rotate, LoggerInfoOutput)
		LoggerErrorOutput = io.MultiWriter(rotate, LoggerErrorOutput)

		logger.Init(cfg.Logger.SystemEventName, cfg.Logger.Verbose, cfg.Logger.SystemEvent, rotate)

		logger.SetLevel(logger.Level(severity(cfg.Logger.Severity)))

		//first log
		logger.Debugf("init logger%v",
			logs.KVL(
				"log-severity", cfg.Logger.Severity,
				"log-system-event", cfg.Logger.SystemEvent,
				"log-system-eventname", cfg.Logger.SystemEventName,
				"log-verbose", cfg.Logger.Verbose,
				"log-filename", cfg.Logger.Filename,
				"log-max-size", cfg.Logger.MaxSize,
				"log-max-age", cfg.Logger.MaxAge,
				"log-max-backups", cfg.Logger.MaxBackups,
				"log-compress", cfg.Logger.Compress,
			))

	}
}

func severity(s string) int {
	switch strings.ToLower(s) {
	case "debug":
		return 0
	case "info", "information":
		return 1
	case "warn", "warning":
		return 2
	case "error":
		return 3
	case "fatal":
		return 4
	default:
		return 0 //default
	}
}
