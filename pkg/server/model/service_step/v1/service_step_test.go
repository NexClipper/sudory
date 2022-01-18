package v1_test

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	servstepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func TestSericeStepSync(t *testing.T) {
	newEngine := func() *xorm.Engine {
		const (
			driver = "mysql"
			dsn    = "root:root@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
		)

		engine, err := xorm.NewEngine(driver, dsn)
		if ErrorWithHandler(err) {
			panic(err)
		}

		engine.SetLogLevel(log.LOG_DEBUG)

		return engine
	}

	sync := func() {

		engine := newEngine()

		model := new(servstepv1.DbSchemaServiceStep)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}
