package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/route"
	"github.com/NexClipper/sudory/pkg/server/status"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/NexClipper/sudory/pkg/version"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

const APP_NAME = "sudory-server"

func init() {
	println("init timezone UTC")
	time.Local = time.UTC //set timezone UTC
}

func main() {
	versionFlag := flag.Bool("version", false, "print the current version")

	cfg := &config.Config{}
	flag.StringVar(&cfg.Database.Host, "db-host", "127.0.0.1", "Database's host")
	flag.StringVar(&cfg.Database.Port, "db-port", "3306", "Database's port")
	flag.StringVar(&cfg.Database.Username, "db-user", "", "Database's username")
	flag.StringVar(&cfg.Database.Password, "db-passwd", "", "Database's password")
	flag.StringVar(&cfg.Database.DBName, "db-dbname", "", "Database's dbname")

	configPath := flag.String("config", "../../conf/sudory-server.yml", "Path to sudory-server's config file")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version.BuildVersion(APP_NAME))
		return
	}

	config.LazyInitLogger(*configPath) //init logger

	cfg, err := config.New(cfg, *configPath)
	if err != nil {
		panic(err)
	}

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//init event
	eventClose, err := newEvent(*configPath)
	if err != nil {
		panic(err)
	}
	defer eventClose() //이벤트 종료

	//init cron
	cronClose, err := newCron(db.Engine())
	if err != nil {
		panic(err)
	}
	defer cronClose() //크론잡 종료

	r := route.New(cfg, db)

	r.Start(cfg.Host.Port)
}

func newEvent(filename string) (func(), error) {
	//에러 핸들러 등록
	errorHandlers := event.HashsetErrorHandlers{}
	errorHandlers.Add(func(err error) {

		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = s
			})
		})

		logger.Error(fmt.Errorf("event notify: %w %s", err,
			logs.KVL(
				"stack", stack,
			)))
	})

	cfgevent, err := event.NewEventConfig(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "make new event config")
	}

	pub := event.NewEventPublish()

	for i := range cfgevent.EventSubscribeConfigs {
		cfgsub := cfgevent.EventSubscribeConfigs[i]

		sub := event.NewEventSubscribe(cfgsub, errorHandlers)

		if err := event.RegistNotifier(sub); err != nil {
			return nil, errors.Wrapf(err, "regist notifier")
		}

		sub.Regist(pub)

	}
	event.PrintEventConfiguation(os.Stdout, pub)

	//return closer
	return pub.Close, nil
}

func newCron(engine *xorm.Engine) (func(), error) {
	const interval = 10 * time.Second

	//환경설정 updater 생성
	envUpdater, err := newEnvironmentUpdate(engine)
	if err != nil {
		return nil, errors.Wrapf(err, "create environment cron updater")
	}

	//에러 핸들러 등록
	errorHandlers := status.HashsetErrorHandlers{}
	errorHandlers.Add(func(err error) {
		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = s
			})
		})

		logger.Error(fmt.Errorf("cron jobs: %w %s", err,
			logs.KVL(
				"stack", stack,
			)))
	})

	//new ticker
	tickerClose := status.NewTicker(interval,
		//environment update
		func() {
			if err := envUpdater.Update(); err != nil {
				errorHandlers.OnError(errors.Wrapf(err, "environment update"))
			}
		},
	)

	return tickerClose, nil
}

func newEnvironmentUpdate(engine *xorm.Engine) (*env.EnvironmentUpdate, error) {
	updator := env.NewEnvironmentUpdate(database.NewXormContext(engine.NewSession()))
	//환경변수 리스트 검사
	if err := updator.WhiteListCheck(); err != nil {
		//빠져있는 환경변수 추가
		if err := updator.Merge(); err != nil {
			return nil, errors.Wrapf(err, "merge")
		}
	}
	//환경변수 업데이트
	if err := updator.Update(); err != nil {
		return nil, errors.Wrapf(err, "update")
	}

	return updator, nil
}
