package main

import (
	"database/sql"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/conf/script/migrations"
	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/database"
	flavor "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/resolvers/mysql"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/event/managed_channel"
	"github.com/NexClipper/sudory/pkg/server/event/managed_event"
	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	dctv1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/NexClipper/sudory/pkg/server/route"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/NexClipper/sudory/pkg/server/status/ticker"
	"github.com/NexClipper/sudory/pkg/version"
	"github.com/fsnotify/fsnotify"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
	"xorm.io/builder"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
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

	enigmaConfigFilename := cfg.Encryption
	if !path.IsAbs(cfg.Encryption) {
		enigmaConfigFilename = path.Join(path.Dir(*configPath), cfg.Encryption)
	}
	if err := newEnigma(enigmaConfigFilename); err != nil {
		panic(err)
	}

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := doMigration(cfg); err != nil {
		panic(err)
	}

	if true {
		//init event
		eventsConfigYaml := cfg.Events
		if !path.IsAbs(cfg.Events) {
			eventsConfigYaml = path.Join(path.Dir(*configPath), cfg.Events)
		}
		eventClose, err := newEvent(eventsConfigYaml)
		if err != nil {
			panic(err)
		}
		defer eventClose() //이벤트 종료
	}
	if true {
		//init managed event
		me := managed_event.NewManagedEvent()
		me.SetEngine(db.Engine())
		me.ErrorHandlers.Add(managed_event.DefaultErrorHandler)
		me.NofitierErrorHandlers.Add(
			managed_event.DefaultErrorHandler_notifier(me),
			func(n managed_event.Notifier, err error) { managed_event.DefaultErrorHandler(err) },
		)
		managed_event.Invoke = me.Invoke
	}

	//init managed channel
	mc := managed_channel.NewEvent(db.Engine().DB().DB, flavor.Dialect())

	mc.ErrorHandlers.Add(managed_event.DefaultErrorHandler)
	mc.NofitierErrorHandlers.Add(
		managed_channel.DefaultErrorHandler_notifier(mc),
		func(n managed_channel.Notifier, err error) { managed_channel.DefaultErrorHandler(err) },
	)
	managed_channel.InvokeByChannelUuid = mc.InvokeByChannelUuid
	managed_channel.InvokeByEventCategory = mc.InvokeByEventCategory

	//init global variant cron
	cronGVClose, err := newGlobalVariablesCron(db.Engine().DB().DB)
	if err != nil {
		panic(err)
	}
	defer cronGVClose() //크론잡 종료

	if false {
		//init purge deleted service
		cronPurgeServiceClose, err := newPurgeDeletedDataCron(db.Engine(), cfg.RespitePeriod)
		if err != nil {
			panic(err)
		}
		defer cronPurgeServiceClose() //크론잡 종료
	}

	r := route.New(cfg, db)
	r.Start(cfg.Host.Port)

	logger.Debugf("%s is DONE", path.Base(strings.ReplaceAll(os.Args[0], "\\", "/")))
}

func newEvent(filename string) (closer func(), err error) {
	//에러 핸들러 등록
	errorHandlers := event.HashsetErrorHandlers{}
	errorHandlers.Add(func(err error) {

		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = logs.KVL(
					"stack", s,
				)
			})
		})

		logger.Error(fmt.Errorf("event notify: %w%s", err, stack))
	})

	hasing := func(filename string) (uint32, error) {

		var (
			buf = make([]byte, 1<<12) //4k
			h   = crc32.NewIEEE()
			// h = sha1.New()
		)

		fd, err := os.Open(filename)
		if err != nil {
			return 0, err
		}
		defer fd.Close()

		if _, err := io.CopyBuffer(h, fd, buf); err != nil {
			return 0, err
		}

		return h.Sum32(), nil
	}

	var (
		oldhash uint32
		oldpub  event.EventPublisher
	)

	setEvent := func() error {
		//read config file
		cfgevent, err := event.NewEventConfig(filename)
		if err != nil {
			return errors.Wrapf(err, "make new event config")
		}
		//new event publisher
		newpub := event.NewEventPublish()
		//regist subscriber
		for _, cfgsub := range cfgevent.EventSubscribeConfigs {
			sub := event.NewEventSubscribe(cfgsub, errorHandlers)

			for _, cfgnotifier := range cfgsub.NotifierConfigs {

				//new notifier
				notifier, err := event.NotifierFactory(cfgnotifier)
				if err != nil {
					return errors.Wrapf(err, "notifier factory%s",
						logs.KVL(
							"config-event", cfgsub,
							"config-notifier", cfgnotifier,
						))
				}

				notifier.Regist(sub)
			}

			sub.Regist(newpub)
		}

		//새로운 publisher invoker 지정
		event.Invoke = newpub.Publish

		//전에있던 이벤트 Pub 종료
		//아래 붙어 있는
		if oldpub != nil {
			oldpub.Close()
		}
		//swap new->old
		oldpub = newpub

		event.PrintEventConfiguation(os.Stdout, newpub)

		return nil
	}
	resetEvent := func() error {
		//file compare
		newhash, err := hasing(filename)
		if err != nil {
			errorHandlers.OnError(errors.Wrapf(err, "event config file hasing"))
		}
		if oldhash == newhash {
			return nil //same hash
		}
		oldhash = newhash //swap hash

		return setEvent()
	}

	// 첫 등록
	if err := setEvent(); err != nil {
		return nil, err
	}
	oldhash, _ = hasing(filename)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					<-time.After(time.Second * 1)

					if err := resetEvent(); err != nil {
						errorHandlers.OnError(errors.Wrapf(err, "event file watcher: file was changed: event set"))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				errorHandlers.OnError(errors.Wrapf(err, "event file watcher: error"))
			}
		}
	}()

	if watcher.Add(filename); err != nil {
		return nil, err
	}

	//return closer
	return func() {
		if watcher != nil {
			watcher.Close()
		}
		if oldpub != nil {
			oldpub.Close()
		}
	}, nil
}

