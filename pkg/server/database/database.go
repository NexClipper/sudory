package database

import (
	"bytes"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/config"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql"
)

const defaultDbConnParams = "charset=utf8mb4&parseTime=True&loc=Local"

type DBManipulator struct {
	engine *xorm.Engine
}

func New(cfg *config.Config) (*DBManipulator, error) {
	db := &DBManipulator{}
	engine, err := xorm.NewEngine(cfg.Database.Type, formatDSN(cfg))
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

	//timezone setting
	engine.SetTZLocation(time.Local) //app timezone Local
	engine.SetTZDatabase(time.UTC)   //db timezone UTC

	db.engine = engine

	return db, nil
}

func (d *DBManipulator) session() *xorm.Session {
	return d.engine.NewSession()
}

func (d *DBManipulator) Close() {
	d.engine.Close()
}
func (d *DBManipulator) Engine() *xorm.Engine {
	return d.engine
}

func formatDSN(cfg *config.Config) string {
	db := cfg.Database
	var buf bytes.Buffer

	if len(db.Username) > 0 {
		buf.WriteString(db.Username)
		if len(db.Password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(db.Password)
		}
		buf.WriteByte('@')
	}

	if len(db.Protocol) > 0 {
		buf.WriteString(db.Protocol)
		if len(db.Host) > 0 && len(db.Port) > 0 {
			buf.WriteByte('(')
			buf.WriteString(db.Host + ":" + db.Port)
			buf.WriteByte(')')
		}
	}

	buf.WriteByte('/')
	buf.WriteString(db.DBName)

	if len(defaultDbConnParams) > 0 {
		buf.WriteByte('?')
		buf.WriteString(defaultDbConnParams)
	}

	return buf.String()
}
