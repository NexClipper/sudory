package main

import (
	"database/sql"
	"flag"
	"fmt"
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
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	mysqlFlavor "github.com/NexClipper/sudory/pkg/server/database/vanilla/excute/dialects/mysql"
	"github.com/NexClipper/sudory/pkg/server/event/managed_channel"
	"github.com/NexClipper/sudory/pkg/server/macro/enigma"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	dctv1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	"github.com/NexClipper/sudory/pkg/server/route"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/NexClipper/sudory/pkg/server/status/ticker"
	"github.com/NexClipper/sudory/pkg/version"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
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

	db, err := database.New(cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := doMigration(cfg); err != nil {
		panic(err)
	}

	//init managed channel
	mc := managed_channel.NewEvent(db, excute.GetSqlExcutor(mysqlFlavor.Dialect()))

	errorHandler := func(err error) {
		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = logs.KVL(
					"stack", s,
				)
			})
		})

		logger.Error(fmt.Errorf("%w%s", err, stack))
	}
	mc.ErrorHandlers.Add(errorHandler)
	mc.NofitierErrorHandlers.Add(
		managed_channel.DefaultErrorHandler_notifier(mc),
		func(n managed_channel.Notifier, err error) { managed_channel.DefaultErrorHandler(err) },
	)
	managed_channel.InvokeByChannelUuid = mc.InvokeByChannelUuid
	managed_channel.InvokeByEventCategory = mc.InvokeByEventCategory

	//init global variant cron
	cronGVClose, err := newGlobalVariablesCron(db, excute.GetSqlExcutor(mysqlFlavor.Dialect()))
	if err != nil {
		panic(err)
	}
	defer cronGVClose() //크론잡 종료

	r := route.New(cfg, db)
	if err := r.Start(); err != nil {
		nullstring := func(p *string) string {
			if p != nil {
				return *p
			}
			return "none"
		}

		var stack *string
		//stack for surface
		logs.StackIter(err, func(s string) {
			stack = &s
		})
		//stack for internal
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = &s
			})
		})

		err = errors.Wrapf(err, "failed to start route")
		logger.Errorf("%v%v", err.Error(), logs.KVL(
			"stack", nullstring(stack),
		))
	}

	logger.Debugf("%s is DONE", path.Base(strings.ReplaceAll(os.Args[0], "\\", "/")))
}

func doMigration(cfg *config.Config) (err error) {
	src := fmt.Sprintf("file://%v", cfg.Migrate.Source)
	dest := fmt.Sprintf("%v://%v", cfg.Database.Type, database.FormatDSN(cfg.Database))

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

func newGlobalVariablesCron(db *sql.DB, dialect excute.SqlExcutor) (func(), error) {
	const interval = 10 * time.Second

	//환경설정 updater 생성
	updator := globvar.NewGlobalVariablesUpdate(db, dialect)
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
