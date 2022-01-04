package database

import (
	"strings"
	"time"

	"github.com/NexClipper/sudory-prototype-r1/pkg/server/config"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql"
)

type DBManipulator struct {
	engine *xorm.Engine
}

func New(cfg *config.Config) (*DBManipulator, error) {
	db := &DBManipulator{}
	engine, err := xorm.NewEngine(cfg.Database.Type, cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	engine.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	engine.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	engine.SetConnMaxLifetime(time.Duration(cfg.Database.MaxConnLifeTime) * time.Second)
	engine.ShowSQL(cfg.Database.ShowSQL)

	switch strings.ToLower(cfg.Database.LogLevel) {
	case "debug":
		engine.SetLogLevel(log.LOG_DEBUG)
	case "info":
		engine.SetLogLevel(log.LOG_INFO)
	case "warn":
		engine.SetLogLevel(log.LOG_WARNING)
	case "error":
		engine.SetLogLevel(log.LOG_ERR)
	default:
		engine.SetLogLevel(log.LOG_OFF)
	}
	engine.SetMapper(names.SnakeMapper{})

	db.engine = engine

	return db, nil
}

func (d *DBManipulator) session() *xorm.Session {
	return d.engine.NewSession()
}

func (d *DBManipulator) Close() {
	d.engine.Close()
}