func doMigration(cfg *config.Config) (err error) {
	src := fmt.Sprintf("file://%v", cfg.Migrate.Source)
	dest := fmt.Sprintf("%v://%v", cfg.Database.Type, database.FormatDSN(cfg))

	// read migration latest file
	var latest interface {
		Version() string
		Err() error
	}
	if true {
		latest = &config.Latest{Source: cfg.Migrate.Source}
	} else {
		latest = migrations.SudoryLatest
	}
	err = latest.Err()
	if err != nil {
		return err
	}

	latest_version, err := strconv.Atoi(latest.Version())
	if err != nil {
		return
	}

	mgrt, err := database.NewMigrate(src, dest)
	if err != nil {
		return
	}
	defer func() {
		// close
		serr, derr := mgrt.Close()
		if err == nil && serr != nil {
			err = errors.WithStack(serr)
		}
		if err == nil && derr != nil {
			err = errors.WithStack(derr)
		}
	}()

	// get migreate version (current)
	cur_ver, cur_dirty, err := mgrt.Version()
	if err != nil && err != migrate.ErrNilVersion {
		err = errors.Wrapf(err, "failed to get current version")
		return
	}

	// check dirty state
	if cur_dirty {
		return &migrate.ErrDirty{Version: int(cur_ver)}
	}

	if cur_ver < uint(latest_version) {
		// do migrate goto V
		err = mgrt.Migrate(uint(latest_version))
		if err != nil && err != migrate.ErrNoChange {
			err = errors.Wrapf(err, "failed to migrate goto version=\"%v\"", latest_version)
			return
		}
	}

	// get migreate version (latest)
	new_ver, new_dirty, err := mgrt.Version()
	if err != nil && err != migrate.ErrNilVersion {
		err = errors.Wrapf(err, "failed to get current version")
		return
	}

	cols := []string{
		"",
		"driver",
		"database",
		"source",
		"version",
		"status",
		"dirty",
	}

	vals := []string{
		"-",
		cfg.Database.Type,
		cfg.Database.DBName,
		cfg.Migrate.Source,
		fmt.Sprintf("v%v", new_ver),
		func() string {
			if cur_ver == new_ver {
				return "no change"
			} else {
				return fmt.Sprintf("v%v->v%v", cur_ver, new_ver)
			}
		}(),
		strconv.FormatBool(new_dirty),
	}

	// print migrate info
	w := os.Stdout
	defer fmt.Fprintln(w, strings.Repeat("_", 40))

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	defer tw.Flush()
	fmt.Fprintln(w, "database migration:")
	tw.Write([]byte(strings.Join(cols, "\t") + "\n"))
	tw.Write([]byte(strings.Join(vals, "\t") + "\n"))

	return
}

func newGlobalVariablesCron(db *sql.DB) (func(), error) {
	const interval = 10 * time.Second

	//환경설정 updater 생성
	updator := globvar.NewGlobalVariablesUpdate(db, flavor.Dialect())
	//환경변수 리스트 검사
	if err := updator.WhiteListCheck(); err != nil {
		//빠져있는 환경변수 추가
		if err := updator.Merge(); err != nil {
			return nil, errors.Wrapf(err, "global variables init merge")
		}
	}
	//환경변수 업데이트
	if err := updator.Update(); err != nil {
		return nil, errors.Wrapf(err, "global variables init update")
	}

	//에러 핸들러 등록
	errorHandlers := ticker.HashsetErrorHandlers{}
	errorHandlers.Add(func(err error) {
		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = logs.KVL(
					"stack", s,
				)
			})
		})

		logger.Error(fmt.Errorf("cron jobs: %w%s", err, stack))
	})

	//new ticker
	tickerClose := ticker.NewTicker(interval,
		//global variables update
		func() {
			if err := updator.Update(); err != nil {
				errorHandlers.OnError(errors.Wrapf(err, "global variables update"))
			}
		},
	)

	return tickerClose, nil
}

