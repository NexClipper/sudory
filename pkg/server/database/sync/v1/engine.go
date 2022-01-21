package v1

import (
	. "github.com/NexClipper/sudory/pkg/server/macro"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func newEngine() *xorm.Engine {
	const (
		driver = "mysql"
		dsn    = "sudory:sudory@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
	)

	engine, err := xorm.NewEngine(driver, dsn)
	if ErrorWithHandler(err) {
		panic(err)
	}

	engine.SetLogLevel(log.LOG_DEBUG)

	return engine
}
