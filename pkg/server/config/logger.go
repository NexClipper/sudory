package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/jinzhu/configor"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Severity        string `env:"SUDORY_LOG_SEVERIY"`
	SystemEvent     bool   `env:"SUDORY_LOG_SYSTEM_EVENT"`
	SystemEventName string `env:"SUDORY_LOG_SYSTEM_EVENT_NAME"`
	Verbose         bool   `env:"SUDORY_LOG_VERBOSE"`
	VerboseLevel    int    `env:"SUDORY_LOG_VERBOSELEVEL"`
	Filename        string `env:"SUDORY_LOG_FILENAME"`
	MaxSize         int    `env:"SUDORY_LOG_MAXSIZE"`
	MaxAge          int    `env:"SUDORY_LOG_MAXAGE"`
	MaxBackups      int    `env:"SUDORY_LOG_MAXBACKUPS"`
	Compress        bool   `env:"SUDORY_LOG_COMPRESS"`
}

var lazyinitlogger func() = func() {
	panic(fmt.Errorf("call me after func init()"))
}

var onceinitlogger sync.Once

var LazyInitLogger func() = func() {
	onceinitlogger.Do(func() {
		lazyinitlogger()
	})
}

// 로그 환경설정 초기화
func init() {

	//실행 명령에서 읽는다
	//logger config
	cfg := LoggerConfig{}

	flag.StringVar(&cfg.Severity, "log-severity", "debug", "severity of log severity=debug,[info|information],[warn|warning],error,fatal")
	flag.BoolVar(&cfg.SystemEvent, "log-system-event", true, "enabled system event")
	flag.StringVar(&cfg.SystemEventName, "log-system-eventname", "nexclipper.io/sudory", "system event name")
	flag.BoolVar(&cfg.Verbose, "log-verbose", false, "enabled verbose")
	flag.IntVar(&cfg.VerboseLevel, "log-verbose-level", 0, "verbose level higher more detail max=5")

	//file rotator (for lumberjack)
	flag.StringVar(&cfg.Filename, "log-filename", "sudory.log", "log file name")
	flag.IntVar(&cfg.MaxSize, "log-max-size", 20, "maximum size in megabytes of the log file (MB)")
	flag.IntVar(&cfg.MaxAge, "log-max-age", 30, "maximum number of days to retain old log files (DAY)")
	flag.IntVar(&cfg.MaxBackups, "log-max-backups", 20, "maximum number of old log files to retain (COUNT)")
	flag.BoolVar(&cfg.Compress, "log-compress", false, "log compress option")

	//swap lazyinitlogger
	lazyinitlogger = func() {

		//환경변수에서 읽는다 syscall.Getenv
		if err := configor.Load(&cfg); err != nil {
			fmt.Fprintln(os.Stderr, "ENV Unmarshal", err.Error())
		}

		rotate := lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			Compress:   cfg.Compress,
		}

		logger.Init(cfg.SystemEventName, cfg.SystemEvent, cfg.Verbose, &rotate)

		logger.SetLevel(logger.Level(severity(cfg.Severity)))

		logs.SetVerbose(cfg.VerboseLevel)

		//first log
		logger.Debugln(logs.WithName("init logger").WithValue(
			"log-severity", cfg.Severity, "log-system-event", cfg.SystemEvent, "log-system-eventname", cfg.SystemEventName, "log-verbose", cfg.Verbose,
			"log-filename", cfg.Filename, "log-max-size", cfg.MaxSize, "log-max-age", cfg.MaxAge, "log-max-backups", cfg.MaxBackups, "log-compress", cfg.Compress).String())
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