func newPurgeDeletedDataCron(engine *xorm.Engine, respitePeriod time.Duration) (func(), error) {
	purgeDeletedrecord := func(tables []names.TableName) error {
		purge := func(tx *xorm.Session, table names.TableName) error {
			cond := builder.Lt{"deleted": time.Now().Add(respitePeriod * -1)}
			if err := database.XormDelete(
				tx.Unscoped().Where(cond), table); err != nil {
				return errors.Wrapf(err, "purge deleted event notifier status%s", logs.KVL(
					"table", table.TableName(),
				))
			}

			return nil
		}
		_, err := engine.Transaction(func(tx *xorm.Session) (interface{}, error) {
			for _, table := range tables {
				if err := purge(tx, table); err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		return err
	}
	taskMapper := func(tables [][]names.TableName, functor func(tables []names.TableName) error) []func() error {
		mapper := func(tables []names.TableName) func() error {
			return func() error {
				return functor(tables)
			}
		}

		tasks := make([]func() error, len(tables))
		for i := range tables {
			tasks[i] = mapper(tables[i])
		}
		return tasks
	}

	taskWrapper := func(tasks []func() error, errorHandler func(error)) []func() {
		wrapper := func(fn func() error) func() {
			return func() {
				if err := fn(); err != nil {
					errorHandler(err)
				}
			}
		}

		wrappedtasks := make([]func(), len(tasks))
		for i := range tasks {
			wrappedtasks[i] = wrapper(tasks[i])
		}
		return wrappedtasks
	}

	//유예 시간 확인
	if respitePeriod == 0 {
		return func() {}, nil //사용하지 않는다
	}

	purgetables := [][]names.TableName{
		{new(stepv1.ServiceStep), new(servicev1.Service)}, //transaction unit; service
		{new(channelv1.NotifierStatus)},                   //transaction unit; event notifier status
	}

	//first call
	for _, fn := range taskMapper(purgetables, purgeDeletedrecord) {
		if err := fn(); err != nil {
			return nil, errors.Wrapf(err, "first call")
		}
	}

	//에러 핸들러 등록
	errorHandlers := ticker.HashsetErrorHandlers{}
	errorHandlers.Add(func(err error) {
		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = logs.KVL(
					"stack", s,
				)
			})
		})

		logger.Error(fmt.Errorf("cron jobs: purge deleted data: %w%s", err, stack))
	})

	//new ticker
	tickerClose := ticker.NewTicker(60*time.Minute, taskWrapper(taskMapper(purgetables, purgeDeletedrecord), errorHandlers.OnError)...)

	return tickerClose, nil
}

func newEnigma(configFilename string) error {
	config := enigma.Config{}
	if err := configor.Load(&config, configFilename); err != nil {
		return errors.Wrapf(err, "read enigma config file %v",
			logs.KVL(
				"filename", configFilename,
			))
	}

	if err := enigma.LoadConfig(config); err != nil {
		b, _ := ioutil.ReadFile(configFilename)

		return errors.Wrapf(err, "load enigma config %v",
			logs.KVL(
				"filename", configFilename,
				"config", string(b),
			))
	}

	if len(config.CryptoAlgorithmSet) == 0 {
		return errors.New("'enigma cripto alg set' is empty")
	}

	for _, k := range dctv1.CiperKeyNames() {
		if _, ok := config.CryptoAlgorithmSet[k]; !ok {
			return errors.Errorf("not found enigma machine name%s",
				logs.KVL(
					"key", k,
				))
		}
	}

	enigma.PrintConfig(os.Stdout, config)

	for key := range config.CryptoAlgorithmSet {
		const quickbrownfox = `the quick brown fox jumps over the lazy dog`
		encripted, err := enigma.CipherSet(key).Encode([]byte(quickbrownfox))
		if err != nil {
			return errors.Wrapf(err, "enigma test: encode %v",
				logs.KVL("config-name", key))
		}
		plain, err := enigma.CipherSet(key).Decode(encripted)
		if err != nil {
			return errors.Wrapf(err, "enigma test: decode %v",
				logs.KVL("config-name", key))
		}

		if strings.Compare(quickbrownfox, string(plain)) != 0 {
			return errors.Errorf("enigma test: diff result %v",
				logs.KVL("config-name", key))
		}
	}

	return nil
}
