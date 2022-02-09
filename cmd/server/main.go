package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/events"
	"github.com/NexClipper/sudory/pkg/server/macro/channels"
	"github.com/NexClipper/sudory/pkg/server/route"
	"github.com/NexClipper/sudory/pkg/server/status"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"xorm.io/xorm"
)

func main() {
	cfg := &config.Config{}
	flag.StringVar(&cfg.Database.Host, "db-host", "127.0.0.1", "Database's host")
	flag.StringVar(&cfg.Database.Port, "db-port", "3306", "Database's port")
	flag.StringVar(&cfg.Database.Username, "db-user", "", "Database's username")
	flag.StringVar(&cfg.Database.Password, "db-passwd", "", "Database's password")
	flag.StringVar(&cfg.Database.DBName, "db-dbname", "", "Database's dbname")

	configPath := flag.String("config", "../../conf/sudory-server.yml", "Path to sudory-server's config file")
	flag.Parse()

	cfg, err := config.New(cfg, *configPath)
	if err != nil {
		panic(err)
	}

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//events
	deactivate := event(*configPath)
	defer func() {
		deactivate() //stop when closing
	}()

	//chrons
	chronStop := chron(db.Engine())
	defer func() {
		chronStop()
	}()

	r := route.New(cfg, db)

	r.Start(cfg.Host.Port)
}

func event(filename string) func() {
	var err error
	//events
	var contexts []events.EventContexter
	var config *events.Config
	//event config
	if config, err = events.NewConfig(filename); err != nil { //config file load
		panic(err)
	}
	//event config vaild
	if err = config.Vaild(); err != nil { //config vaild
		panic(err)
	}
	//event config make listener
	if contexts, err = config.MakeEventListener(); err != nil { //events regist listener
		panic(err)
	}
	//event manager
	sender := channels.NewSafeChannel(0)
	manager := events.NewManager(sender, contexts, log.Printf)
	deactivate := events.Activate(manager, len(contexts)) //manager activate

	events.Invoke = manager.Invoker //setting invoker

	return deactivate
}

func chron(engine *xorm.Engine) func() {
	const ChronInterval = 10
	chronStop := status.NewChron(os.Stdout, ChronInterval*time.Second,
		//chron environment
		func() status.ChronUpdater {
			inst := env.NewEnvironmentChron(database.NewContext(engine))
			//환경변수 리스트 검사
			if err := inst.(*env.EnvironmentChron).WhiteListCheck(); err != nil {
				log.Println(err)
				log.Println("merge environment setting") //환경변수 병합
				if err := inst.(*env.EnvironmentChron).Merge(); err != nil {
					panic(err)
				}
			}
			//환경변수 업데이트
			if err := inst.Update(); err != nil {
				panic(err)
			}
			return inst
		},
	)

	return chronStop
}
