package v1_test

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func TestSericeSync(t *testing.T) {
	newEngine := func() *xorm.Engine {
		const (
			driver = "mysql"
			dsn    = "root:root@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
		)

		engine, err := xorm.NewEngine(driver, dsn)
		if ErrorHandle(err) {
			panic(err)
		}

		engine.SetLogLevel(log.LOG_DEBUG)

		return engine
	}

	sync := func() {

		engine := newEngine()

		model := new(servicev1.DbSchemaService)

		err := engine.Sync(model)

		if ErrorHandle(err) {
			t.Fatal(err)
		}
	}
	sync()
}
