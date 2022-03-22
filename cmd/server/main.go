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
	eventCancel, err := newEvent(*configPath)
	if err != nil {
		panic(err)
	}
	defer eventCancel() //이벤트 종료

	//chrons
	chronStop := chron(db.Engine())
	defer func() {
		chronStop()
	}()

	r := route.New(cfg, db)

	r.Start(cfg.Host.Port)
}

func newEvent(filename string) (func(), error) {
	cfgevent, err := event.NewEventConfig(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "make new event config")
	}

	pub := event.NewEventPublish()

	for i := range cfgevent.EventSubscribeConfigs {
		cfgsub := cfgevent.EventSubscribeConfigs[i]

		sub := event.NewEventSub(cfgsub)
		sub.ErrorHandlers().Add(func(sub event.EventSubscriber, err error) {

			var stack string
			logs.CauseIter(err, func(err error) {
				logs.StackIter(err, func(s string) {
					stack = s
				})
			})

			logger.Error(fmt.Errorf("%w %s", err,
				logs.KVL(
					"stack", stack,
				)))
		})

		if err := event.RegistNotifier(sub); err != nil {
			return nil, errors.Wrapf(err, "regist notifier")
		}

		sub.Regist(pub)

	}
	event.PrintEventConfiguation(os.Stdout, pub)

	//return closer
	return pub.Close, nil
}

func chron(engine *xorm.Engine) func() {
	const ChronInterval = 10
	chronStop := status.NewChron(os.Stdout, ChronInterval*time.Second,
		//chron environment
		func() status.ChronUpdater {
			sink := logs.WithName("Chron Environment")
			updator := env.NewEnvironmentChron(database.NewXormContext(engine.NewSession()))
			//환경변수 리스트 검사
			if err := updator.WhiteListCheck(); err != nil {
				logger.Info(sink.WithError(errors.Wrapf(err, "WhiteListCheck")).String())
				logger.Debugln("merge environment setting") //환경변수 병합
				if err := updator.Merge(); err != nil {
					logger.Error(sink.WithError(errors.Wrapf(err, "Merge")).String())
				}
			}
			//환경변수 업데이트
			if err := updator.Update(); err != nil {
				logger.Fatal(err) //first work error //(exit)
			}
			return updator
		},
	)

	return chronStop
}
