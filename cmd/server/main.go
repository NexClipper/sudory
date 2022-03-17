package main

import (
	"bytes"
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
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

func init() {
	println("init timezone UTC")
	time.Local = time.UTC //set timezone UTC
}

func main() {
	cfg := &config.Config{}
	flag.StringVar(&cfg.Database.Host, "db-host", "127.0.0.1", "Database's host")
	flag.StringVar(&cfg.Database.Port, "db-port", "3306", "Database's port")
	flag.StringVar(&cfg.Database.Username, "db-user", "", "Database's username")
	flag.StringVar(&cfg.Database.Password, "db-passwd", "", "Database's password")
	flag.StringVar(&cfg.Database.DBName, "db-dbname", "", "Database's dbname")

	configPath := flag.String("config", "../../conf/sudory-server.yml", "Path to sudory-server's config file")
	flag.Parse()
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

			err_ := errors.Cause(err)

			type stackTracer interface {
				StackTrace() errors.StackTrace
			}

			sink := logs.WithError(err)
			sink = sink.WithValue(
				"event-name", sub.Config().Name,
			)

			buff := &bytes.Buffer{}
			if err, ok := err_.(stackTracer); ok {
				for _, f := range err.StackTrace() {
					fmt.Fprintf(buff, "%+s:%d\n", f, f)
				}

				sink = sink.WithValue(
					"stack", buff.Bytes(),
				)
			}

			logger.Errorf("%w %s", err, sink.String())
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
